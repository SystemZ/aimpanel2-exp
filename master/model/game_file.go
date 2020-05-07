package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameFile struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id" example:"1238206236281802752"`

	GameId string `bson:"game_id" json:"game_id"`

	GameVersion string `bson:"game_version" json:"game_version"`

	DownloadUrl string `bson:"download_url" json:"download_url"`
}

func (g *GameFile) GetCollectionName() string {
	return gameFileCollection
}

func (g *GameFile) GetID() primitive.ObjectID {
	return g.ID
}

func (g *GameFile) SetID(id primitive.ObjectID) {
	g.ID = id
}

func GetGameFileByGameIdAndVersion(gameId uint, version string) (*GameFile, error) {
	var gf GameFile

	err := DB.Collection(gameFileCollection).FindOne(context.TODO(), bson.D{
		{"game_id", fmt.Sprintf("%v", gameId)},
		{"$or", []bson.D{
			bson.D{{"game_version", version}},
			bson.D{{"game_version", "0"}},
		}},
	}).Decode(&gf)
	if err != nil {
		return nil, err
	}

	return &gf, nil
}
