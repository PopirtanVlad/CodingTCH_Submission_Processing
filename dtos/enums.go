package dtos

type ProgrammingLanguage string
type ProblemDifficulty string

const (
	Python3 ProgrammingLanguage = "python3"
	Java    ProgrammingLanguage = "java"
	C       ProgrammingLanguage = "c"
)

const (
	Easy   ProblemDifficulty = "Easy"
	Normal ProblemDifficulty = "Normal"
	Hard   ProblemDifficulty = "Hard"
)
