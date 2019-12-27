package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/handler/gs"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	//router.StrictSlash(true)

	v1 := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		var handler http.Handler
		handler = CommonMiddleware(route.HandlerFunc)

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
		handler.GetSpec,
		false,
		false,
	},
	Route{
		"SwaggerUi",
		"GET",
		"/swagger",
		handler.GetSwaggerUi,
		false,
		false,
	},
	Route{
		"SwaggerDocsRedirect",
		"GET",
		"/docs",
		handler.GetDocsRedirect,
		false,
		false,
	},
	Route{
		"SwaggerDocs",
		"GET",
		"/docs/",
		handler.GetDocs,
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
		gs.Data,
		true,
		true,
	},
	Route{
		"Game server data",
		"POST",
		"/events/{host_token}/{server_id}",
		gs.Data,
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
		"Host List",
		"GET",
		"/host",
		handler.HostList,
		true,
		false,
	},
	Route{
		"Host Authentication",
		"GET",
		"/host/auth/{token}",
		handler.HostAuth,
		false,
		false,
	},
	Route{
		"Host Details",
		"GET",
		"/host/{id}",
		handler.HostDetails,
		true,
		false,
	},
	Route{
		"Host Remove",
		"DELETE",
		"/host/{id}",
		handler.HostRemove,
		true,
		false,
	},
	Route{
		"Host Create",
		"POST",
		"/host",
		handler.HostCreate,
		true,
		false,
	},
	Route{
		"Host Update",
		"GET",
		"/host/{id}/update",
		handler.HostUpdate,
		true,
		false,
	},
	Route{
		"Host Metric",
		"GET",
		"/host/{id}/metric",
		handler.HostMetric,
		true,
		false,
	},

	//GameServers
	Route{
		"List game servers by user id",
		"GET",
		"/host/my/server",
		gs.ListByUser,
		true,
		false,
	},
	Route{
		"Get game server",
		"GET",
		"/host/{host_id}/server/{server_id}",
		gs.Get,
		true,
		false,
	},
	Route{
		"Remove game server",
		"DELETE",
		"/host/{host_id}/server/{server_id}",
		gs.Remove,
		true,
		false,
	},
	Route{
		"Create game server",
		"POST",
		"/host/{host_id}/server",
		gs.Create,
		true,
		false,
	},
	Route{
		"List game servers by host id",
		"GET",
		"/host/{id}/server",
		gs.ListByHostId,
		true,
		false,
	},
	Route{
		"Install game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/install",
		gs.Install,
		true,
		false,
	},
	Route{
		"Start game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/start",
		gs.Start,
		true,
		false,
	},
	Route{
		"Restart game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/restart",
		gs.Restart,
		true,
		false,
	},
	Route{
		"Stop game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/stop",
		gs.Stop,
		true,
		false,
	},
	Route{
		"Send command to game server",
		"PUT",
		"/host/{host_id}/server/{server_id}/command",
		gs.SendCommand,
		true,
		false,
	},
	Route{
		"Game server logs",
		"GET",
		"/host/{host_id}/server/{server_id}/logs",
		gs.ConsoleLog,
		true,
		false,
	},
	Route{
		"Game server logs",
		"GET",
		"/host/{host_id}/server/{server_id}/logs",
		gs.ConsoleLog,
		true,
		true,
	},

	//User
	Route{
		"Change password",
		"POST",
		"/user/change_password",
		handler.UserChangePassword,
		true,
		false,
	},
	Route{
		"Change email",
		"POST",
		"/user/change_email",
		handler.UserChangeEmail,
		true,
		false,
	},
	Route{
		"User profile",
		"GET",
		"/me",
		handler.UserProfile,
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
