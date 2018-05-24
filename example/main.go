package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"go-jwt-interface/api"
	"go-jwt-interface/example/mockimpl"
	"encoding/json"
)

func main(){

	provider := mockimpl.UserProviderImpl{}
	provider.Init()

	router := mux.NewRouter()
	router.HandleFunc("/foo",testEndpoint)
	router.HandleFunc("/secured/test",testEndpoint)
	router.HandleFunc("/secured",testEndpoint)
	router.HandleFunc("/test/test",testEndpoint)
	s := router.PathPrefix("/foo").Subrouter()
	s.HandleFunc("/bar",testEndpoint)
	api.ManageMuxRouter(provider,router,s)


	/*router.HandleFunc("/foo",testEndpoint)
	router.HandleFunc("/secured/foo",testEndpoint)
	s := router.PathPrefix("/foo").Subrouter()
	s.HandleFunc("/bar",testEndpoint)
	s.HandleFunc("/bar2",testEndpoint)
	s.HandleFunc("/bar/test",testEndpoint)
	api.secureMuxRouterEndpoints(router)
	api.secureRoutes(s)*/
	http.ListenAndServe(":8080", router)
}



func testEndpoint(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	w.Write([]byte("SUCCESS!"))
}

func endpoint(w http.ResponseWriter, r *http.Request){
	jwtUser,err := api.ExtractUserFromValidToken(r)

	if err == nil {
		json.NewEncoder(w).Encode(jwtUser)
	}
}




