package model

type GameFile struct {
	Base

	GameId string `json:"game_id"`

	GameVersion string `json:"game_version"`

	DownloadUrl string `json:"download_url"`
}

func GetGameFileByGameIdAndVersion(gameId uint, version string) *GameFile {
	var gf GameFile

	err := GetOneS(&gf, map[string]interface{}{
		"doc_type": "game_file",
		"game_id":  gameId,
		"$or": []map[string]interface{}{
			{
				"game_version": version,
			},
			{
				"game_version": "0",
			},
		},
	})
	if err != nil {
		return nil
	}

	return &gf
}
