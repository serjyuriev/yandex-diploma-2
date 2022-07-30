package main

import (
	"github.com/rs/zerolog/log"

	"github.com/serjyuriev/yandex-diploma-2/internal/app/gokeepersrv"
)

func main() {
	srv, err := gokeepersrv.NewServer()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize new server")
	}

	if err = srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("unable to start server")
	}
}
