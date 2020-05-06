package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"gitlab.com/systemz/aimpanel2/master/model"
	"strconv"
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
		logrus.Info("Task seed started.")
		logrus.Info("Connecting to database.")
		model.DB = model.InitDB()

		logrus.Info("Adding game files...")

		logrus.Info("Adding BungeeCord")
		gameFile := &model.GameFile{
			GameId:      strconv.Itoa(game.GAME_BUNGEECORD),
			GameVersion: "1420",
			DownloadUrl: "https://ci.md-5.net/job/BungeeCord/1420/artifact/bootstrap/target/BungeeCord.jar",
		}
		err := model.Put(gameFile)
		if err != nil {
			logrus.Error(err)
		}

		gameFile = &model.GameFile{
			GameId:      strconv.Itoa(game.GAME_BUNGEECORD),
			GameVersion: "1421",
			DownloadUrl: "https://ci.md-5.net/job/BungeeCord/1421/artifact/bootstrap/target/BungeeCord.jar",
		}
		err = model.Put(gameFile)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("Added BungeeCord files successfully.")

		logrus.Info("Adding Spigot")
		gameFile = &model.GameFile{
			GameId:      strconv.Itoa(game.GAME_SPIGOT),
			GameVersion: "0",
			DownloadUrl: "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
		}
		err = model.Put(gameFile)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("Added Spigot files successfully.")

		logrus.Info("Adding TeamSpeak3")
		gameFile = &model.GameFile{
			GameId:      strconv.Itoa(game.GAME_TEAMSPEAK3),
			GameVersion: "3.9.1",
			DownloadUrl: "https://files.teamspeak-services.com/releases/server/3.9.1/teamspeak3-server_linux_amd64-3.9.1.tar.bz2",
		}
		err = model.Put(gameFile)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("Added TeamSpeak3 files successfully.")

		logrus.Info("Adding TS3AudioBot")
		gameFile = &model.GameFile{
			GameId:      strconv.Itoa(game.GAME_TS3AUDIOBOT),
			GameVersion: "master",
			DownloadUrl: "https://splamy.de/api/nightly/ts3ab/master/download",
		}
		err = model.Put(gameFile)
		if err != nil {
			logrus.Error(err)
		}

		gameFile = &model.GameFile{
			GameId:      strconv.Itoa(game.GAME_TS3AUDIOBOT),
			GameVersion: "develop",
			DownloadUrl: "https://splamy.de/api/nightly/ts3ab/develop/download",
		}
		err = model.Put(gameFile)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("Added TS3Audiobot files successfully.")

		logrus.Info("Added game files successfully.")

		logrus.Info("Task seed finished successfully.")

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
