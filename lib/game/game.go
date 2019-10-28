package game

import (
	"github.com/softbrewery/gojoi/pkg/joi"
	"strconv"
	"strings"
)

const (
	GAME_SPIGOT = iota
	GAME_BUNGEECORD

	GAME_TEAMSPEAK3
	GAME_TEAMSPEAK3_BOT
)

type GameDefinition struct {
	Id   uint
	Name string
}

var Games = []GameDefinition{
	{
		Id:   GAME_SPIGOT,
		Name: "Spigot",
	},
	{
		Id:   GAME_BUNGEECORD,
		Name: "BungeeCord",
	},
	{
		Id:   GAME_TEAMSPEAK3,
		Name: "TeamSpeak3",
	},
	{
		Id:   GAME_TEAMSPEAK3_BOT,
		Name: "Teamspeak3 Bot",
	},
}

type Game struct {
	Id uint

	DownloadUrl string

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

	case GAME_TEAMSPEAK3_BOT:
		command = []string{
			"mono",
			"TS3AudioBot.exe",
		}
	}

	return strings.Join(command, " "), nil
}

func (game *Game) GetInstallCmds() (cmds []string, err error) {
	switch game.Id {
	case GAME_SPIGOT:
		cmds = append(cmds, "wget "+game.DownloadUrl+" -O {storageDir}/spigot.jar")
		cmds = append(cmds, "cp {storageDir}/spigot.jar {gsDir}/{uuid}/")
	}

	return cmds, nil
}
