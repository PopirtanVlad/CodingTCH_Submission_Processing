package dtos

import (
	"github.com/satori/go.uuid"
	"time"
)

type SolutionResult struct {
	solutionID     uuid.UUID
	solutionStatus bool
	executionTime  time.Duration
	memoryUsage    int
}
