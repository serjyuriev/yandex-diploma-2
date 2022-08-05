package main

import (
	"github.com/rs/zerolog/log"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/gokeeperclt"
)

func main() {
	clt, err := gokeeperclt.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize client")
	}

	if err = clt.Run(); err != nil {
		log.Fatal().Err(err).Msg("unable to start client")
	}
}
