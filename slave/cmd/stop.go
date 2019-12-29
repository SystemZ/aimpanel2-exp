package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop <GS ID> <GS ID> ...",
	Short: "Stop game server",
	Long:  "Gracefully stop 1 or more game servers\nUse \"all\" to stop all GS on this host",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		for _, gsId := range args {
			tasks.GsStop(gsId)
		}
	},
}
