package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/supervisor"
)

func init() {
	rootCmd.AddCommand(supervisorCmd)
}

var supervisorCmd = &cobra.Command{
	Use:   "supervisor",
	Short: "Start supervisor",
	Run: func(cmd *cobra.Command, args []string) {
		supervisor.Start()
	},
}
