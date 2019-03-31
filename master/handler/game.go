package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/master/db"
	"gitlab.com/systemz/aimpanel2/master/model"
	"net/http"
)

func ListGames(w http.ResponseWriter, r *http.Request) {
	var games []model.Game
	db.DB.Find(&games)
	json.NewEncoder(w).Encode(games)
}
