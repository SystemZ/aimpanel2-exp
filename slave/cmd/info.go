package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show info about installed version of slave",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", config.GIT_COMMIT)
		fmt.Fprintf(cmd.OutOrStdout(), "Token: %s\n", model.GetHostToken())
		fmt.Fprintf(cmd.OutOrStdout(), "HW ID: %s\n", model.GetHwId())
		fmt.Fprintf(cmd.OutOrStdout(), "Instance URL: %s\n", config.MASTER_URLS[0])
		fmt.Fprintf(cmd.OutOrStdout(), "\nReinstall CMD:\n")
		fmt.Fprintf(cmd.OutOrStdout(), "wget https://%s/i/%s -O- | bash -\n",
			config.MASTER_URLS[0], model.GetHostToken())
	},
}
