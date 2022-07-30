package gokeepersrv

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/handlers"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	g "github.com/serjyuriev/yandex-diploma-2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

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
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port))
	if err != nil {
		s.logger.
			Err(err).
			Caller().
			Msgf("unable to listen on %s:%d", s.cfg.Address, s.cfg.Port)
		return err
	}
	srv := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		},
	))
	g.RegisterGokeeperServer(srv, s.rpc)

	s.logger.Info().Msgf("go-keeper server listening on tcp %s:%d", s.cfg.Address, s.cfg.Port)
	if err := srv.Serve(listen); err != nil {
		s.logger.
			Err(err).
			Caller().
			Msg("unexpected error occured")
		return err
	}
	return nil
}
