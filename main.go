package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go-jwt-interface/api"
	"go-jwt-interface/mockImpl"
)

func main(){

	provider := mockImpl.UserProviderImpl{}
	provider.Init()
	api.LoadUserProvider(provider)

	router := mux.NewRouter()
	router.HandleFunc("/secured",endpoint)
	router.HandleFunc("/foo",endpoint)
	api.AdjustMuxEndpoints(router)
	log.Fatal(http.ListenAndServe(":8080",router))
}

func endpoint(w http.ResponseWriter, req *http.Request){
	w.Write([]byte("TEST"))
}


