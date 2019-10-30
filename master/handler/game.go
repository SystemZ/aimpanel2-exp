package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"net/http"
)

func ListGames(w http.ResponseWriter, r *http.Request) {
	lib.MustEncode(json.NewEncoder(w), game.Games)
}
