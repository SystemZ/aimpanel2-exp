package cmd

import (
	"github.com/michaelklishin/rabbit-hole"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/rabbit"
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
		rabbit.SetupRabbitHole()

		xs, err := rabbit.Client.ListUsers()
		if err != nil {
			logrus.Error(err)
		}

		logrus.Info(xs)

		resp, err := rabbit.Client.PutUser("test", rabbithole.UserSettings{Password: "test", Tags: "administrator"})
		if err != nil {
			logrus.Error(err)
		}

		logrus.Info(resp.Status)

		xs, err = rabbit.Client.ListUsers()
		if err != nil {
			logrus.Error(err)
		}

		resp, err = rabbit.Client.UpdatePermissionsIn("/", "test", rabbithole.Permissions{
			Configure: ".*",
			Write:     ".*",
			Read:      ".*",
		})

		logrus.Info(xs)
	},
}
