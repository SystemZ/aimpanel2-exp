package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/handler/game_server"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	v1 := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if route.AuthRequired {
			handler = AuthMiddleware(PermissionMiddleware(route.HandlerFunc))
		}

		v1.Path(route.Pattern).Handler(handler).Name(route.Name).Methods(route.Method)
	}

	return router
}

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
		"SwaggerSpec",
		"GET",
		"/swagger.json",
		handler.SwaggerSpec,
		false,
	},

	//Auth
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

	//Host
	Route{
		"List hosts",
		"GET",
		"/host",
		handler.ListHosts,
		true,
	},
	Route{
		"Get host",
		"GET",
		"/host/{id}",
		handler.GetHost,
		true,
	},
	Route{
		"Create host",
		"POST",
		"/host",
		handler.CreateHost,
		true,
	},
	Route{
		"Get host metric",
		"GET",
		"/host/{id}/metric",
		handler.GetHostMetric,
		true,
	},

	//TODO add to swagger from here

	//GameServers
	Route{
		"Create game server",
		"POST",
		"/host/{host_id}/server",
		game_server.Create,
		true,
	},
	Route{
		"List game servers by host id",
		"GET",
		"/host/{id}/server",
		game_server.ListByHostId,
		true,
	},
	Route{
		"List game servers by user id",
		"GET",
		"/host/my/server",
		game_server.ListByUser,
		true,
	},
	Route{
		"Install game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/install",
		game_server.Install,
		true,
	},
	Route{
		"Start game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/start",
		game_server.Start,
		true,
	},
	Route{
		"Restart game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/restart",
		game_server.Restart,
		true,
	},
	Route{
		"Stop game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/stop",
		game_server.Stop,
		true,
	},
	Route{
		"Send command to game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/command",
		game_server.SendCommand,
		true,
	},
	Route{
		"Game server logs",
		"PUT",
		"/host/{host_id}/server/{server_id}/logs",
		game_server.ConsoleLog,
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
		"/game",
		handler.ListGames,
		true,
	},
}
