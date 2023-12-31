package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(backupListCmd)
}

var backupListCmd = &cobra.Command{
	Use:   "backup-list",
	Short: "List available backups",
	Long:  "Show all backups that can be restored",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("WIP, it's not working yet")
		//model.InitRedis()
		//service := gdrive.ClientInit()
		//gdrive.ListFiles(service)
	},
}
