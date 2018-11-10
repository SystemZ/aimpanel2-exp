package router

import (
	"gitlab.com/systemz/aimpanel2/master/handlers"
	"net/http"
)

type Route struct {
	Name         string
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
	AuthRequired bool
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
		true,
	},
	Route{
		"Register",
		"POST",
		"/auth/register",
		handlers.Register,
		false,
	},
	Route{
		"Login",
		"POST",
		"/auth/login",
		handlers.Login,
		false,
	},

	//Hosts
	Route{
		"List hosts",
		"GET",
		"/hosts",
		handlers.ListHosts,
		false,
	},
	Route{
		"Get host",
		"GET",
		"/hosts/{id}",
		handlers.GetHost,
		false,
	},
	Route{
		"Create host",
		"POST",
		"/hosts",
		handlers.CreateHost,
		false,
	},
}
