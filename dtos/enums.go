package dtos

type ProgrammingLanguage int
type ProblemDifficulty string

const (
	Python3 ProgrammingLanguage = 1
	Java    ProgrammingLanguage = 0
	C       ProgrammingLanguage = 2
)

const (
	Easy   ProblemDifficulty = "Easy"
	Normal ProblemDifficulty = "Normal"
	Hard   ProblemDifficulty = "Hard"
)
