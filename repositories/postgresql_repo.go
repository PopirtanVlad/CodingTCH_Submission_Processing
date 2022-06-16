package repositories

import (
	"Licenta_Processing_Service/dtos"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresSQLRepo struct {
	db *gorm.DB
}

func NewPostgresSQLRepo() *PostgresSQLRepo {
	db := initTables()
	return &PostgresSQLRepo{
		db: db,
	}
}

func initTables() *gorm.DB {
	postgresConf := dtos.PostgresSQLConfig{
		PostgresDialect:  "postgres",
		PostgresHost:     "ec2-3-248-121-12.eu-west-1.compute.amazonaws.com",
		PostgresDBport:   5432,
		PostgresUser:     "mvllytxvxxpayt",
		PostgresName:     "d98shrm1soi4q9",
		PostgresPassword: "f401ca25b286bd014990ccc47c38a91d5031ef6a3066accdb0d21ac1542ff802",
	}
	// Database connection string
	_ = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", postgresConf.PostgresHost, postgresConf.PostgresDBport, postgresConf.PostgresUser, postgresConf.PostgresPassword, postgresConf.PostgresName)

	//Opening connection to db
	db, err := gorm.Open(postgres.Open("postgres://lwgfseegzupekp:f401ca25b286bd014990ccc47c38a91d5031ef6a3066accdb0d21ac1542ff802@ec2-34-246-227-219.eu-west-1.compute.amazonaws.com:5432/d7hbs5jbnqlcjd"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to db!")

	db.AutoMigrate(&dtos.User{})
	db.AutoMigrate(&dtos.Problem{})
	db.AutoMigrate(&dtos.Submission{})
	db.AutoMigrate(&dtos.TestCase{})
	db.AutoMigrate(&dtos.TestResult{})

	return db
}

func (postgresSQLRepo *PostgresSQLRepo) UpdateSubmission(submission dtos.Submission) error {
	if err := postgresSQLRepo.db.Model(&dtos.Submission{}).Where("id = ?", submission.Id).Update("submission_status", submission.SubmissionStatus).Error; err != nil {
		return errors.Wrapf(err, "Couldn't update the field: , for the submission: %s", submission.Id)
	}

	return nil
}

func (postgresSQLRepo *PostgresSQLRepo) GetSubmission(submissionId string) (*dtos.Submission, error) {
	submission := dtos.Submission{}
	if err := postgresSQLRepo.db.Model(&dtos.Submission{}).Find(&submission, "id = ?", submissionId).Limit(1).Error; err != nil {
		return nil, errors.Wrapf(err, "The submission with id: %s couldn't be found", submissionId)
	}
	return &submission, nil
}

func (postgresSQLRepo *PostgresSQLRepo) GetTests(problemId string) ([]dtos.TestCase, error) {
	var tests []dtos.TestCase
	if err := postgresSQLRepo.db.Model(&dtos.TestCase{}).Find(&tests, "problem_id = ?", problemId).Error; err != nil {
		return nil, errors.Wrapf(err, "The submission with id: %s couldn't be found", problemId)
	}
	return tests, nil
}

func (postgresSQLRepo *PostgresSQLRepo) GetProblem(problemId string) (*dtos.Problem, error) {
	problem := &dtos.Problem{}

	if err := postgresSQLRepo.db.Model(&dtos.Problem{}).Where("id = ?", problemId).Scan(&problem).Limit(1).Error; err != nil {
		return nil, errors.Wrapf(err, "The problem with id: %s couldn't be found", problemId)
	}

	return problem, nil
}

func (postgresSQLRepo *PostgresSQLRepo) SaveTestResults(testResults []*dtos.TestResult) error {
	if err := postgresSQLRepo.db.Model(&dtos.TestResult{}).Save(testResults).Error; err != nil {
		return errors.Wrap(err, "Couldn't save the test results")
	}

	return nil
}
