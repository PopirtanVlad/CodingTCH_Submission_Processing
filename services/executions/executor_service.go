package executions

import (
	"Licenta_Processing_Service/custom_errors"
	"Licenta_Processing_Service/dtos"
	"bytes"
	_ "fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	_ "os"
	"os/exec"
	"time"
)

type ExecutionRunner struct {
	timeLimit   time.Duration
	memoryLimit int
}

func NewExecutionRunner() *ExecutionRunner {
	return &ExecutionRunner{
		timeLimit:   3000009900,
		memoryLimit: 30,
	}
}

//func main() {
//	_, err := NewExecutionRunner().runSolution(dtos.CommandConfig{
//		CommandArgs: []string{"hello_world.py", "1", "2"},
//		CommandName: "py",
//	})
//	if err == custom_errors.TimeLimitExceededError {
//		fmt.Println(err)
//	}
//}

func (executionRunner *ExecutionRunner) RunCommand(cmdConfig dtos.CommandConfig) (*dtos.SolutionResult, error) {
	cmd := exec.Command(cmdConfig.CommandName, cmdConfig.CommandArgs...)
	startTime := time.Now()
	var errBuff bytes.Buffer

	cmd.Stderr = &errBuff
	cmd.Stdout = cmdConfig.StdOut
	cmd.Stdin = cmdConfig.StdIn
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "Could not run the command: %s \nError: %s", cmd.String(), errBuff.String())
	}

	cmd.Wait()
	endTime := time.Since(startTime)
	if endTime > executionRunner.timeLimit {
		logrus.WithFields(logrus.Fields{
			"Time limit":  executionRunner.timeLimit,
			"Actual time": endTime,
		}).Debugf("Time limit exceeded")
		return nil, custom_errors.TimeLimitExceededError
	}

	//if memoryUsed > executionRunner.memoryLimit {
	//	logrus.WithFields(logrus.Fields{
	//		"Memory limit":       executionRunner.memoryLimit,
	//		"Actual memory used": memoryUsed,
	//	}).Debugf("Memory limit exceeded")
	//	return nil, custom_errors.MemoryLimitExceededError
	//}

	return &dtos.SolutionResult{
		//MemoryUsage: memoryUsed
		ExecutionTime: endTime,
		Error:         nil,
	}, nil
}
