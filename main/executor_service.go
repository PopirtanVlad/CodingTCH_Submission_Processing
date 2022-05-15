package main

import (
	"Licenta_Processing_Service/dtos"
	_ "fmt"
	"github.com/sirupsen/logrus"
	_ "os"
	"os/exec"
	"time"
)

type ExecutionRunner struct {
	timeLimit      time.Duration
	memoryLimit    int
	solutionResult dtos.SolutionResult
}

func NewExecutionRunner() *ExecutionRunner {
	return &ExecutionRunner{
		timeLimit:   1,
		memoryLimit: 30,
	}
}

func main() {
	NewExecutionRunner().runSolution()
}

func (executionRunner *ExecutionRunner) runSolution() error {

	cmd := exec.Command("py", "hello_world.py")

	out, err := cmd.Output()
	if err != nil {
		logrus.
			WithFields(logrus.Fields{"error": err}).Fatal("DIDNT WORK")
		return err
	}

	println(string(out))

	return nil
}
