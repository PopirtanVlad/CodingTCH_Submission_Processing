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
	Problem        Problem
	Test           TestCase
	OutputFileName string
}

type SolutionRequest struct {
	File        io.ReadCloser
	Problem     Problem
	Submission  Submission
	Tests       []TestCase
	TimeOut     time.Duration
	MemoryLimit uint64
}

type CommandConfig struct {
	CommandName string
	CommandArgs []string
	TimeOut     time.Duration
	StdIn       io.Reader
	StdOut      io.Writer
}
