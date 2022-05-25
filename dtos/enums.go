package dtos

type ProgrammingLanguage string
type ProblemDifficulty string

const (
	Python3 ProgrammingLanguage = "Python"
	Java    ProgrammingLanguage = "Java"
	C       ProgrammingLanguage = "C"
)

const (
	Easy   ProblemDifficulty = "Easy"
	Normal ProblemDifficulty = "Normal"
	Hard   ProblemDifficulty = "Hard"
)
