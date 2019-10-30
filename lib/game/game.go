package game

import (
	"github.com/sirupsen/logrus"
	"github.com/softbrewery/gojoi/pkg/joi"
	"gitlab.com/systemz/aimpanel2/lib"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	GAME_SPIGOT = iota + 1
	GAME_BUNGEECORD

	GAME_TEAMSPEAK3
	GAME_TEAMSPEAK3_BOT
)

type GameDefinition struct {
	Id       uint     `json:"id"`
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

var Games = []GameDefinition{
	{
		Id:   GAME_SPIGOT,
		Name: "Spigot",
		Versions: []string{
			"1.14.4",
			"1.14.3",
			"1.14.2",
			"1.14.1",
			"1.14",
			"1.13.2",
			"1.13.1",
			"1.13"},
	},
	{
		Id:   GAME_BUNGEECORD,
		Name: "BungeeCord",
		Versions: []string{
			"1421", "1420",
		},
	},
	{
		Id:   GAME_TEAMSPEAK3,
		Name: "TeamSpeak3",
		Versions: []string{
			"1.0.0",
		},
	},
	{
		Id:   GAME_TEAMSPEAK3_BOT,
		Name: "Teamspeak3 Bot",
		Versions: []string{
			"1.0.0",
		},
	},
}

type Game struct {
	Id uint `json:"id,omitempty"`

	Version     string `json:"version,omitempty"`
	DownloadUrl string `json:"download_url,omitempty"`

	RamMinM     int    `json:"ram_min_m,omitempty"`
	RamMaxM     int    `json:"ram_max_m,omitempty"`
	JarFilename string `json:"jar_filename,omitempty"`
}

func (game *Game) SetDefaults() {
	switch game.Id {
	case GAME_SPIGOT:
		game.RamMinM = 1024
		game.RamMaxM = 2048
		game.JarFilename = "spigot.jar"
	case GAME_BUNGEECORD:
		game.RamMinM = 1024
		game.RamMaxM = 2048
		game.JarFilename = "BungeeCord.jar"
	}
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

func (game *Game) Install(storagePath string, gsPath string) (err error) {
	switch game.Id {
	case GAME_SPIGOT:
		if _, err := os.Stat(storagePath + "/BuildTools/"); os.IsNotExist(err) {
			_ = os.Mkdir(storagePath+"/BuildTools/", 0777)
		}

		if _, err := os.Stat(storagePath + "/BuildTools/BuildTools.jar"); os.IsNotExist(err) {
			err = lib.DownloadFile(game.DownloadUrl, storagePath+"/BuildTools/BuildTools.jar")
			if err != nil {
				return err
			}
		}

		if _, err := os.Stat(storagePath + "/spigot-" + game.Version + ".jar"); os.IsNotExist(err) {
			cmd := exec.Command("java", "-jar", "BuildTools.jar", "--output-dir="+storagePath, " --rev="+game.Version)
			cmd.Dir = storagePath + "/BuildTools"

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err = cmd.Run(); err != nil {
				logrus.Error(err)
			}

			cmd.Wait()
		}

		err = lib.CopyFile(storagePath+"/spigot-"+game.Version+".jar", gsPath+"/"+game.JarFilename)
		if err != nil {
			return err
		}

	case GAME_BUNGEECORD:
		if _, err := os.Stat(storagePath + "/BungeeCord-" + game.Version + ".jar"); os.IsNotExist(err) {
			err = lib.DownloadFile(game.DownloadUrl, storagePath+"/BungeeCord-"+game.Version+".jar")
			if err != nil {
				return err
			}
		}

		err = lib.CopyFile(storagePath+"/BungeeCord-"+game.Version+".jar", gsPath+"/"+game.JarFilename)
		if err != nil {
			return err
		}
	}
	return nil
}
