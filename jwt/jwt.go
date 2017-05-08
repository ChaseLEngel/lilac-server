package jwt

import (
	"crypto/rsa"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
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
