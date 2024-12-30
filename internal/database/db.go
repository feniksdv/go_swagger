package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func Connect() *sql.DB {
	driverName, dataSourceName := getDataSourceName()
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Println("Connect %v", err.Error())
	}

	return db
}

func getDataSourceName() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	driver := os.Getenv("DRIVER")
	dbName := os.Getenv("DB_NAME")

	return driver, dbName
}
