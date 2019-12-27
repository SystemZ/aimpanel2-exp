package model

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/systemz/aimpanel2/lib"
	"log"
	"time"
)

// Single host
// swagger:model Host
type Host struct {
	// ID of the host
	// Example: 100112233-4455-6677-8899-aabbccddeeff
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)" json:"id"`

	// User assigned name
	// Example: My Great Linux server
	Name string `gorm:"column:name" json:"name"`

	// User ID
	UserId uuid.UUID `gorm:"column:user_id" json:"user_id"`

	// User assigned host ip
	// Example: 192.51.100.128
	Ip string `gorm:"column:ip" json:"ip"`

	// Generated token for host
	Token string `gorm:"column:token" json:"token"`

	// Metric frequency
	// Example: 30
	MetricFrequency int `gorm:"column:metric_frequency" json:"metric_frequency"`

	// Host OS
	// Example: linux
	OS string `gorm:"column:os" json:"os"`

	// Host platform
	// Example: ubuntu
	Platform string `gorm:"column:platform" json:"platform"`

	// Host platform family
	// Example: debian
	PlatformFamily string `gorm:"column:platform_family" json:"platform_family"`

	// Host platform version
	// Example: 18.04
	PlatformVersion string `gorm:"column:platform_version" json:"platform_version"`

	// Host kernel version
	// Example: 5.3.0-1-generic
	KernelVersion string `gorm:"column:kernel_version" json:"kernel_version"`

	// Host arch
	// Example: x86_64
	KernelArch string `gorm:"column:kernel_arch" json:"kernel_arch"`

	// State
	// 0 off, 1 running
	// Example: 1
	// required: false
	State uint `gorm:"column:state" json:"state"`

	// Date of host creation
	// Example: 2019-09-29T03:16:27+02:00
	CreatedAt time.Time `json:"created_at"`

	//swagger:ignore
	UpdatedAt time.Time `json:"-"`

	//swagger:ignore
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
