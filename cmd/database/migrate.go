package database

import (
	"go-forum-backend/config"
	"go-forum-backend/database"
	"go-forum-backend/internal"
	"go-forum-backend/utils/log"

	"github.com/joho/godotenv"
)

func initialize() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Config Initialization
	config.InitDatabase()

	err = database.Forum.Ping()

	if err != nil {
		log.Fatal(err)
	}

	internal.CreateTableMigrations(database.Forum)

}

func Migrate() {

	initialize()

	internal.Migrate(database.Forum)

}
