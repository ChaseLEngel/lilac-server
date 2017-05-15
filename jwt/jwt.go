package jwt

import (
	"crypto/rsa"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strings"
)

type JwtData struct {
	token       *jwtgo.Token
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	TokenString string
}

func Init(privateKeyFile, publicKeyFile string) (*JwtData, error) {
	data := new(JwtData)
	var err error
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	data.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKeyBytes, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return nil, err
	}
	data.publicKey, err = jwtgo.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	data.token = jwtgo.New(jwtgo.SigningMethodRS256)
	data.TokenString, err = data.token.SignedString(data.privateKey)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (data *JwtData) Authenticate(r *http.Request) (bool, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, nil
	}

	authorization := strings.Split(authHeader, " ")

	if len(authorization) != 2 {
		return false, nil
	}

	token, err := jwtgo.Parse(authorization[1], func(token *jwtgo.Token) (interface{}, error) {
		if token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("Incorrent JWT signing method: %v\n", token.Method.Alg())
		}
		return data.publicKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("Failed to parse token: %v", err)
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}
