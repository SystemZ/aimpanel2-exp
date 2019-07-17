package model

import (
	"time"
)

// Group represents the group for this application
// swagger:model group
type Group struct {
	// ID of the group
	//
	// required: true
	ID uint `gorm:"primary_key" json:"id"`

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
