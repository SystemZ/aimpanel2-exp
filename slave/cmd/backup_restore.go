package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func init() {
	rootCmd.AddCommand(backupRestoreCmd)
}

var backupRestoreCmd = &cobra.Command{
	Use:   "backup-restore <GS ID> <BACKUP FILENAME>",
	Short: "Restore backup of game server",
	Long:  "Restore specific game server backup to gs dir",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		tasks.GsBackupRestoreTrigger(args[0], args[1])
	},
}
