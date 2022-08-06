package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/gokeeperclt"
)

var (
	// win-flags:
	buildVersion string = "N/A" // -X main.buildVersion=v1.0.0
	buildDate    string = "N/A" // -X 'main.buildDate=$(Get-Date)'
)

func main() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	clt, err := gokeeperclt.NewClient(buildVersion, buildDate)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to initialize client")
	}

	if err = clt.Run(); err != nil {
		logger.Fatal().Err(err).Msg("unable to launch client")
	}
}
