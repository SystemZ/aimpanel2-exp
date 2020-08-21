package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
)

func init() {
	rootCmd.AddCommand(adminSetCmd)
	adminSetCmd.Flags().BoolVar(&adminSetCmdAdminStatus, "admin", false, "Set admin permissions")
}

var (
	adminSetCmdAdminStatus bool
)

var adminSetCmd = &cobra.Command{
	Use:     "admin-set --admin <username>",
	Example: "admin-set --admin joe\nadmin-set joe",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.DB = model.InitDB()
		for _, arg := range args {
			logrus.Info(arg)
			user, err := model.GetUserByUsername(arg)
			if err != nil {
				logrus.Warningf("user %v not found, skipping...", arg)
				continue
			}
			logrus.Infof("before update: admin=%v", user.Admin)
			user.Admin = adminSetCmdAdminStatus
			err = model.Update(user)
			if err != nil {
				logrus.Errorf("update failed for user %v", arg)
			}
			logrus.Infof("after update: admin=%v", user.Admin)
		}
	},
}
