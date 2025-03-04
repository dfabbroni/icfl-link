package models

import "time"

type Node struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"type:varchar(255);unique;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	PublicKey string `gorm:"type:varchar(512);unique;not null"`
	Approved  bool   `gorm:"default:false"`
	LastSeen  time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
