package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/mocks"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/repository"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	os.Args = []string{"test", "-c", "A:\\dev\\yandex\\yandex-diploma-2\\dev_srv_config.yaml"}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	repo, err := repository.NewRepository(logger)
	require.NoError(t, err)

	_, err = NewService(logger, repo)
	require.NoError(t, err)
}

func TestSignUpUser(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		newUser := &models.User{
			Login:     "test",
			Password:  "somepwd",
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(nil, repository.ErrNoUser)
		create := mr.EXPECT().CreateUser(
			context.Background(),
			newUser,
		).Return(nil)

		gomock.InOrder(read, create)

		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}
		_, err := svc.SignUpUser(context.Background(), newUser)
		require.NoError(t, err)
	})

	t.Run("read err", func(t *testing.T) {
		newUser := &models.User{
			Login: "test",
		}

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(nil, fmt.Errorf("some err"))
		gomock.InOrder(read)

		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}
		_, err := svc.SignUpUser(context.Background(), newUser)
		require.Error(t, err)
	})

	t.Run("existing user", func(t *testing.T) {
		newUser := &models.User{
			Login: "test",
		}

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(newUser, nil)
		gomock.InOrder(read)

		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}
		_, err := svc.SignUpUser(context.Background(), newUser)
		require.ErrorIs(t, err, ErrUserExists)
	})

	t.Run("create err", func(t *testing.T) {
		newUser := &models.User{
			Login:     "test",
			Password:  "somepwd",
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(nil, repository.ErrNoUser)
		create := mr.EXPECT().CreateUser(
			context.Background(),
			newUser,
		).Return(fmt.Errorf("some err"))

		gomock.InOrder(read, create)

		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}
		_, err := svc.SignUpUser(context.Background(), newUser)
		require.Error(t, err)
	})
}

func TestLoginUser(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}

		user := &models.User{
			Login:     "test",
			Password:  "somepwd",
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}
		uid := uuid.New()

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(&models.User{
			ID:        uid,
			Login:     "test",
			Password:  svc.hashUserPassword("somepwd"),
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}, nil)

		gomock.InOrder(read)
		userID, err := svc.LoginUser(context.Background(), user)
		require.NoError(t, err)
		require.Equal(t, uid.String(), userID)
	})

	t.Run("read err", func(t *testing.T) {
		newUser := &models.User{
			Login: "test",
		}

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(nil, fmt.Errorf("some err"))
		gomock.InOrder(read)

		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}
		_, err := svc.LoginUser(context.Background(), newUser)
		require.Error(t, err)
	})

	t.Run("no user", func(t *testing.T) {
		newUser := &models.User{
			Login: "test",
		}

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(nil, repository.ErrNoUser)
		gomock.InOrder(read)

		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}
		_, err := svc.LoginUser(context.Background(), newUser)
		require.ErrorIs(t, err, ErrUserNotExists)
	})

	t.Run("invalid creds", func(t *testing.T) {
		svc := &service{
			cfg:    config.ServerConfig{Salt: "testsalt"},
			repo:   mr,
			logger: logger,
		}

		user := &models.User{
			Login:     "test",
			Password:  "somepwd",
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}
		uid := uuid.New()

		read := mr.EXPECT().ReadUserByLogin(
			context.Background(),
			gomock.Eq("test"),
		).Return(&models.User{
			ID:        uid,
			Login:     "test",
			Password:  svc.hashUserPassword("s0m3pwd"),
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}, nil)

		gomock.InOrder(read)
		_, err := svc.LoginUser(context.Background(), user)
		require.ErrorIs(t, err, ErrInvalidCredentials)

	})
}
