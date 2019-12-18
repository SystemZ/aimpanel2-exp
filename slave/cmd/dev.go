package cmd

import (
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
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
		//lib.GetGameCmd()
		client := sse.NewClient("http://localhost:9000/v1/events/x")
		client.Headers = map[string]string{
			"Authorization": "Bearer x",
		}
		err := client.SubscribeRaw(func(msg *sse.Event) {
			logrus.Info(string(msg.ID))
			logrus.Info(string(msg.Data))
		})
		if err != nil {
			logrus.Info(err)
		}
	},
}
