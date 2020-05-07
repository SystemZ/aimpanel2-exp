package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func (g *GameFile) GetID() primitive.ObjectID {
	return g.ID
}

func GetGameFileByGameIdAndVersion(gameId uint, version string) (*GameFile, error) {
	var gf GameFile

	err := DB.Collection(gameFileCollection).FindOne(context.TODO(), bson.D{
		{"game_id", gameId},
		{"$or", []interface{}{
			bson.D{{"game_version", version}},
			bson.D{{"game_version", "0"}},
		}},
	}).Decode(&gf)
	if err != nil {
		return nil, err
	}

	return &gf, nil
}
