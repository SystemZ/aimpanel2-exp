package model

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type GameServer struct {
	// ID of the game server
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)" json:"id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// User assigned name
	Name string `gorm:"column:name" json:"name" example:"Ultra MC Server"`

	// Host ID
	//
	// required: true
	HostId uuid.UUID `gorm:"column:host_id" json:"host_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// State
	// 0 off, 1 running
	State uint `gorm:"column:state" json:"state" example:"0"`

	// State Last Changed
	StateLastChanged time.Time `gorm:"default:CURRENT_TIMESTAMP;column:state_last_changed" json:"state_last_changed" example:"2019-09-29T03:16:27+02:00"`

	// Game ID
	GameId uint `gorm:"column:game_id" json:"game_id" example:"1"`

	// Game Version
	GameVersion string `gorm:"column:game_version" json:"game_version" example:"1.14.2"`

	// Game
	GameJson string `gorm:"column:game_json" json:"-"`

	// Metric Frequency
	MetricFrequency int `gorm:"column:metric_frequency" json:"metric_frequency" example:"30"`

	// Stop Timeout
	StopTimeout int `gorm:"column:stop_timeout" json:"stop_timeout" example:"30"`

	// Date of game server creation
	CreatedAt time.Time `json:"created_at" example:"2019-09-29T03:16:27+02:00"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"-"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"-"`
}

func (gs *GameServer) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)

	return nil
}

func GetGameServer(db *gorm.DB, gsId string) *GameServer {
	var gs GameServer

	if db.Where("id = ?", gsId).First(&gs).RecordNotFound() {
		return nil
	}

	return &gs
}

func GetGameServersByHostId(db *gorm.DB, hostId string) *[]GameServer {
	var gs []GameServer

	if db.Where("host_id = ?", hostId).Find(&gs).RecordNotFound() {
		return nil
	}

	return &gs
}
