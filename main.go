package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)

const (
	keyPath = "/home/tighs/Development/go/src/gojwt/keys/app.rsa"
	pubkeyPath = "/home/tighs/Development/go/src/gojwt/keys/app.rsa.pub"
)

var signingKey, verifyKey []byte

func main(){

	LoadKeys()

	router := mux.NewRouter()
	router.HandleFunc("/login",LoginEndpoint).Methods("POST")
	router.HandleFunc("/secured",SecuredEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080",router))

}


func LoadKeys(){
	var err error
	signingKey, err = ioutil.ReadFile(keyPath)
	verifyKey, err = ioutil.ReadFile(pubkeyPath)

	if err != nil {
		log.Fatal("Could not read rsa signingKey-pair")
		return
	}
}

func generateToken() Token{

	key,_ := jwt.ParseRSAPrivateKeyFromPEM(signingKey)

	signer := jwt.New(jwt.SigningMethodRS256)

	signer.Claims = &CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
		"Level1",
		User{"Name",13},
	}

	tokenString,err := signer.SignedString(key)

	if err != nil {
		log.Printf("Error signing token: %v\n", err)
	}

	token := Token{Token:tokenString}
	return token
}

func LoginEndpoint(w http.ResponseWriter, req *http.Request){

	var user LoginCredentials

	json.NewDecoder(req.Body).Decode(&user)

	if user.Username == "test" && user.Password == "test"{
		token := generateToken()
		json.NewEncoder(w).Encode(token)
	}
}

func SecuredEndpoint(w http.ResponseWriter, req *http.Request){

	key,_ := jwt.ParseRSAPublicKeyFromPEM(verifyKey)

	tokenString := req.Header["Authorization"][0][7:]

	claims := CustomClaims{}

	_,err := jwt.ParseWithClaims(tokenString,&claims,func(token *jwt.Token)(interface {},error){
		return key,nil
	})

	fmt.Println(claims)

	if err == nil {
		w.WriteHeader(http.StatusOK)

	}else{
		w.WriteHeader(http.StatusUnauthorized)
	}
}

type CustomClaims struct{

	*jwt.StandardClaims
	TokenType string
	User
}

type LoginCredentials struct{

	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims

}

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

type Token struct {
	Token string `json:"Token"`
}

type Response struct {
	Data string `json:"data"`
}

