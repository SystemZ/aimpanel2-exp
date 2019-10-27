package lib

//
//type Game struct {
//	DeveloperName string
//	DeveloperUrl  string
//}
//
//type JavaGame struct {
//	Game
//	RamMinM int
//	RamMaxM int
//}
//
//func (game JavaGame) getCmd() string {
//	return "java"
//}
//
//func (game JavaGame) validate() (err error) {
//	err = joi.Validate(game.RamMinM, joi.Int().Min(16))
//	if err != nil {
//		return err
//	}
//	return joi.Validate(game.RamMaxM, joi.Int().Min(64))
//}
//
//type Minecraft struct {
//	JavaGame
//	JarFilename string
//	MaxPlayers  uint
//}
//
//func (game Minecraft) validate() (err error) {
//	err = game.JavaGame.validate()
//	if err != nil {
//		return err
//	}
//	return joi.Validate(game.JarFilename, joi.String().Min(3))
//}
//
//type Spigot struct {
//	Minecraft
//	MobsDisabled bool //just an example
//}
//
//func (game Spigot) validate() (err error) {
//	err = game.Minecraft.validate()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (game Spigot) getCmd() (err error, cmd string) {
//	err = game.validate()
//	if err != nil {
//		return err, ""
//	}
//	cmd = "java -Xmx" + strconv.Itoa(int(game.RamMaxM)) + " -Djline.terminal=jline.UnsupportedTerminal -jar " + game.JarFilename
//	return nil, cmd
//}

//func GetGameCmd() {
//	var game Spigot
//	err, cmd := game.getCmd()
//	if err != nil {
//		log.Printf("%v", err)
//		return
//	}
//	log.Printf("CMD: %v", cmd)
//}
