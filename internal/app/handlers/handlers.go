package handlers

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/serjyuriev/yandex-diploma-2/internal/app/repository"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/service"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	g "github.com/serjyuriev/yandex-diploma-2/proto"
)

var errNotImplemented = errors.New("not implemented yet")

// RPC holds objects for grpc implementation.
type RPC struct {
	g.UnimplementedGokeeperServer

	cfg    config.ServerConfig
	repo   *repository.Repository
	svc    *service.Service
	logger zerolog.Logger
}

// MakeRPC initializes app's grpc service.
func MakeRPC(logger zerolog.Logger) (*RPC, error) {
	logger.Debug().Str("module", "gRPC").Msg("getting app's configuration")
	cfg := config.GetServerConfig()

	logger.Debug().Str("module", "gRPC").Msg("initializing data layer")
	repo, err := repository.NewRepository(logger)
	if err != nil {
		logger.
			Err(err).
			Caller().
			Msg("unable to initialize data layer")
		return nil, err
	}

	logger.Debug().Str("module", "gRPC").Msg("initializing service layer")
	svc, err := service.NewService(logger, repo)
	if err != nil {
		logger.
			Err(err).
			Caller().
			Msg("unable to initialize service layer")
		return nil, err
	}

	logger.Info().Msg("gRPC layer was successfully initialized")
	return &RPC{
		cfg:    cfg,
		repo:   repo,
		svc:    svc,
		logger: logger,
	}, nil
}

// SignUpUser signs new user up.
func (r *RPC) SignUpUser(ctx context.Context, in *g.SignUpUserRequest) (*g.SignUpUserResponse, error) {
	r.logger.Info().Str("user", in.User.Login).Msg("received new user sign up request")
	user := &models.User{
		Login:    in.User.Login,
		Password: in.User.Password,
	}
	res := new(g.SignUpUserResponse)

	r.logger.Debug().Msg("passing user's info to service layer")
	userID, err := r.svc.SignUpUser(ctx, user)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("login", user.Login).
			Msg("unable to sign user up")
		res.UserID = ""
		res.Error = err.Error()
		return res, err
	}

	r.logger.Info().Str("user", in.User.Login).Msg("new user was successfully signed up")
	res.UserID = userID
	res.Error = ""
	return res, nil
}
