package store

import (
	"sync"

	"link/internal/models"
)

type NodeInstruction struct {
	NodeID      uint
	Instruction models.Instruction
}

type InstructionStore struct {
	instructions map[uint][]models.Instruction
	mu           sync.RWMutex
}

var GlobalInstructionStore = &InstructionStore{
	instructions: make(map[uint][]models.Instruction),
}

func (s *InstructionStore) AddInstructions(nodeInstructions []NodeInstruction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ni := range nodeInstructions {
		s.instructions[ni.NodeID] = append(s.instructions[ni.NodeID], ni.Instruction)
	}
}

func (s *InstructionStore) GetInstructions(nodeID uint) []models.Instruction {
	s.mu.Lock()
	defer s.mu.Unlock()
	instructions := s.instructions[nodeID]
	delete(s.instructions, nodeID)
	return instructions
}
