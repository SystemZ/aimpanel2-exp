package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/tasks/gdrive"
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
		//model.InitRedis()
		service := gdrive.ClientInit()
		gdrive.ListFiles(service)
	},
}
