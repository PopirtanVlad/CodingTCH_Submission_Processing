package main

import (
	"Licenta_Processing_Service/dtos"
	"Licenta_Processing_Service/services/executions"
	"fmt"
	"io/ioutil"
)

type CplusplusRunner struct {
}

func NewCplusPlusRunner() *CplusplusRunner {
	return &CplusplusRunner{}
}

func main() {
	_, err := NewCplusPlusRunner().compileSolution("Solution.c")

	fmt.Println(err)
}

func (cplusplusRunner *CplusplusRunner) compileSolution(fileName string) (*dtos.SolutionResult, error) {
	cmdConfig := dtos.CommandConfig{
		CommandName: "gcc",
		CommandArgs: []string{fileName, "-o", "c_solution.exe"},
		TimeOut:     1400000,
		StdIn:       nil,
		StdOut:      ioutil.Discard,
	}

	return executions.NewExecutionRunner().RunCommand(cmdConfig)
}
