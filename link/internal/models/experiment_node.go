package models

type ExperimentNodeStatus string

const (
	ExperimentNodeStatusPending  ExperimentNodeStatus = "PENDING"
	ExperimentNodeStatusAccepted ExperimentNodeStatus = "ACCEPTED"
	ExperimentNodeStatusRejected ExperimentNodeStatus = "REJECTED"
	ExperimentNodeStatusTraining ExperimentNodeStatus = "TRAINING"
	ExperimentNodeStatusStopped  ExperimentNodeStatus = "STOPPED"
	ExperimentNodeStatusCompleted ExperimentNodeStatus = "COMPLETED"
	ExperimentNodeStatusFailed    ExperimentNodeStatus = "FAILED"
	ExperimentNodeStatusPreparing ExperimentNodeStatus = "PREPARING"
	ExperimentNodeStatusChecksumMismatch ExperimentNodeStatus = "CHECKSUM_MISMATCH"
)

type ExperimentNode struct {
	ExperimentID uint `gorm:"primaryKey"`
	NodeID       uint `gorm:"primaryKey"`
	MetadataID   uint `gorm:"primaryKey"`
	Status       ExperimentNodeStatus
	Experiment   Experiment `gorm:"foreignKey:ExperimentID"`
	Node         Node       `gorm:"foreignKey:NodeID"`
	Metadata     Metadata   `gorm:"foreignKey:MetadataID"`
}
