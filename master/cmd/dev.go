package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
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
		tree, _ := filemanager.NewTree("/opt/aimpanel/gs/448c29e3-0a2f-40b4-ad83-35fe0847fded")
		logrus.Info(tree.String())
	},
}
