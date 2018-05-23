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

	keyPath := xmlKey.Path + xmlKey.Name
	pubkeyPath := keyPath+".pub"

	signingKey, err = ioutil.ReadFile(keyPath)
	verifyKey, err = ioutil.ReadFile(pubkeyPath)

	if err != nil {
		log.Fatal("Could not read rsa signingKey-pair")
		return
	}
}

func GenerateToken(user JWTUser) Token{

	key,_ := jwt.ParseRSAPrivateKeyFromPEM(signingKey)

	signer := jwt.New(jwt.SigningMethodRS256)

	signer.Claims = &CustomClaims{
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

	token := Token{Token:tokenString}
	return token
}

func VerifyTokenExtractClaims(req *http.Request) (CustomClaims,error) {

	claims := CustomClaims{}
	var err error

	tokenString := extractJWTFromHttpRequest(req)

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

type CustomClaims struct{

	*jwt.StandardClaims
	TokenType string
	JWTUser
}

type LoginCredentials struct{

	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims

}

type JWTUser struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	FirstName string `json:"firstName"`
	Title string `json:"title"`
	Age int `json:"age"`
	Address string `json:"address"`
	Position string `json:"position"`
	Department string `json:"department"`
	Additional []string `json:"additional"`
}

type Token struct {
	Token string `json:"Token"`
}