package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// Host represents the host for this application
//
// swagger:model
type Host struct {
	// ID of the host
	//
	// required: true
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)"`

	// Name
	//
	// required: true
	Name string `gorm:"column:name"`

	// User ID
	//
	// required: true
	UserId uuid.UUID `gorm:"column:user_id"`

	// Host IP address
	//
	// required: true
	Ip string `gorm:"column:ip"`

	//Created at timestamp
	CreatedAt time.Time

	//Updated at timestamp
	UpdatedAt time.Time

	//Deleted at timestamp
	DeletedAt *time.Time
}

func (u *Host) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)
	return nil
}
