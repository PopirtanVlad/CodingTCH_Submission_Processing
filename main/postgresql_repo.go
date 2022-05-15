package main

import (
	"Licenta_Processing_Service/dtos"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

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

}
