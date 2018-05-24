package api

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

var provider UserProvider

func ManageMuxRouter(p UserProvider, routerList ...*mux.Router){

	provider = p

	loginEndpoint := findLoginEndpoint()
	if len(routerList) == 0{
		log.Fatal("no router found")
	}
	routerList[0].HandleFunc(loginEndpoint, createLoginEndpoint)

	for _, router := range routerList{

		router.Use(secureRouteMiddleware)
	}
}

func ManageNegroni(p UserProvider) *negroni.Negroni{

	provider = p

	var secured negroni.Negroni

	secured.Use(negroni.HandlerFunc(secureRouteNegroniMiddleware))

	return &secured
}

func GenerateLoginEndpointForHttprouter(router *httprouter.Router){

	loginEndpoint := findLoginEndpoint()

	router.POST(loginEndpoint, httpRouterHandle())
}

func createLoginEndpoint(w http.ResponseWriter, r *http.Request){
	var user loginCredentials
	var matchFlag bool

	json.NewDecoder(r.Body).Decode(&user)
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
		log.Fatal("You have to load/inject a UserProvider via the loadUserProvider function")
	}
}

func httpRouterHandle() httprouter.Handle{

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		createLoginEndpoint(w,r)
	}
}

func securedEndpointHandler(req *http.Request) error{

	_,err := verifyTokenExtractClaims(req)

	return err
}

//middleware to secure all routes after the given router/subRouter
func secureRouteMiddleware(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if isSecuredPartRoute(r.RequestURI) || isSecuredEndpoint(r.RequestURI){
			err := securedEndpointHandler(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}else if h != nil{
				h.ServeHTTP(w,r)
			}
		}else if h != nil{
			h.ServeHTTP(w,r)
		}
	})
}

func secureRouteNegroniMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if isSecuredPartRoute(r.RequestURI) || isSecuredEndpoint(r.RequestURI){
		err := securedEndpointHandler(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}else {
			next(w,r)
		}
	}else{
		next(w,r)
	}
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
