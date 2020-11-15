package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/cert"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "For testing on dev various things",
	Long:  "",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.DB = model.InitDB()
		cert.InitLego()

		cert.CreateCertForDomain("ksawicki.pl")
	},
}
