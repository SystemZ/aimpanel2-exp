package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func init() {
	rootCmd.AddCommand(killCmd)
}

var killCmd = &cobra.Command{
	Use:   "kill <GS ID> <GS ID> ...",
	Short: "Kill game server",
	Long:  "Instantly and forcefully hard stop 1 or more game servers\nUse \"all\" to kill all GS on this host",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		for _, gsId := range args {
			tasks.GsKill(gsId)
		}
	},
}
