package router

import (
	"gitlab.com/systemz/aimpanel2/master/handler"
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
		handler.Index,
		true,
	},
	Route{
		"Register",
		"POST",
		"/auth/register",
		handler.Register,
		false,
	},
	Route{
		"Login",
		"POST",
		"/auth/login",
		handler.Login,
		false,
	},

	//Hosts
	Route{
		"List hosts",
		"GET",
		"/hosts",
		handler.ListHosts,
		true,
	},
	Route{
		"Get host",
		"GET",
		"/hosts/{id}",
		handler.GetHost,
		true,
	},
	Route{
		"Create host",
		"POST",
		"/hosts",
		handler.CreateHost,
		true,
	},

	//User
	Route{
		"Change password",
		"POST",
		"/user/change_password",
		handler.ChangePassword,
		true,
	},
	Route{
		"Change email",
		"POST",
		"/user/change_email",
		handler.ChangeEmail,
		true,
	},
	Route{
		"User profile",
		"GET",
		"/user/profile",
		handler.Profile,
		true,
	},
}
