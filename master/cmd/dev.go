package cmd

import (
	"github.com/spf13/cobra"
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

		//model.Snowflake = model.InitSnowflake()
		//model.DB = model.InitDB()
		//
		//taskMsg := task.Message{
		//	TaskId:       task.GAME_COMMAND,
		//	GameServerID: "1256587000778067968",
		//	Body:         "alert test",
		//}
		//
		//hostJob := &model.HostJob{
		//	Base: model.Base{
		//		DocType: "host_job",
		//	},
		//	Name:           "Testowy job",
		//	HostId:         "1256576238273695744",
		//	CronExpression: "* * * * *",
		//	TaskMessage:    taskMsg,
		//}
		//err := hostJob.Put(&hostJob)
		//if err != nil {
		//	logrus.Error(err)
		//}
	},
}
