package cmd

import (
	"fmt"
	"go-forum-backend/cmd/database"
	"go-forum-backend/cmd/server"
	"os"
)

func Cmd() {

	if len(os.Args) < 2 {
		fmt.Println("Commande manquante. Utilisation : monexecutable [start|serve]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		server.Start()

	case "migrate":
		database.Migrate()

	default:
		fmt.Println("Commande inconnue. Utilisation : go main serve")
		os.Exit(1)
	}
}
