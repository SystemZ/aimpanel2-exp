package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameFile struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"1238206236281802752"`

	GameId string `json:"game_id"`

	GameVersion string `json:"game_version"`

	DownloadUrl string `json:"download_url"`
}

func (g *GameFile) GetCollectionName() string {
	return "games_file"
}

func GetGameFileByGameIdAndVersion(gameId uint, version string) *GameFile {
	var gf GameFile

	err := GetOneS(&gf, map[string]interface{}{
		"doc_type": "game_file",
		"game_id":  fmt.Sprint(gameId),
		"$or": []map[string]interface{}{
			{
				"game_version": version,
			},
			{
				"game_version": "0",
			},
		},
	})
	if err != nil {
		return nil
	}

	return &gf
}
