package router

import (
	"gitlab.com/systemz/aimpanel2/master/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"Register",
		"POST",
		"/auth/register",
		handlers.Register,
	},
	Route{
		"Login",
		"POST",
		"/auth/login",
		handlers.Login,
	},
}
