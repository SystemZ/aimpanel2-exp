package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with mock data",
	Long:  "",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		//model.DB = model.InitMysql()
		//
		//var count int
		//var users []model.User
		//model.DB.Find(&users).Count(&count)
		//if count == 0 {
		//	file, err := os.Open("./dump/seed.sql")
		//	if err != nil {
		//		logrus.Fatal(err)
		//	}
		//	defer file.Close()
		//
		//	reader := bufio.NewReader(file)
		//	var line string
		//	for {
		//		line, err = reader.ReadString('\n')
		//
		//		if len(line) > 1 {
		//			model.DB.Exec(line)
		//		}
		//
		//		if err != nil {
		//			break
		//		}
		//	}
		//}
	},
}
