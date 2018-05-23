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

	router.Use(SecureRoute)

	xmlEndpoints := provideEndpoints()
	for _,endpoint := range xmlEndpoints.Endpoints{
		router.NewRoute().Path(endpoint.Path).Methods(endpoint.Method)
	}
}

func provideLoginEndpoint(w http.ResponseWriter, req *http.Request){

	var user LoginCredentials
	var matchFlag bool

	json.NewDecoder(req.Body).Decode(&user)
	if provider != nil{
		if user.Username == "" && user.Password == ""{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Could not parse login credentials"))
		}else{
			for _,validUser := range provider.ProvideValidUsers(){
				if user.Username == validUser.username && user.Password == validUser.hashedPassword{
					matchFlag = true
					json.NewEncoder(w).Encode(GenerateToken(mapMinimalUserToInternalUser(validUser)))
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

	_,err := VerifyTokenExtractClaims(req)

	return err
}

//Middleware to secure the given route
func SecureRoute(h http.Handler) http.Handler {
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

func GenerateMinimalUser(id int64, username string ,hashedPassword string) MinimalUser {
	return MinimalUser{id,username,hashedPassword}
}

func mapMinimalUserToInternalUser(user MinimalUser) JWTUser {

	var jwtUser JWTUser
	jwtUser.Id = user.databaseId
	return jwtUser

}

type UserProvider interface {

	ProvideValidUsers() []MinimalUser
}

type MinimalUser struct{
	databaseId int64
	username string
	hashedPassword string
}
