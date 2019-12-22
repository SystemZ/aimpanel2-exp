package cmd

import (
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/cron"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
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
		model.DB = model.InitMysql()
		model.InitRedis()

		events.SSE = events.InitSSE()

		rabbit.Listen()
		rabbit.SetupRabbitAPI()

		go cron.CheckHostsHeartbeat()
		go cron.CheckGSHeartbeat()

		logrus.Info("Starting API on port :" + args[0])
		r := router.NewRouter()

		logrus.Fatal(http.ListenAndServe(
			":"+args[0],
			router.CommonMiddleware(
				router.CorsMiddleware(
					handlers.LoggingHandler(os.Stdout, r),
				),
			),
		))
	},
}
