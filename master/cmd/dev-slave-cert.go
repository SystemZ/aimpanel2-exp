package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/cert"
)

func init() {
	rootCmd.AddCommand(devSlaveCert)
}

var devSlaveCert = &cobra.Command{
	Use:     "dev-slave-cert [ip address or domain]",
	Example: "domain-add 192.168.2.2",
	Short:   "Generate self signed cert for slave",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.DB = model.InitDB()
		for _, arg := range args {
			domain, err := model.GetCertDomainByName(arg)

			newDomain := &model.CertDomain{
				Name:     arg,
				PoolSize: 1,
			}

			if domain == nil {
				err = model.Put(newDomain)
				if err != nil {
					logrus.Errorf("put failed for domain %v", arg)
				}

				domain = newDomain

				logrus.Infof("added %v domain to database", arg)
			}

			err = cert.CreateSelfSignedCert(*domain)
			if err != nil {
				logrus.Errorf("create self signed cert failed %v", err)
			}

			logrus.Infof("created self signed cert for %v", arg)
		}
	},
}
