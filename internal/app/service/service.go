// Package service provides service layer implementation of gokeeper server app.
package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/repository"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
)

var (
	// ErrInvalidCredentails is raised when provided user's credentials are not correct.
	ErrInvalidCredentials = errors.New("login and/or password incorrect")
	// ErrUserNotExists is raised when client tries to login with non-existing user.
	ErrUserNotExists = errors.New("user doesn't exist")
	// ErrNilArgument is raised when client makes a call for a method without
	// providing enough information.
	ErrNilArgument = errors.New("argument can't be empty")
)

// Service provides service layer methods.
type Service interface {
	SignUpUser(ctx context.Context, user *models.User) (string, error)
	LoginUser(ctx context.Context, user *models.User) (string, error)
}

// Service holds objects for service layer implementation.
type service struct {
	cfg    config.ServerConfig
	repo   repository.Repository
	logger zerolog.Logger
}

// NewService initializes app's service layer.
func NewService(logger zerolog.Logger, repo repository.Repository) (Service, error) {
	if repo == nil {
		logger.Err(ErrNilArgument).Str("arg", "repo").Msg("repository can't be nil")
		return nil, ErrNilArgument
	}

	logger.Debug().Str("module", "service").Msg("getting app's configuration")
	cfg := config.GetServerConfig()

	logger.Info().Msg("service layer was successfully initialized")
	return &service{
		cfg:    cfg,
		repo:   repo,
		logger: logger,
	}, nil
}

// SignUpUser hashes user password and adds user to the database,
// returning user's uuid.
func (s *service) SignUpUser(ctx context.Context, user *models.User) (string, error) {
	if user == nil {
		s.logger.Err(ErrNilArgument).Str("arg", "user").Msg("user can't be nil")
		return "", ErrNilArgument
	}

	s.logger.Debug().Str("user", user.Login).Msg("generating user's uuid")
	user.ID = uuid.New()

	s.logger.Debug().Str("user", user.Login).Msg("hashing user's password")
	user.Password = s.hashUserPassword(user.Password)

	s.logger.Debug().Str("user", user.Login).Msg("passing user's info to data layer")
	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.
			Err(err).
			Caller().
			Str("user", user.Login).
			Msg("unable to create user entry in database")
		return "", err
	}

	return user.ID.String(), nil
}

// LoginUser checks whether user exists in the database and
// if user's credentials are equals, logins user.
func (s *service) LoginUser(ctx context.Context, user *models.User) (string, error) {
	if user == nil {
		s.logger.Err(ErrNilArgument).Str("arg", "user").Msg("user can't be nil")
		return "", ErrNilArgument
	}

	s.logger.Debug().Str("user", user.Login).Msg("checking if such user exists")
	dbUser, err := s.repo.ReadUserByLogin(ctx, user.Login)
	if err != nil {
		if err != repository.ErrNoUser {
			s.logger.
				Err(err).
				Caller().
				Str("user", user.Login).
				Msg("unable to check if user exists")
			return "", err
		} else {
			s.logger.Info().Str("user", user.Login).Msg("user with provided login doesn't exist in the system")
			return "", ErrUserNotExists
		}
	}
	s.logger.Debug().Str("user", user.Login).Msg("user with provided login exists")

	s.logger.Debug().Str("user", user.Login).Msg("hashing provided user's password")
	user.Password = s.hashUserPassword(user.Password)

	s.logger.Debug().Str("user", user.Login).Msg("checking credentials")
	if user.Login == dbUser.Login && user.Password == dbUser.Password {
		s.logger.Debug().Str("user", user.Login).Msg("provided credentials are correct")
		return dbUser.ID.String(), nil
	}

	s.logger.Debug().Str("user", user.Login).Msg("provided credentials aren't correct")
	return "", ErrInvalidCredentials
}

// hashUserPassword returns hashed with sha256 algorythm password.
func (s *service) hashUserPassword(password string) string {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.cfg.Salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
