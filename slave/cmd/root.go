package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	quiet bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "No output")
}

var rootCmd = &cobra.Command{
	Use:   "slave",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
