package routers

import (
	"ccl/ccl-auth-api/controllers"


	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetHelloRoutes(router *mux.Router) *mux.Router {
	router.Handle("/test/hello",
		negroni.New(
			negroni.HandlerFunc(controllers.HelloController),
		)).Methods("GET")

	return router
}
