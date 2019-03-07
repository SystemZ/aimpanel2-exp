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
	HostId uuid.UUID `gorm:"column:user_id" json:"host_id"`

	// Game ID
	//
	// required: true
	GameId uuid.UUID `gorm:"column:ip" json:"ip"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}

func (u *GameServer) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)

	return nil
}
