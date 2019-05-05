package model

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// Game Server represents the game server of the host
// swagger:model game_server
type GameServer struct {
	// ID of the game server
	//
	// required: true
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)" json:"id"`

	// Name
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// Host ID
	//
	// required: true
	HostId uuid.UUID `gorm:"column:host_id" json:"host_id"`

	// State
	// 0 off, 1 running
	// required: false
	State uint `gorm:"column:state" json:"state"`

	// State Last Changed
	//
	// required: false
	StateLastChanged time.Time `gorm:"default:CURRENT_TIMESTAMP column:state_last_changed" json:"state_last_changed"`

	// Game ID
	//
	// required: true
	GameId uint `gorm:"column:game_id" json:"game_id"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}

func (gs *GameServer) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)

	return nil
}

func (gs *GameServer) GetGame(db *gorm.DB) *Game {
	var game Game

	if db.Where("id = ?", gs.GameId).First(&game).RecordNotFound() {
		return nil
	}

	return &game
}
