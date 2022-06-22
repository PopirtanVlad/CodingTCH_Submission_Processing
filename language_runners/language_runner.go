package language_runners

import "Licenta_Processing_Service/entities"

type LanguageRunerConf struct {
}

type LanguageRunner interface {
	RunSubmission(solutionReq *entities.SolutionRequest) ([]*entities.TestResult, error)
}
