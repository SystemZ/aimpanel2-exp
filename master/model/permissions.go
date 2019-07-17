package model

import (
	"time"
)

// Group represents the group for this application
// swagger:model group
type Permission struct {
	// ID of the group
	//
	// required: true
	ID uint `gorm:"primary_key" json:"id"`

	// Name of the group
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// ID of the group
	//
	// required: true
	GroupId int `gorm:"column:group_id" json:"id"`

	Verb uint8 `gorm:"column:verb" json:"verb"`

	Endpoint string `gorm:"column:endpoint" json:"endpoint"`

	// Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	// Updated at timestamp
	UpdatedAt time.Time `json:"-"`

	// Deleted at timestamp
	DeletedAt *time.Time `json:"-"`
}
