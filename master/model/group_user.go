package model

import (
	"github.com/gofrs/uuid"
	"time"
)

// User group represents the group for this application
// swagger:model userGroup
type GroupUser struct {
	// ID of the group
	//
	// required: true
	ID int `gorm:"primary_key" json:"id"`

	// ID of the group
	//
	// required: true
	GroupId int `gorm:"column:group_id" json:"group_id"`

	// ID of the user
	//
	// required: true
	UserId uuid.UUID `gorm:"column:user_id" json:"user_id"`

	// Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	// Updated at timestamp
	UpdatedAt time.Time `json:"-"`

	// Deleted at timestamp
	DeletedAt *time.Time `json:"-"`
}
