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

type SolutionRequest struct {
	File     io.ReadCloser
	Solution Submission
	Tests    []TestCase
}

type CommandConfig struct {
	CommandName string
	CommandArgs []string
	TimeOut     time.Duration
	StdIn       io.Reader
	StdOut      io.Writer
}
