package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/cert"
)

func init() {
	rootCmd.AddCommand(generateCertsCmd)
}

var generateCertsCmd = &cobra.Command{
	Use:   "generate-certs",
	Short: "Generate certs for all domains in database",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		model.DB = model.InitDB()

		err := cert.InitLego()
		if err != nil {
			logrus.Fatal(err)
		}

		err = cert.CreateCerts()
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Info("Certs generated successfully.")
	},
}
