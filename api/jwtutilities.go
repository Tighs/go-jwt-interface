package api

import (
	"io/ioutil"
	"log"
	"github.com/dgrijalva/jwt-go"
	"time"
	"net/http"
	"errors"
)

var signingKey, verifyKey []byte

func init(){

	loadKeys()
}


func loadKeys(){
	var err error

	keyPath := key.Path + key.Name
	pubkeyPath := keyPath+".pub"

	signingKey, err = ioutil.ReadFile(keyPath)
	verifyKey, err = ioutil.ReadFile(pubkeyPath)

	if err != nil {
		log.Fatal("Could not read rsa signingKey-pair")
		return
	}
}

func generateToken(user JWTUser) token {

	key,_ := jwt.ParseRSAPrivateKeyFromPEM(signingKey)

	signer := jwt.New(jwt.SigningMethodRS256)

	signer.Claims = &customClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
		"Level1",
		user,
	}

	tokenString,err := signer.SignedString(key)

	if err != nil {
		log.Printf("Error signing token: %v\n", err)
	}

	token := token{Token:tokenString}
	return token
}

func verifyTokenExtractClaims(r *http.Request) (customClaims,error) {

	claims := customClaims{}
	var err error

	tokenString := extractJWTFromHttpRequest(r)

	if len(tokenString) >0 {

		key,_ := jwt.ParseRSAPublicKeyFromPEM(verifyKey)

		_,err = jwt.ParseWithClaims(tokenString,&claims,func(token *jwt.Token)(interface {},error){
			return key,nil
		})
	}else{
		err = errors.New("could not extract valid jwt")
	}


	return claims,err
}

func extractJWTFromHttpRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	var tokenString string

	if len(auth) > 7{
		tokenString = auth[7:]

	}
	return tokenString
}

func ExtractUserFromValidToken(r * http.Request) (JWTUser,error){
	claims,err := verifyTokenExtractClaims(r)

	if err == nil {
		return claims.JWTUser,nil
	}
	return JWTUser{},err
}

type customClaims struct{

	*jwt.StandardClaims
	TokenType string
	JWTUser
}

type loginCredentials struct{

	Username string `json:"Username"`
	Password string `json:"password"`
	jwt.StandardClaims

}

type JWTUser struct {
	ID   uint64  `json:"id"`
	Role []string `json:"Role"`
}

type token struct {
	Token string `json:"token"`
}