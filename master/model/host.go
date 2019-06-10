package model

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/systemz/aimpanel2/lib"
	"log"
	"time"
)

// Host represents the host for this application
// swagger:model host
type Host struct {
	// ID of the host
	//
	// required: true
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)" json:"id"`

	// Name
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// User ID
	//
	// required: true
	UserId uuid.UUID `gorm:"column:user_id" json:"user_id"`

	// Host IP address
	//
	// required: true
	Ip string `gorm:"column:ip" json:"ip"`

	// Token
	Token string `gorm:"column:token" json:"token"`

	MetricFrequency int `gorm:"column:metric_frequency" json:"metric_frequency"`

	//Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	//Updated at timestamp
	UpdatedAt time.Time `json:"updated_at"`

	//Deleted at timestamp
	DeletedAt *time.Time `json:"deleted_at"`
}

func (h *Host) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)

	scope.SetColumn("Token", lib.RandomString(32))

	return nil
}

func GetHost(db *gorm.DB, hostId string) *Host {
	var host Host
	if db.Where("id = ?", hostId).First(&host).RecordNotFound() {
		return nil
	}
	return &host
}

func GetHostToken(db *gorm.DB, hostId string) string {
	var token string
	if db.Select("token").Where("id = ?", hostId).First(&token).RecordNotFound() {
		return ""
	}

	return token
}
