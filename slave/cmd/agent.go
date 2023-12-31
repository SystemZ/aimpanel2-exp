package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/agent"
)

func init() {
	rootCmd.AddCommand(agentCmd)
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start agent",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		agent.Start()
	},
}
