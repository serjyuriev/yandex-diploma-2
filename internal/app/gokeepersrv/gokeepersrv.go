package gokeepersrv

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/handlers"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
)

var errNotImplemented = errors.New("not implemented yet")

// Server holds app's server-side related objects.
type Server struct {
	cfg    config.ServerConfig
	rpc    *handlers.RPC
	logger zerolog.Logger
}

// NewServer initializes app's server.
func NewServer() (*Server, error) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	logger.Debug().Msg("initializing go-keeper server")

	logger.Debug().Msg("getting app's configuration")
	cfg := config.GetServerConfig()

	logger.Debug().Msg("initializing gRPC layer")
	rpc, err := handlers.MakeRPC(logger)
	if err != nil {
		logger.
			Err(err).
			Caller().
			Msg("unable to initialize gRPC layer")
		return nil, err
	}

	logger.Info().Msg("go-keeper server was successfully initialized")
	return &Server{
		cfg:    cfg,
		rpc:    rpc,
		logger: logger,
	}, nil
}

// Start launches app's server.
func (s *Server) Start() error {
	return errNotImplemented
}
