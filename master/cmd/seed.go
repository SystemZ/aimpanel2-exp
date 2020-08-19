package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with mock data",
	Long:  "",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Task seed started.")
		logrus.Info("Connecting to database.")
		model.DB = model.InitDB()

		logrus.Info("Adding game files...")

		model.SeedGames()

		logrus.Info("Added game files successfully.")
		logrus.Info("Task seed finished successfully.")
	},
}
