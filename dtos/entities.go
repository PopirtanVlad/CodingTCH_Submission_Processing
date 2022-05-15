package dtos

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	gorm.Model

	UserID         uint64
	DisplayName    string
	UserEmail      string
	Enabled        bool
	UserPassword   string
	Provider       string
	ProviderUserId string
}

type Problem struct {
	gorm.Model

	ProblemID         uuid.UUID
	ProblemDifficulty ProblemDifficulty
	ProblemStatement  string
	ProblemTitle      string
	ProblemProposerID uuid.UUID
	TimeLimit         time.Duration
	MemoryLimit       uint64
	TestCases         []TestCase
}

type TestCase struct {
	gorm.Model

	InputFileName          string
	ExpectedOutputFileName string
}

type Submission struct {
	gorm.Model

	SubmissionID        uuid.UUID
	ProblemID           uuid.UUID
	UserId              uuid.UUID
	ProgrammingLanguage ProgrammingLanguage
	UploadTime          time.Time
	TestResults         []*TestResult
}

type TestResult struct {
	TestName     string
	Correct      bool
	TimeElapsed  time.Duration
	MemoryUsed   uint64
	ErrorMessage string
}
