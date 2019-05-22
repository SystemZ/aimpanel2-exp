package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
	"gitlab.com/systemz/aimpanel2/master/router"
	"net/http"
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

		rabbit.Listen()
		rabbit.ListenWrapperData()
		rabbit.ListenAgentData()

		logrus.Info("Starting API on port :" + args[0])
		r := router.NewRouter()
		logrus.Fatal(http.ListenAndServe(":"+args[0],
			router.CommonMiddleware(router.CorsMiddleware(r))))
	},
}
