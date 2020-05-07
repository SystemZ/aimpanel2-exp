package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		d := []model.GameFile{}
		oid, err := primitive.ObjectIDFromHex("5eb32724646b16777be3b394")
		if err != nil {
			logrus.Error(err)
			return
		}

		err = model.Find(&d, bson.D{{"_id", oid}})
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info(d)

		//d.DownloadUrl = "lalalalala2"
		//err = model.Update(&d)
		//if err != nil {
		//	logrus.Error(err.Error())
		//	return
		//}
	},
}
