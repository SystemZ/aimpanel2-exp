package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/model"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [token]",
	Short: "Configure slave to use provided token",
	Long:  "",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		model.SetHostToken(args[0])
		fmt.Fprintf(cmd.OutOrStdout(), "Installed token successfully.")
	},
}
