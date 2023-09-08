package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Db *sql.DB
)

func ConnectDatabase() {
	var err error

	godotenv.Load()

	Db, err = sql.Open("sqlite3", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Db Connected")
}
