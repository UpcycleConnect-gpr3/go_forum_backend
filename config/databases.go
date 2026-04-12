package config

import (
	"go-forum-backend/database"
	"go-forum-backend/internal"
	"os"
)

func InitDatabase() {

	database.Forum = internal.NewDatabase(
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
}
