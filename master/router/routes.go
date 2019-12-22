package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/events"
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
			if route.SlaveOnly {
				handler = SlavePermissionMiddleware(handler)
			} else {
				handler = AuthMiddleware(PermissionMiddleware(handler))
			}
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
	SlaveOnly    bool
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handler.Index,
		true,
		false,
	},
	Route{
		"SwaggerSpec",
		"GET",
		"/swagger.json",
		handler.SwaggerSpec,
		false,
		false,
	},

	Route{
		"New version",
		"POST",
		"/version",
		handler.NewVersion,
		false,
		false,
	},

	Route{
		"Host events",
		"GET",
		"/events/{host_token}",
		events.Handler,
		true,
		true,
	},
	Route{
		"Game server events",
		"GET",
		"/events/{host_token}/{server_id}",
		events.Handler,
		true,
		true,
	},
	Route{
		"Host data",
		"POST",
		"/events/{host_token}",
		game_server.Data,
		true,
		true,
	},
	Route{
		"Game server data",
		"POST",
		"/events/{host_token}/{server_id}",
		game_server.Data,
		true,
		true,
	},

	//Auth
	Route{
		"Register",
		"POST",
		"/auth/register",
		handler.Register,
		false,
		false,
	},
	Route{
		"Login",
		"POST",
		"/auth/login",
		handler.Login,
		false,
		false,
	},

	//Host
	Route{
		"List hosts",
		"GET",
		"/host",
		handler.ListHosts,
		true,
		false,
	},
	Route{
		"Host auth",
		"GET",
		"/host/auth/{token}",
		handler.HostAuth,
		false,
		false,
	},
	Route{
		"Get host",
		"GET",
		"/host/{id}",
		handler.GetHost,
		true,
		false,
	},
	Route{
		"Remove host",
		"DELETE",
		"/host/{id}",
		handler.RemoveHost,
		true,
		false,
	},
	Route{
		"Create host",
		"POST",
		"/host",
		handler.CreateHost,
		true,
		false,
	},
	Route{
		"Get host metric",
		"GET",
		"/host/{id}/metric",
		handler.GetHostMetric,
		true,
		false,
	},
	//TODO add to swagger from here

	//GameServers
	Route{
		"List game servers by user id",
		"GET",
		"/host/my/server",
		game_server.ListByUser,
		true,
		false,
	},
	Route{
		"Get game server",
		"GET",
		"/host/{host_id}/server/{server_id}",
		game_server.Get,
		true,
		false,
	},
	Route{
		"Remove game server",
		"DELETE",
		"/host/{host_id}/server/{server_id}",
		game_server.Remove,
		true,
		false,
	},
	Route{
		"Create game server",
		"POST",
		"/host/{host_id}/server",
		game_server.Create,
		true,
		false,
	},
	Route{
		"List game servers by host id",
		"GET",
		"/host/{id}/server",
		game_server.ListByHostId,
		true,
		false,
	},
	Route{
		"Install game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/install",
		game_server.Install,
		true,
		false,
	},
	Route{
		"Start game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/start",
		game_server.Start,
		true,
		false,
	},
	Route{
		"Restart game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/restart",
		game_server.Restart,
		true,
		false,
	},
	Route{
		"Stop game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/stop",
		game_server.Stop,
		true,
		false,
	},
	Route{
		"Send command to game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/command",
		game_server.SendCommand,
		true,
		false,
	},
	Route{
		"Game server logs",
		"GET",
		"/host/{host_id}/server/{server_id}/logs",
		game_server.ConsoleLog,
		true,
		false,
	},
	Route{
		"Game server logs",
		"GET",
		"/host/{host_id}/server/{server_id}/logs",
		game_server.ConsoleLog,
		true,
		true,
	},

	//User
	Route{
		"Change password",
		"POST",
		"/user/change_password",
		handler.ChangePassword,
		true,
		false,
	},
	Route{
		"Change email",
		"POST",
		"/user/change_email",
		handler.ChangeEmail,
		true,
		false,
	},
	Route{
		"User profile",
		"GET",
		"/user/profile",
		handler.Profile,
		true,
		false,
	},

	//Games
	Route{
		"List games",
		"GET",
		"/game",
		handler.ListGames,
		true,
		false,
	},
}
