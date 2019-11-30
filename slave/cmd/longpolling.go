package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(longpollingCmd)
}

var longpollingCmd = &cobra.Command{
	Use:   "longpolling",
	Short: "For testing longpolling",
	Long:  "",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//lib.GetGameCmd()
	},
}
