package lib

type Game struct {
	Name        string
	Command     []string
	FileName    string
	DownloadUrl string
	InstallCmds [][]string
	StopCommand string
}

var GAMES = map[string]Game{
	"minecraft": {
		Name:        "Minecraft",
		Command:     []string{"java", "-Djline.terminal=jline.UnsupportedTerminal", "-jar", "BungeeCord.jar"},
		FileName:    "BungeeCord.jar",
		DownloadUrl: "https://ci.md-5.net/job/BungeeCord/lastSuccessfulBuild/artifact/bootstrap/target/BungeeCord.jar",
		InstallCmds: [][]string{
			{"cp", "/opt/aimpanel/storage/{fileName}", "/opt/aimpanel/gs/{uuid}/"},
		},
		StopCommand: "exit",
	},
	"teamspeak3": {
		Name:        "TeamSpeak3",
		Command:     []string{"sh", "teamspeak3-server_linux_amd64/ts3server_minimal_runscript.sh"},
		FileName:    "teamspeak3-server_linux_amd64-3.5.0.tar.bz2",
		DownloadUrl: "http://dl.4players.de/ts/releases/3.5.0/teamspeak3-server_linux_amd64-3.5.0.tar.bz2",
		InstallCmds: [][]string{},
		StopCommand: "",
	},
}
