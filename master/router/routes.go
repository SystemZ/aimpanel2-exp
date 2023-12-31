package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/handler"
	"gitlab.com/systemz/aimpanel2/master/handler/gs"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	//router.StrictSlash(true)

	// frontend
	frontendDir := config.HTTP_FRONTEND_DIR
	// frontend home
	router.HandleFunc("/", handler.Index).Methods("GET", "HEAD")
	router.HandleFunc("/robots.txt", handler.RobotsTxt).Methods("GET", "HEAD")
	router.HandleFunc("/service-worker.js", handler.ServiceWorker).Methods("GET", "HEAD")
	// frontend assets
	frontendAssetDirs := []string{"css", "js", "fonts", "img"}
	for _, fAssetDir := range frontendAssetDirs {
		dir := "/" + fAssetDir + "/"
		// FIXME e2e test for FS disclosure bugs
		router.PathPrefix(dir).Handler(
			http.StripPrefix(dir,
				http.FileServer(
					http.Dir(frontendDir+dir),
				),
			),
		).Methods("GET", "HEAD")
	}
	// slave deployment
	router.HandleFunc("/i/{hostToken}", handler.DeploymentScript).Methods("GET", "HEAD")

	// API endpoints
	v1 := router.PathPrefix("/v1").Subrouter()
	for _, route := range v1Routes {
		var h http.Handler
		h = CommonMiddleware(DBCheckMiddleware(ExitMiddleware(route.HandlerFunc)))

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

// /v1/*
var v1Routes = Routes{
	// API docs
	Route{
		"Index",
		"GET",
		"",
		handler.GetDocsRedirect,
		false,
		false,
	},
	Route{
		"Index",
		"GET",
		"/",
		handler.GetDocsRedirect,
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
	// SSE endpoint for browser
	Route{
		"GameServer Console",
		"GET",
		"/host/{hostId}/server/{gsId}/console",
		events.Handler,
		true,
		false,
	},
	// SSE endpoint for slave
	Route{
		"Host events",
		"GET",
		"/events/{hostToken}",
		events.Handler,
		true,
		true,
	},
	// SSE endpoint for slave
	Route{
		"Game server events",
		"GET",
		"/events/{hostToken}/{gsId}",
		events.Handler,
		true,
		true,
	},
	// internal, not documented in API docs
	Route{
		"Host data",
		"POST",
		"/events/{hostToken}",
		handler.ReceiveData,
		true,
		true,
	},
	// internal, not documented in API docs
	Route{
		"Batch host data",
		"POST",
		"/events/{hostToken}/batch",
		handler.ReceiveBatchData,
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
		"Host Details",
		"GET",
		"/host/{hostId}",
		handler.HostDetails,
		true,
		false,
	},
	Route{
		"Host Edit",
		"PUT",
		"/host/{hostId}",
		handler.HostEdit,
		true,
		false,
	},
	Route{
		"Host Remove",
		"DELETE",
		"/host/{hostId}",
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
		"/host/{hostId}/update",
		handler.HostUpdate,
		true,
		false,
	},
	Route{
		"Host Metric",
		"GET",
		"/host/{hostId}/metric",
		handler.HostMetric,
		true,
		false,
	},
	Route{
		"Host Create Job",
		"POST",
		"/host/{hostId}/job",
		handler.HostCreateJob,
		true,
		false,
	},
	Route{
		"Host Jobs",
		"GET",
		"/host/{hostId}/job",
		handler.HostJobList,
		true,
		false,
	},
	Route{
		"Host Job Remove",
		"DELETE",
		"/host/{hostId}/job/{jobId}",
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
		"/host/{hostId}/server/{gsId}",
		gs.Get,
		true,
		false,
	},
	Route{
		"GameServer Edit",
		"PUT",
		"/host/{hostId}/server/{gsId}",
		gs.Edit,
		true,
		false,
	},
	Route{
		"GameServer Remove",
		"DELETE",
		"/host/{hostId}/server/{gsId}",
		gs.Remove,
		true,
		false,
	},
	Route{
		"GameServer Create",
		"POST",
		"/host/{hostId}/server",
		gs.Create,
		true,
		false,
	},
	Route{
		"GameServer ListByHostId",
		"GET",
		"/host/{hostId}/server",
		gs.ListByHostId,
		true,
		false,
	},
	Route{
		"GameServer Install",
		"PUT",
		"/host/{hostId}/server/{gsId}/install",
		gs.Install,
		true,
		false,
	},
	Route{
		"GameServer Start",
		"PUT",
		"/host/{hostId}/server/{gsId}/start",
		gs.Start,
		true,
		false,
	},
	Route{
		"GameServer Restart",
		"PUT",
		"/host/{hostId}/server/{gsId}/restart",
		gs.Restart,
		true,
		false,
	},
	Route{
		"GameServer Stop",
		"PUT",
		"/host/{hostId}/server/{gsId}/stop",
		gs.Stop,
		true,
		false,
	},
	Route{
		"GameServer Send Command",
		"PUT",
		"/host/{hostId}/server/{gsId}/command",
		gs.SendCommand,
		true,
		false,
	},
	Route{
		"GameServer Console Log",
		"GET",
		"/host/{hostId}/server/{gsId}/logs",
		gs.ConsoleLog,
		true,
		false,
	},
	Route{
		"GameServer Put Logs",
		"PUT",
		"/host/{hostId}/server/{gsId}/logs",
		gs.PutLogs,
		true,
		true,
	},
	Route{
		"GameServer Shutdown",
		"PUT",
		"/host/{hostId}/server/{gsId}/shutdown",
		gs.Shutdown,
		true,
		false,
	},
	Route{
		"GameServer Backup",
		"PUT",
		"/host/{hostId}/server/{gsId}/backup",
		gs.Backup,
		true,
		false,
	},
	Route{
		"GameServer Backup list",
		"GET",
		"/host/{hostId}/server/{gsId}/backup/list",
		gs.BackupList,
		true,
		false,
	},
	Route{
		"GameServer Backup restore",
		"PUT",
		"/host/{hostId}/server/{gsId}/backup/restore",
		gs.BackupRestore,
		true,
		false,
	},

	//GameServer Files
	Route{
		"GameServer File list",
		"GET",
		"/host/{hostId}/server/{gsId}/file/list",
		gs.FileList,
		true,
		false,
	},
	Route{
		"GameServer File remove",
		"DELETE",
		"/host/{hostId}/server/{gsId}/file",
		gs.FileRemove,
		true,
		false,
	},
	Route{
		"GameServer File server",
		"PUT",
		"/host/{hostId}/server/{gsId}/file/server",
		gs.FileServer,
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
