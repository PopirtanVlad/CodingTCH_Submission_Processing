package dtos

import (
	"bytes"
	"io"
	"time"
)

type SolutionResult struct {
	ExecutionTime time.Duration
	MemoryUsage   uint64
	StdErr        *bytes.Buffer
	ExitCode      int
	Error         error
}

type RunTestRequest struct {
	Submission     Submission
	Test           TestCase
	OutputFileName string
}

type SolutionRequest struct {
	File       io.ReadCloser
	Submission Submission
	Tests      []TestCase
}

type CommandConfig struct {
	CommandName string
	CommandArgs []string
	TimeOut     time.Duration
	StdIn       io.Reader
	StdOut      io.Writer
}
