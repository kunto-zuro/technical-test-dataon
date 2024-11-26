package models

import "time"

type Node struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"unique;not null" json:"code"`
	Name      string    `gorm:"not null" json:"name"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
