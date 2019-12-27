package model

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/systemz/aimpanel2/lib"
	"log"
	"time"
)

type Host struct {
	// ID of the host
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)" json:"id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// User assigned name
	Name string `gorm:"column:name" json:"name" example:"My Great Linux server"`

	// User ID
	UserId uuid.UUID `gorm:"column:user_id" json:"user_id" example:"100112233-4455-6677-8899-aabbccddeeff"`

	// User assigned host ip
	Ip string `gorm:"column:ip" json:"ip" example:"192.51.100.128"`

	// Generated token for host
	Token string `gorm:"column:token" json:"token" example:"VRAKKUBHNIMKLXSXLWTQAOFGOMSSCXOO"`

	// Metric frequency
	MetricFrequency int `gorm:"column:metric_frequency" json:"metric_frequency" example:"30"`

	// Host OS
	OS string `gorm:"column:os" json:"os" example:"linux"`

	// Host platform
	Platform string `gorm:"column:platform" json:"platform" example:"ubuntu"`

	// Host platform family
	PlatformFamily string `gorm:"column:platform_family" json:"platform_family" example:"debian"`

	// Host platform version
	PlatformVersion string `gorm:"column:platform_version" json:"platform_version" example:"18.04"`

	// Host kernel version
	KernelVersion string `gorm:"column:kernel_version" json:"kernel_version" example:"5.3.0-1-generic"`

	// Host arch
	KernelArch string `gorm:"column:kernel_arch" json:"kernel_arch" example:"x86_64"`

	// State
	// 0 off, 1 running
	State uint `gorm:"column:state" json:"state" example:"1"`

	// Date of host creation
	CreatedAt time.Time `json:"created_at" example:"2019-09-29T03:16:27+02:00"`

	UpdatedAt time.Time `json:"-"`

	DeletedAt *time.Time `json:"-"`
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
	var host Host
	if db.Where("id = ?", hostId).First(&host).RecordNotFound() {
		return ""
	}

	return host.Token
}

func GetHostByToken(db *gorm.DB, token string) *Host {
	var host Host
	if db.Where("token = ?", token).First(&host).RecordNotFound() {
		return nil
	}
	return &host
}
