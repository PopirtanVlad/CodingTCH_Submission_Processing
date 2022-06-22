package entities

import (
	"time"
)

type User struct {
	Id             uint64 `gorm:"primaryKey"`
	DisplayName    string `gorm:"column:display_name"`
	UserEmail      string `gorm:"column:user_email"`
	Enabled        bool   `gorm:"enabled"`
	UserPassword   string `gorm:"user_password"`
	ProviderUserId string `gorm:"provider_user_id"`
}

type Problem struct {
	Id                   string            `gorm:"primaryKey"`
	ProblemDifficulty    ProblemDifficulty //`gorm:"column:problem_difficulty"`
	ProblemExampleInput  string            //`gorm:"column:problem_example_input"`
	ProblemExampleOutput string            //`gorm:"column:problem_example_output"`
	ProblemTitle         string            //`gorm:"column:problem_title"`
	TimeLimit            time.Duration     //`gorm:"column:time_limit"`
	MemoryLimit          uint64            //`gorm:"column:memory_limit"`
	TestCases            []TestCase
}

type TestCase struct {
	Id                     string `gorm:"primaryKey"`
	InputFileName          string `gorm:"column:test_case_input"`
	ExpectedOutputFileName string `gorm:"column:test_case_output"`
	ProblemId              string
}

type Submission struct {
	Id                  string `gorm:"primaryKey"`
	ProblemID           string
	UserId              uint64
	ProgrammingLanguage ProgrammingLanguage `gorm:"column:programming_language"`
	UploadTime          time.Time           `gorm:"column:upload_time"`
	SubmissionStatus    uint64
	TestResults         []*TestResult
}

type TestResult struct {
	Id           string `gorm:"primaryKey"`
	SubmissionId string
	Correct      bool          `gorm:"column:test_status"`
	TimeElapsed  time.Duration `gorm:"test_time_elapsed"`
	MemoryUsed   uint64        `gorm:"test_memory_used"`
	ErrorMessage string        `gorm:"error_message"`
}
