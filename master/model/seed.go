package model

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func CheckGameInfoInDb(gf GameFile) (needInsert bool, needUpdate bool, id string) {
	count, err := DB.Collection(userCollection).CountDocuments(context.TODO(),
		bson.D{
			{Key: "game_id", Value: gf.GameId},
			{Key: "game_version", Value: gf.GameVersion},
		})
	if err != nil {
		logrus.Error(err)
	}
	// no entry in DB
	if count < 0 {
		return true, false, ""
	}
	// entry present, update just in case
	return false, true, ""
}

func SeedGames() {
	logrus.Info("Adding BungeeCord")
	gameFile := &GameFile{
		GameId:      strconv.Itoa(game.GAME_BUNGEECORD),
		GameVersion: "1420",
		DownloadUrl: "https://ci.md-5.net/job/BungeeCord/1420/artifact/bootstrap/target/BungeeCord.jar",
	}
	err := Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}

	gameFile = &GameFile{
		GameId:      strconv.Itoa(game.GAME_BUNGEECORD),
		GameVersion: "1421",
		DownloadUrl: "https://ci.md-5.net/job/BungeeCord/1421/artifact/bootstrap/target/BungeeCord.jar",
	}
	err = Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Added BungeeCord files successfully.")

	logrus.Info("Adding Spigot")
	gameFile = &GameFile{
		GameId:      strconv.Itoa(game.GAME_SPIGOT),
		GameVersion: "0",
		DownloadUrl: "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar",
	}
	err = Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Added Spigot files successfully.")

	logrus.Info("Adding TeamSpeak3")
	gameFile = &GameFile{
		GameId:      strconv.Itoa(game.GAME_TEAMSPEAK3),
		GameVersion: "3.9.1",
		DownloadUrl: "https://files.teamspeak-services.com/releases/server/3.9.1/teamspeak3-server_linux_amd64-3.9.1.tar.bz2",
	}
	err = Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Added TeamSpeak3 files successfully.")

	logrus.Info("Adding TS3AudioBot")
	gameFile = &GameFile{
		GameId:      strconv.Itoa(game.GAME_TS3AUDIOBOT),
		GameVersion: "master",
		DownloadUrl: "https://splamy.de/api/nightly/ts3ab/master/download",
	}
	err = Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}

	gameFile = &GameFile{
		GameId:      strconv.Itoa(game.GAME_TS3AUDIOBOT),
		GameVersion: "develop",
		DownloadUrl: "https://splamy.de/api/nightly/ts3ab/develop/download",
	}
	err = Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Added TS3Audiobot files successfully.")

	gameFile = &GameFile{
		GameId:      strconv.Itoa(game.GAME_MC_VANILLA),
		GameVersion: "1.16.2",
		DownloadUrl: "https://launcher.mojang.com/v1/objects/c5f6fb23c3876461d46ec380421e42b289789530/server.jar",
	}
	err = Put(gameFile)
	if err != nil {
		logrus.Error(err)
	}
}
