package handler

import (
	"encoding/json"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/game"
	"net/http"
)

// @Summary List
// @Tags Game
// @Description List supported games
// @Accept json
// @Produce json
// @Success 200 {array} response.Game
// @Failure 400 {object} JsonError
// @Router /game [get]
func ListGames(w http.ResponseWriter, r *http.Request) {
	lib.MustEncode(json.NewEncoder(w), game.Games)
}
