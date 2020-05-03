package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/lib/task"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "For testing on dev various things",
	Long:  "",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//ahttp.HttpClient = ahttp.InitHttpClient()
		//
		//_, err := ahttp.Get("https://aimpanel.local/", nil)
		//if err != nil {
		//	logrus.Error(err.Error())
		//}
		//
		//_, err = ahttp.Get("https://google.com/", nil)
		//if err != nil {
		//	logrus.Error(err.Error())
		//}
		//
		//_, err = ahttp.Get("https://api-lab.aimpanel.pro", nil)
		//if err != nil {
		//	logrus.Error(err.Error())
		//}

		logrus.Info(task.AGENT_BACKUP_GS)
		logrus.Info(task.AGENT_BACKUP_GS.StringValue())
	},
}
