package game

import "github.com/softbrewery/gojoi/pkg/joi"

type JavaGame struct {
	Game
	RamMinM int
	RamMaxM int
}

func (game JavaGame) getCmd() string {
	return "java"
}

func (game JavaGame) validate() (err error) {
	err = joi.Validate(game.RamMinM, joi.Int().Min(16))
	if err != nil {
		return err
	}
	return joi.Validate(game.RamMaxM, joi.Int().Min(64))
}
