package models

import "time"

type Metadata struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	NodeID         uint `gorm:"uniqueIndex:idx_node_metadata"`
	NodeMetadataID uint `gorm:"uniqueIndex:idx_node_metadata"`
	Name           string
	Type           string
	Tags           string
	Description    string
	Extras         string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	Node           Node `gorm:"foreignKey:NodeID"`
}
