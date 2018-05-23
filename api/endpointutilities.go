package api

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
)

var provider UserProvider

func LoadUserProvider(p UserProvider){
	provider = p
}

func AdjustMuxEndpoints(router *mux.Router){

	router.Use(secureRoute)

	xmlEndpoints := provideEndpoints()
	for _,endpoint := range xmlEndpoints.EndpointList {
		router.NewRoute().Path(endpoint.Path).Methods(endpoint.Method)
	}
}

func provideLoginEndpoint(w http.ResponseWriter, req *http.Request){

	var user loginCredentials
	var matchFlag bool

	json.NewDecoder(req.Body).Decode(&user)
	if provider != nil{
		if user.Username == "" && user.Password == ""{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Could not parse login credentials"))
		}else{
			for _,validUser := range provider.ProvideValidUsers(){
				if user.Username == validUser.Username && user.Password == validUser.HashedPassword {
					matchFlag = true
					json.NewEncoder(w).Encode(generateToken(mapMinimalUserToInternalUser(validUser)))
				}
			}
			if !matchFlag{
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	}else{
		log.Fatal("You have to load/inject a UserProvider via the LoadUserProvider function")
	}
}

func provideSecuredEndpoint(w http.ResponseWriter, req *http.Request) error{

	_,err := verifyTokenExtractClaims(req)

	return err
}

//Middleware to secure the given route
func secureRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoints := findEndpoints(r.RequestURI)

		for _,endpoint := range endpoints{
			if endpoint.Type == "secured"{
				err := provideSecuredEndpoint(w,r)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
				}else {
					h.ServeHTTP(w, r)
				}
			}else if endpoint.Type =="login"{
				provideLoginEndpoint(w,r)
			}else{
				h.ServeHTTP(w,r)
			}
		}
	})
}

func mapMinimalUserToInternalUser(user MinimalUser) JWTUser {

	var jwtUser JWTUser
	jwtUser.ID = user.DatabaseId
	jwtUser.Role = user.Role
	return jwtUser

}

type UserProvider interface {

	ProvideValidUsers() []MinimalUser
}

type MinimalUser struct{
	DatabaseId     uint64
	Username       string
	HashedPassword string
	Role           []string
}
