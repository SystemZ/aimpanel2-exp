package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Group represents the group for this application
// swagger:model group
type Group struct {
	// ID of the group
	//
	// required: true
	ID int `gorm:"primary_key" json:"id"`

	// Name of the group
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	// Updated at timestamp
	UpdatedAt time.Time `json:"-"`

	// Deleted at timestamp
	DeletedAt *time.Time `json:"-"`
}

func GetGroup(db *gorm.DB, name string) *Group {
	var group Group
	if db.Where("name = ?", name).First(&group).RecordNotFound() {
		return nil
	}
	return &group
}
