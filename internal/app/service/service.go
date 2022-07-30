package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/repository"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
)

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
func (s *Service) SignUpUser(ctx context.Context, user *models.User) error {
	s.logger.Debug().Str("user", user.Login).Msg("signing new user up")

	s.logger.Debug().Str("user", user.Login).Msg("hashing user's password")
	user.Password = s.hashUserPassword(user.Password)

	s.logger.Debug().Str("user", user.Login).Msg("passing user info to data layer")
	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.
			Err(err).
			Caller().
			Str("user", user.Login).
			Msg("unable to create user entry in database")
		return err
	}

	s.logger.Info().Str("user", user.Login).Msg("new user was successfully signed up")
	return nil
}

// hashUserPassword returns hashed with sha256 algorythm password.
func (s *Service) hashUserPassword(password string) string {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.cfg.Salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
