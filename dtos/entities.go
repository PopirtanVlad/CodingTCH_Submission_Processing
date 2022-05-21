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
	ProviderUserId string
}

type Problem struct {
	gorm.Model

	Id                uuid.UUID
	ProblemDifficulty ProblemDifficulty
	ProblemStatement  string
	ProblemTitle      string
	TimeLimit         time.Duration
	MemoryLimit       uint64
	TestCases         []TestCase `gorm:"foreignKey:Id"`
}

type TestCase struct {
	gorm.Model

	Id                     uuid.UUID
	InputFileName          string
	ExpectedOutputFileName string
}

type Submission struct {
	gorm.Model

	Id                  uuid.UUID
	ProblemID           uuid.UUID `gorm:"referecnes:problem_id"`
	UserId              uint64
	ProgrammingLanguage ProgrammingLanguage
	UploadTime          time.Time
	TestResults         []*TestResult `gorm:"foreignKey:Id"`
}

type TestResult struct {
	Id           uuid.UUID
	Correct      bool
	TimeElapsed  time.Duration
	MemoryUsed   uint64
	ErrorMessage string
}
