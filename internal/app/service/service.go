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
	ErrInvalidCredentials = errors.New("login and/or password incorrect")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotExists      = errors.New("user doesn't exist")
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

// SignUpUser hashes user password and adds user to the database,
// returning user's uuid.
func (s *Service) SignUpUser(ctx context.Context, user *models.User) (string, error) {
	s.logger.Debug().Str("user", user.Login).Msg("checking if such user already exists")
	dbUser, err := s.repo.ReadUserByLogin(ctx, user.Login)
	if err != nil {
		if err != repository.ErrNoUser {
			s.logger.
				Err(err).
				Caller().
				Str("user", user.Login).
				Msg("unable to check if user already exists")
			return "", err
		}
	}
	if err == nil && user.Login == dbUser.Login {
		s.logger.Info().Str("user", user.Login).Msg("user with provided login already exists in the system")
		return "", ErrUserExists
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
func (s *Service) LoginUser(ctx context.Context, user *models.User) (string, error) {
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
func (s *Service) hashUserPassword(password string) string {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(s.cfg.Salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
