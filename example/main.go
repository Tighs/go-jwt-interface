package main

import (
	"net/http"
	"go-jwt-interface/api"
	"encoding/json"
	"go-jwt-interface/example/mockimpl"
	"github.com/julienschmidt/httprouter"
	"log"
)

func main(){

	provider := mockimpl.UserProviderImpl{}
	provider.Init()

	/*router := mux.NewRouter()
	router.HandleFunc("/normal",Testhandler)
	router.HandleFunc("/test",Testhandler)
	router.HandleFunc("/blubb/dorn",Testhandler)
	s := router.PathPrefix("/foo").Subrouter()
	s.HandleFunc("/bar",Testhandler)
	api.ManageMuxRouter(provider,router,s)


	http.ListenAndServe(":8080", router)*/

	router := httprouter.New()
	api.GenerateLoginEndpointForHttprouter(router)
	security := api.ManageNegroni(provider)
	security.UseHandlerFunc(Testhandler)
	router.Handler("GET","/secured",security)
	router.Handler("GET","/test", security)
	router.Handler("GET","/foo", security)
	router.Handler("GET","/foo/bar",security)

	log.Fatal(http.ListenAndServe(":8080", router))

}


func Testhandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SUCCESS"))
}

func endpoint(w http.ResponseWriter, r *http.Request){
	jwtUser,err := api.ExtractUserFromValidToken(r)

	if err == nil {
		json.NewEncoder(w).Encode(jwtUser)
	}
}




