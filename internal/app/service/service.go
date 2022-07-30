package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/repository"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
)

var errNotImplemented = errors.New("not implemented yet")

// Service holds objects for service layer implementation.
type Service struct {
	cfg    config.ServerConfig
	repo   *repository.Repository
	logger zerolog.Logger
}

// NewService initializes app's service layer.
func NewService(logger zerolog.Logger, repo *repository.Repository) (*Service, error) {
	logger.Debug().Str("module", "service").Msg("getting app's configuration")
	cfg := config.GetServerConfig()

	logger.Info().Msg("service layer was successfully initialized")
	return &Service{
		cfg:    cfg,
		repo:   repo,
		logger: logger,
	}, nil
}

// SignUpUser hashes user password and adds user to the database.
func (s *Service) SignUpUser(ctx context.Context, user models.User) error {
	return errNotImplemented
}
