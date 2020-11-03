package routers

import (
	authentication "ccl/ccl-auth-api/core"
	"ccl/ccl-auth-api/services"

	"github.com/gorilla/mux"

	"github.com/urfave/negroni"
)

func SetAdminRoutes(router *mux.Router) *mux.Router {
	router.Handle("/openadmin/g/u",
		negroni.New(
			negroni.HandlerFunc(services.GetUser()),
		)).Methods("POST")

	router.Handle("/openadmin/g/us",
		negroni.New(
			negroni.HandlerFunc(services.GetUsers()),
		)).Methods("POST")

	router.Handle("/openadmin/p/u",
		negroni.New(
			negroni.HandlerFunc(services.CreateUser()),
		)).Methods("POST")

	/*router.Handle("/openadmin/p/u/r",
	negroni.New(
		negroni.HandlerFunc(services.RegisterUser()),
	)).Methods("POST")*/

	router.Handle("/openadmin/d/u",
		negroni.New(
			negroni.HandlerFunc(services.DeleteUser()),
		)).Methods("POST")

	router.Handle("/openadmin/u/us",
		negroni.New(
			negroni.HandlerFunc(services.UpdateUsers()),
		)).Methods("POST")

	router.Handle("/openadmin/u/u",
		negroni.New(
			negroni.HandlerFunc(services.UpdateUser()),
		)).Methods("POST")

	router.Handle("/test",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireAdminPermissions),
		)).Methods("POST")

	router.Handle("/admin/g/u",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(services.GetUser()),
		)).Methods("POST")

	router.Handle("/admin/g/us",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(services.GetUsers()),
		)).Methods("POST")

	router.Handle("/admin/p/u",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireAdminPermissions),
			negroni.HandlerFunc(services.CreateUser()),
		)).Methods("POST")

	/*router.Handle("/admin/p/u/r",
	negroni.New(
		negroni.HandlerFunc(authentication.RequireTokenAuthentication),
		negroni.HandlerFunc(services.RegisterUser()),
	)).Methods("POST")*/

	router.Handle("/admin/d/u",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(authentication.RequireAdminPermissions),
			negroni.HandlerFunc(services.DeleteUser()),
		)).Methods("POST")

	router.Handle("/admin/u/us",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(services.UpdateUsers()),
		)).Methods("POST")

	router.Handle("/admin/u/u",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(services.UpdateUser()),
		)).Methods("POST")

	return router
}
