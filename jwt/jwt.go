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

func Init() (*JwtData, error) {
	data := new(JwtData)
	var err error
	privateKeyBytes, err := ioutil.ReadFile("lilac.rsa")
	if err != nil {
		return nil, err
	}
	data.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKeyBytes, err := ioutil.ReadFile("lilac.rsa.pub")
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

func (data *JwtData) Authenticate(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	authorization := strings.Split(authHeader, " ")

	if len(authorization) != 2 {
		return false
	}

	token, err := jwtgo.Parse(authorization[1], func(token *jwtgo.Token) (interface{}, error) {
		if token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("Incorrent JWT signing method: %v\n", token.Method.Alg())
		}
		return data.publicKey, nil
	})

	if err != nil {
		//fmt.Printf("Failed to parse JWT token from request: %v\n", err)
		return false
	}

	if token.Valid {
		return true
	}

	return false
}
