package cmd

import (
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/exit"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/router"
	"net/http"
	"os"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [port]",
	Short: "Start master server",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exit.CheckForExitSignal()

		model.DB = model.InitDB()
		// helper for sending slave tasks via SSE
		model.EmitterInit()
		//go model.EventChanges()

		//model.InitRedis()
		events.SSE = events.InitSSE()

		router.InitRateLimit()

		//go cron.SseHeartbeat()
		//go cron.CheckHostsHeartbeat()
		//go cron.CheckGSHeartbeat()

		logrus.Info("Starting API on port :" + args[0])
		r := router.NewRouter()

		// enable CORS only in dev mode
		if config.DEV_MODE {
			logrus.Fatal(http.ListenAndServeTLS(
				":"+args[0],
				config.HTTP_TLS_CERT_PATH,
				config.HTTP_TLS_KEY_PATH,
				router.Limiter.RateLimit(
					router.CorsMiddleware(
						handlers.LoggingHandler(os.Stdout, r),
					),
				),
			))
		} else {
			logrus.Fatal(http.ListenAndServe(
				":"+args[0],
				router.Limiter.RateLimit(
					handlers.LoggingHandler(os.Stdout, r)),
			))
		}
	},
}
