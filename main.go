package main

import (
	"go-forum-backend/cmd"
	"os"
)

func main() {
	if len(os.Args) > 0 {
		cmd.Cmd()
	}
}
