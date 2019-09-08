package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/master/model"
	"net/http"
)

func ListGames(w http.ResponseWriter, r *http.Request) {
	var games []model.Game
	model.DB.Find(&games)

	lib.MustEncode(json.NewEncoder(w), games)
}

func ListGameVersions(w http.ResponseWriter, r *http.Request) {
	var gameVersions []model.GameVersion
	model.DB.Find(&gameVersions)

	lib.MustEncode(json.NewEncoder(w), gameVersions)
}