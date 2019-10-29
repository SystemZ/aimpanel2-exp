package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/lib/game"
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
		g := game.Game{
			Id:          game.GAME_SPIGOT,
			RamMinM:     1024,
			RamMaxM:     2048,
			JarFilename: "spigot.jar",
			DownloadUrl: "https://cdn.getbukkit.org/spigot/spigot-1.14.4.jar",
		}

		c, err := g.GetCmd()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info(c)

	},
}
