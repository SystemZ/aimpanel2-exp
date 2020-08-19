package game

import (
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
	GAME_TS3AUDIOBOT
	GAME_MC_VANILLA
	GAME_FACTORIO
)

type GameDefinition struct {
	Id          uint     `json:"id"`
	Name        string   `json:"name"`
	Versions    []string `json:"versions"`
	StopCommand string   `json:"stop_command"`
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
			"3.9.1",
		},
	},
	{
		Id:   GAME_TS3AUDIOBOT,
		Name: "TS3AudioBot",
		Versions: []string{
			"master", "develop",
		},
	},
	{
		Id:   GAME_MC_VANILLA,
		Name: "Minecraft Java Edition Vanilla",
		Versions: []string{
			"1.16.2",
		},
	},
	{
		Id:   GAME_FACTORIO,
		Name: "Factorio",
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

	StopCommand     string `json:"stop_command,omitempty"`
	StopTimeout     int    `json:"stop_timeout,omitempty"`
	StopHardTimeout int    `json:"stop_hard_timeout,omitempty"`
}

func (game *Game) SetDefaults() {
	switch game.Id {
	case GAME_SPIGOT:
		game.RamMinM = 1024
		game.RamMaxM = 2048
		game.JarFilename = "spigot.jar"
		game.StopCommand = "stop"
		game.StopTimeout = 15
		game.StopHardTimeout = 30
	case GAME_BUNGEECORD:
		game.RamMinM = 1024
		game.RamMaxM = 2048
		game.JarFilename = "BungeeCord.jar"
		game.StopCommand = "end"
		game.StopTimeout = 15
		game.StopHardTimeout = 30
	case GAME_MC_VANILLA:
		game.RamMinM = 1024
		game.RamMaxM = 2048
		game.JarFilename = "minecraft_server.jar"
		game.StopCommand = "stop"
		game.StopTimeout = 15
		game.StopHardTimeout = 30
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
	case GAME_MC_VANILLA:
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
	case GAME_MC_VANILLA:
		command = []string{
			"java",
			"-Xms" + strconv.Itoa(game.RamMinM) + "M",
			"-Xmx" + strconv.Itoa(game.RamMaxM) + "M",
			"-Djline.terminal=jline.UnsupportedTerminal",
			"-jar",
			game.JarFilename,
			"nogui",
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
	case GAME_TS3AUDIOBOT:
		command = []string{
			"mono",
			"TS3AudioBot.exe",
			"--non-interactive",
		}
	case GAME_FACTORIO:
		command = []string{
			"./bin/x64/factorio",
			"--start-server-load-latest",
			//"--start-server",
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
			cmd := exec.Command("java", "-jar", "BuildTools.jar", "--output-dir", storagePath, "--rev", game.Version)
			cmd.Dir = storagePath + "/BuildTools"

			if err = cmd.Run(); err != nil {
				return err
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
	case GAME_MC_VANILLA:
		if _, err := os.Stat(storagePath + "/minecraft_server." + game.Version + ".jar"); os.IsNotExist(err) {
			err = lib.DownloadFile(game.DownloadUrl, storagePath+"/minecraft_server."+game.Version+".jar")
			if err != nil {
				return err
			}
		}

		err = lib.CopyFile(storagePath+"/minecraft_server."+game.Version+".jar", gsPath+"/"+game.JarFilename)
		if err != nil {
			return err
		}
	case GAME_TEAMSPEAK3:
		storageFilePath := storagePath + "/teamspeak3-server-" + game.Version + ".tar.bz2"
		if _, err := os.Stat(storageFilePath); os.IsNotExist(err) {
			err = lib.DownloadFile(game.DownloadUrl, storageFilePath)
		}

		cmd := exec.Command("tar", "xfvj", storageFilePath, "--strip-components=1", "--directory="+gsPath)
		if err = cmd.Run(); err != nil {
			return err
		}
		cmd.Wait()

		licenseFile, err := os.Create(gsPath + "/.ts3server_license_accepted")
		if err != nil {
			return err
		}
		licenseFile.Close()
	case GAME_TS3AUDIOBOT:
		storageFilePath := storagePath + "/ts3audiobot-" + game.Version + ".zip"
		err = lib.DownloadFile(game.DownloadUrl, storageFilePath)

		cmd := exec.Command("unzip", storageFilePath, "-d", gsPath)
		if err = cmd.Run(); err != nil {
			return err
		}
		cmd.Wait()
	case GAME_FACTORIO:
		storageFilePath := storagePath + "/factorio_headless_x64_" + game.Version + ".tar.xz"
		if _, err := os.Stat(storageFilePath); os.IsNotExist(err) {
			err = lib.DownloadFile(game.DownloadUrl, storageFilePath)
		}

		cmd := exec.Command("tar", "xfv", storageFilePath, "--strip-components=1", "--directory="+gsPath)
		if err = cmd.Run(); err != nil {
			return err
		}
		cmd.Wait()
	}
	return nil
}
