package router

import (
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/handler/game_server"
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

	//GameServers
	Route{
		"Create game server",
		"POST",
		"/hosts/{host_id}/servers",
		game_server.Create,
		true,
	},
	Route{
		"List game servers by host id",
		"GET",
		"/hosts/{id}/servers",
		game_server.ListByHostId,
		true,
	},
	Route{
		"List game servers by user id",
		"GET",
		"/hosts/my/servers",
		game_server.ListByUser,
		true,
	},
	Route{
		"Install game server",
		"GET",
		"/hosts/{host_id}/servers/{server_id}/install",
		game_server.Install,
		true,
	},
	Route{
		"Start game server",
		"GET",
		"/hosts/{host_id}/servers/{server_id}/start",
		game_server.Start,
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

	//Games
	Route{
		"List games",
		"GET",
		"/games",
		handler.ListGames,
		true,
	},
}
