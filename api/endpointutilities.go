package api

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
)

var provider UserProvider
var securedRouteMap map[string]bool

func ManageMuxRouter(p UserProvider, routerList ...*mux.Router){

	provider = p

	endpoint,err := findLoginEndpoint()
	var loginEndpoint string
	if err == nil {
		loginEndpoint = endpoint.Path
	}else{
		loginEndpoint = "/login"
	}
	routerList[0].HandleFunc(loginEndpoint, provideLoginEndpoint)
	securedRouteMap = make(map[string]bool)

	for _, router := range routerList{
		initSecurePathList(router)

		router.Use(secureRoute)
	}
	log.Println(securedRouteMap)
}

func initSecurePathList(router *mux.Router){
	for _,route := range readRoutes(router){
		tpl,err := route.GetPathTemplate()
		if err == nil {
			if isSecuredPartRoute(tpl) || isSecuredEndpoint(tpl) {
				securedRouteMap[tpl] = true
			}else{
				securedRouteMap[tpl] = false
			}
		}
	}
}

func readRoutes(router *mux.Router) []mux.Route{

	var routes []mux.Route

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		routes = append(routes,*route)
		return nil
	})

	return routes
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
		log.Fatal("You have to load/inject a UserProvider via the loadUserProvider function")
	}
}

func provideSecuredEndpoint(req *http.Request) error{

	_,err := verifyTokenExtractClaims(req)

	return err
}

//middleware to secure all routes after the given router/subRouter
func secureRoute(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if securedRouteMap[r.RequestURI]{
			err := provideSecuredEndpoint(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}else if h != nil{
				h.ServeHTTP(w,r)
			}
		}else{
			h.ServeHTTP(w,r)
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
