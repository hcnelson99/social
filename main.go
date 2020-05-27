package main

import (
	"log"
	"os"

	"github.com/hcnelson99/social/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		// default port
		port = "8080"
	}

	app.Start(":" + port)
}
