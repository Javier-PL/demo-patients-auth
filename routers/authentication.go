package routers

import (
	"ccl/ccl-auth-api/controllers"
	authentication "ccl/ccl-auth-api/core"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

//func SetAuthenticationRoutes(router *mux.Router, client *mongo.Client) *mux.Router {
//router.HandleFunc("/token-auth", controllers.Login(client)).Methods("POST")
func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/token-auth", controllers.Login).Methods("POST")
	router.Handle("/refresh-token-auth",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("GET")
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("GET")
	return router
}
