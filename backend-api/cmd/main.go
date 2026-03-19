package main

import (
	"log"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/app"
)

func main() {
	srv := app.New()

	if err := srv.ConfigureService(); err != nil {
		log.Fatal(err)
	}

	if err := srv.StartServe(); err != nil {
		log.Fatal(err)
	}
}
