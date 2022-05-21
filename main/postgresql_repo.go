package main

import (
	"Licenta_Processing_Service/dtos"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

var submission dtos.Submission

type PostgresSQLRepo struct {
	db *gorm.DB
}

func NewPostgreSQLRepo() *PostgresSQLRepo {
	return &PostgresSQLRepo{}
}

func (postgresSQLRepo *PostgresSQLRepo) getSubmission(db *gorm.DB) {
	solution := db.Find(&submission, "id = ?", "f467b4ff-dd3b-40a5-810a-61dbcbba4237")
	json.Marshal(solution)
	fmt.Println(solution)

}

func main() {
	//Loading env variables
	postgresConf := dtos.PostgresSQLConfig{
		PostgresDialect:  "postgres",
		PostgresHost:     "localhost",
		PostgresDBport:   5432,
		PostgresUser:     "postgres",
		PostgresName:     "postgres",
		PostgresPassword: "123",
	}

	// Database connection string
	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", postgresConf.PostgresHost, postgresConf.PostgresDBport, postgresConf.PostgresUser, postgresConf.PostgresPassword, postgresConf.PostgresName)
	println(dbURI)
	//Opening connection to db
	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to db!")

	//Close connection to database when main function finishes
	//Make database migrations
	db.AutoMigrate(&dtos.User{})
	db.AutoMigrate(&dtos.Problem{})
	db.AutoMigrate(&dtos.Submission{})
	db.AutoMigrate(&dtos.TestCase{})

	NewPostgreSQLRepo().getSubmission(db)
}
