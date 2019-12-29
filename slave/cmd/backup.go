package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup <GS ID> <GS ID> ...",
	Short: "Backup game server",
	Long:  "Archive and compress all game server files to one file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		for _, gsId := range args {
			tasks.GsBackupTrigger(gsId)
		}
	},
}
