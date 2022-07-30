package main

import (
	"github.com/rs/zerolog/log"

	"github.com/serjyuriev/yandex-diploma-2/internal/app/gokeepertui"
)

func main() {
	clt, err := gokeepertui.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize client")
	}

	if err = clt.Start(); err != nil {
		log.Fatal().Err(err).Msg("unable to start client")
	}
}
