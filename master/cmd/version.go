package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/config"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "",
	Long:  "",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.GIT_COMMIT)
	},
}
