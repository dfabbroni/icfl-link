package models

import (
	"time"
)

type Experiment struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Name        string
	Description string
	BasePath    string
	Status      string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	User        User `gorm:"foreignKey:UserID"`
	ExperimentNodes []ExperimentNode `gorm:"foreignKey:ExperimentID"`
}
