package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/wrapper"
)

func init() {
	rootCmd.AddCommand(wrapperCmd)
}

var wrapperCmd = &cobra.Command{
	Use:   "wrapper [game server id] [game]",
	Short: "Start wrapper",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wrapper.Start(args[0], args[1])
	},
}
