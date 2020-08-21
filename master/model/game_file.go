package model

import (
	"context"
	"fmt"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type GameFile struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"1238206236281802752"`

	GameId string `bson:"game_id" json:"game_id"`

	GameVersion string `bson:"game_version" json:"game_version"`

	DownloadUrl string `bson:"download_url" json:"download_url"`
}

func (g *GameFile) GetCollectionName() string {
	return GameFileCollection
}

func (g *GameFile) GetID() primitive.ObjectID {
	return g.ID
}

func (g *GameFile) SetID(id primitive.ObjectID) {
	g.ID = id
}

func GetGameFileByGameIdAndVersionFromMongo(gameId uint, version string) (*GameFile, error) {
	var gf GameFile

	err := DB.Collection(GameFileCollection).FindOne(context.TODO(), bson.D{
		{Key: "game_id", Value: fmt.Sprintf("%v", gameId)},
		{Key: "$or", Value: []bson.D{
			bson.D{{Key: "game_version", Value: version}},
			bson.D{{Key: "game_version", Value: "0"}},
		}},
	}).Decode(&gf)
	if err != nil {
		return nil, err
	}

	return &gf, nil
}

func GetGameFileByGameIdAndVersion(gameId uint, version string) (gameFile *GameFile, err error) {
	downloadUrl := ""
	switch gameId {
	case game.GAME_BUNGEECORD:
		downloadUrl = fmt.Sprintf("https://ci.md-5.net/job/BungeeCord/%s/artifact/bootstrap/target/BungeeCord.jar", version)
	case game.GAME_SPIGOT:
		downloadUrl = "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
	case game.GAME_TEAMSPEAK3:
		downloadUrl = fmt.Sprintf("https://files.teamspeak-services.com/releases/server/%s/teamspeak3-server_linux_amd64-%s.tar.bz2", version, version)
	case game.GAME_TS3AUDIOBOT:
		downloadUrl = fmt.Sprintf("https://splamy.de/api/nightly/ts3ab/%s/download", version)
	case game.GAME_MC_VANILLA:
		switch version {
		case "1.16.2":
			downloadUrl = "https://launcher.mojang.com/v1/objects/c5f6fb23c3876461d46ec380421e42b289789530/server.jar"
		}
	// bedrock
	// https://www.minecraft.net/en-us/download/server/bedrock
	case game.GAME_FACTORIO:
		downloadUrl = fmt.Sprintf("https://www.factorio.com/get-download/%s/headless/linux64", version)
	}

	gameFile = &GameFile{
		GameId:      strconv.Itoa(int(gameId)),
		GameVersion: version,
		DownloadUrl: downloadUrl,
	}
	return
}
