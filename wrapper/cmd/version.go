package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var appVersion string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Version of this build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(appVersion)
	},
}
