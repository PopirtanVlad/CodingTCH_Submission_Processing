package custom_errors

import "errors"

var (
	MemoryLimitExceededError = errors.New("memory limit exceeded")
	TimeLimitExceededError   = errors.New("time limit exceeded")
)
