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
		var h http.Handler
		h = CommonMiddleware(ExitMiddleware(route.HandlerFunc))

		if route.AuthRequired {
			if route.SlaveOnly {
				h = SlavePermissionMiddleware(h)
			} else {
				h = AuthMiddleware(PermissionMiddleware(h))
			}
		}

		v1.Path(route.Pattern).Handler(h).Name(route.Name).Methods(route.Method)
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
		false,
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
	Route{
		"Host Create Job",
		"POST",
		"/host/{id}/job",
		handler.HostCreateJob,
		true,
		false,
	},
	Route{
		"Host Jobs",
		"GET",
		"/host/{id}/job",
		handler.HostJobList,
		true,
		false,
	},
	Route{
		"Host Job Remove",
		"DELETE",
		"/host/{id}/job/{job_id}",
		handler.HostJobRemove,
		true,
		false,
	},

	//GameServers
	Route{
		"GameServer ListByUser",
		"GET",
		"/host/my/server",
		gs.ListByUser,
		true,
		false,
	},
	Route{
		"GameServer Details",
		"GET",
		"/host/{host_id}/server/{server_id}",
		gs.Get,
		true,
		false,
	},
	Route{
		"GameServer Remove",
		"DELETE",
		"/host/{host_id}/server/{server_id}",
		gs.Remove,
		true,
		false,
	},
	Route{
		"GameServer Create",
		"POST",
		"/host/{host_id}/server",
		gs.Create,
		true,
		false,
	},
	Route{
		"GameServer ListByHostId",
		"GET",
		"/host/{id}/server",
		gs.ListByHostId,
		true,
		false,
	},
	Route{
		"GameServer Install",
		"PUT",
		"/host/{host_id}/server/{server_id}/install",
		gs.Install,
		true,
		false,
	},
	Route{
		"GameServer Start",
		"PUT",
		"/host/{host_id}/server/{server_id}/start",
		gs.Start,
		true,
		false,
	},
	Route{
		"GameServer Restart",
		"PUT",
		"/host/{host_id}/server/{server_id}/restart",
		gs.Restart,
		true,
		false,
	},
	Route{
		"GameServer Stop",
		"PUT",
		"/host/{host_id}/server/{server_id}/stop",
		gs.Stop,
		true,
		false,
	},
	Route{
		"GameServer Send Command",
		"PUT",
		"/host/{host_id}/server/{server_id}/command",
		gs.SendCommand,
		true,
		false,
	},
	Route{
		"GameServer Console Log",
		"GET",
		"/host/{host_id}/server/{server_id}/logs",
		gs.ConsoleLog,
		true,
		false,
	},
	Route{
		"GameServer Put Logs",
		"PUT",
		"/host/{host_id}/server/{server_id}/logs",
		gs.PutLogs,
		true,
		true,
	},
	Route{
		"GameServer Console",
		"GET",
		"/host/{host_id}/server/{server_id}/console",
		events.Handler,
		true,
		false,
	},

	//GameServer Files
	Route{
		"GameServer File list",
		"GET",
		"/host/{host_id}/server/{server_id}/file/list",
		gs.FileList,
		true,
		false,
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
