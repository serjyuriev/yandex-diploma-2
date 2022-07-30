package main

import (
	"log"

	"github.com/serjyuriev/yandex-diploma-2/internal/app/gokeepersrv"
)

func main() {
	srv, err := gokeepersrv.NewServer()
	if err != nil {
		log.Fatalf("unable to initialize server: %v", err)
	}

	if err = srv.Start(); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
