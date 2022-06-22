package repositories

import (
	"Licenta_Processing_Service/entities"
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
	postgresConf := entities.PostgresSQLConfig{
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

	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Problem{})
	db.AutoMigrate(&entities.Submission{})
	db.AutoMigrate(&entities.TestCase{})
	db.AutoMigrate(&entities.TestResult{})

	return db
}

func (postgresSQLRepo *PostgresSQLRepo) UpdateSubmission(submission entities.Submission) error {
	if err := postgresSQLRepo.db.Model(&entities.Submission{}).Where("id = ?", submission.Id).Update("submission_status", submission.SubmissionStatus).Error; err != nil {
		return errors.Wrapf(err, "Couldn't update the field: , for the submission: %s", submission.Id)
	}

	return nil
}

func (postgresSQLRepo *PostgresSQLRepo) GetSubmission(submissionId string) (*entities.Submission, error) {
	submission := entities.Submission{}
	if err := postgresSQLRepo.db.Model(&entities.Submission{}).Find(&submission, "id = ?", submissionId).Limit(1).Error; err != nil {
		return nil, errors.Wrapf(err, "The submission with id: %s couldn't be found", submissionId)
	}
	return &submission, nil
}

func (postgresSQLRepo *PostgresSQLRepo) GetTests(problemId string) ([]entities.TestCase, error) {
	var tests []entities.TestCase
	if err := postgresSQLRepo.db.Model(&entities.TestCase{}).Find(&tests, "problem_id = ?", problemId).Error; err != nil {
		return nil, errors.Wrapf(err, "The submission with id: %s couldn't be found", problemId)
	}
	return tests, nil
}

func (postgresSQLRepo *PostgresSQLRepo) GetProblem(problemId string) (*entities.Problem, error) {
	problem := &entities.Problem{}

	if err := postgresSQLRepo.db.Model(&entities.Problem{}).Where("id = ?", problemId).Scan(&problem).Limit(1).Error; err != nil {
		return nil, errors.Wrapf(err, "The problem with id: %s couldn't be found", problemId)
	}

	return problem, nil
}

func (postgresSQLRepo *PostgresSQLRepo) SaveTestResults(testResults []*entities.TestResult) error {
	if err := postgresSQLRepo.db.Model(&entities.TestResult{}).Save(testResults).Error; err != nil {
		return errors.Wrap(err, "Couldn't save the test results")
	}

	return nil
}
