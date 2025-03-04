package models

import "time"

type InstructionType string

const (
	InstructionNewExperiment    InstructionType = "NEW_EXPERIMENT"
	InstructionStartTraining    InstructionType = "START_TRAINING"
	InstructionStopTraining     InstructionType = "STOP_TRAINING"
	InstructionUpdateExperiment InstructionType = "UPDATE_EXPERIMENT"
)

type Instruction struct {
	Type      InstructionType
	Payload   interface{}
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
