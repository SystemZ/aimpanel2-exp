package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/lib/metric"
	"gitlab.com/systemz/aimpanel2/master/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
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

		oid, _ := primitive.ObjectIDFromHex("5eb6039ec00d8680d63dc114")
		now := time.Now()
		from := time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		to := time.Date(now.Year()+1, now.Month(), now.Day(), 23, 59, 0, 0, now.Location())

		res, _ := model.GetTimeSeries(oid, from, to, metric.RamFree)
		log.Printf("%v", res)
	},
}
