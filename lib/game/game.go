package game

import (
	"github.com/softbrewery/gojoi/pkg/joi"
	"strconv"
	"strings"
)

const (
	GAME_SPIGOT     = iota
	GAME_TEAMSPEAK3 = iota
	GAME_BUNGEECORD = iota
)

type Game struct {
	Id int

	DeveloperName string
	DeveloperUrl  string

	RamMinM     int
	RamMaxM     int
	JarFilename string
}

func (game *Game) Validate() (err error) {
	switch game.Id {
	case GAME_SPIGOT:
		err = joi.Validate(game.RamMinM, joi.Int().Min(16))
		if err != nil {
			return err
		}

		err = joi.Validate(game.RamMaxM, joi.Int().Min(64))
		if err != nil {
			return err
		}

		return joi.Validate(game.JarFilename, joi.String().Min(3))
	case GAME_BUNGEECORD:
		err = joi.Validate(game.RamMinM, joi.Int().Min(16))
		if err != nil {
			return err
		}

		err = joi.Validate(game.RamMaxM, joi.Int().Min(64))
		if err != nil {
			return err
		}

		return joi.Validate(game.JarFilename, joi.String().Min(3))
	default:
		return nil
	}
}

func (game *Game) GetCmd() (cmd string, err error) {
	err = game.Validate()
	if err != nil {
		return "", err
	}

	var command []string
	switch game.Id {
	case GAME_SPIGOT:
		command = []string{
			"java",
			"-Xms" + strconv.Itoa(game.RamMinM) + "M",
			"-Xmx" + strconv.Itoa(game.RamMaxM) + "M",
			"-Djline.terminal=jline.UnsupportedTerminal",
			"-jar",
			game.JarFilename,
		}

	case GAME_BUNGEECORD:
		command = []string{
			"java",
			"-Xms" + strconv.Itoa(game.RamMinM) + "M",
			"-Xmx" + strconv.Itoa(game.RamMaxM) + "M",
			"-Djline.terminal=jline.UnsupportedTerminal",
			"-jar",
			game.JarFilename,
		}

	case GAME_TEAMSPEAK3:
		command = []string{
			"sh",
			"ts3server_minimal_runscript.sh",
		}
	}

	return strings.Join(command, " "), nil
}
