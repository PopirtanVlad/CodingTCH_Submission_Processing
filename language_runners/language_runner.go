package language_runners

import "Licenta_Processing_Service/dtos"

type LanguageRunerConf struct {
}

type LanguageRunner interface {
	RunSubmission(solutionReq *dtos.SolutionRequest) ([]*dtos.TestResult, error)
}
