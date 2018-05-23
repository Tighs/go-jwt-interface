package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go-jwt-interface/api"
	"go-jwt-interface/example/mockimpl"
	"encoding/json"
)

func main(){

	provider := mockimpl.UserProviderImpl{}
	provider.Init()
	api.LoadUserProvider(provider)

	router := mux.NewRouter()
	router.HandleFunc("/secured",endpoint)
	router.HandleFunc("/foo",endpoint)
	router.HandleFunc("/bar", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("blubb"))
	})
	api.AdjustMuxEndpoints(router)
	log.Fatal(http.ListenAndServe(":8080",router))
}

func endpoint(w http.ResponseWriter, r *http.Request){
	jwtUser,err := api.ExtractUserFromValidToken(r)

	if err == nil {
		json.NewEncoder(w).Encode(jwtUser)
	}
}




