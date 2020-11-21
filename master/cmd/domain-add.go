package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func init() {
	rootCmd.AddCommand(domainAddCmd)
}

var domainAddCmd = &cobra.Command{
	Use:     "domain-add <domain> <domain>",
	Example: "domain-add aimpanel.pro",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.DB = model.InitDB()
		for _, arg := range args {
			logrus.Info(arg)
			domain, err := model.GetCertDomainByName(arg)
			if domain != nil || err != nil {
				logrus.Warningf("domain %v found, skipping...", arg)
				continue
			}

			newDomain := &model.CertDomain{
				Name:     arg,
				PoolSize: 1,
			}

			err = model.Put(newDomain)
			if err != nil {
				logrus.Errorf("put failed for domain %v", arg)
			}

			logrus.Infof("added %v domain to database", arg)
		}
	},
}
