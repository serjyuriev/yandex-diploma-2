package repository

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewRepository(t *testing.T) {
	os.Args = []string{"test", "-c", "A:\\dev\\yandex\\yandex-diploma-2\\dev_srv_config.yaml"}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	_, err := NewRepository(logger)
	require.NoError(t, err)
}

func TestReadUserByLogin(t *testing.T) {
	cfg := config.ServerConfig{}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}
		uid := uuid.New()
		mongoID := primitive.NewObjectID()
		expectedUser := &models.User{
			ID:        uid,
			Login:     "tester",
			Password:  "somepwd",
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			"foo.bar",
			mtest.FirstBatch, bson.D{
				{Key: "_id", Value: mongoID},
				{Key: "id", Value: expectedUser.ID},
				{Key: "login", Value: expectedUser.Login},
				{Key: "password", Value: expectedUser.Password},
				{Key: "logins", Value: expectedUser.Logins},
				{Key: "cards", Value: expectedUser.BankCards},
				{Key: "texts", Value: expectedUser.Texts},
				{Key: "binaries", Value: expectedUser.Binaries},
			},
		))

		user, err := repo.ReadUserByLogin(context.Background(), expectedUser.Login)
		require.NoError(t, err)
		require.Equal(t, expectedUser, user)
	})

	mt.Run("no user", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}

		login := "test"
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		_, err := repo.ReadUserByLogin(context.Background(), login)
		require.Error(t, err)
	})
}

func TestReadUserByID(t *testing.T) {
	cfg := config.ServerConfig{}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}
		uid := uuid.New()
		mongoID := primitive.NewObjectID()
		expectedUser := &models.User{
			ID:        uid,
			Login:     "tester",
			Password:  "somepwd",
			Logins:    make([]*models.LoginPasswordItem, 0),
			BankCards: make([]*models.BankCardItem, 0),
			Texts:     make([]*models.TextItem, 0),
			Binaries:  make([]*models.BinaryItem, 0),
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			"foo.bar",
			mtest.FirstBatch, bson.D{
				{Key: "_id", Value: mongoID},
				{Key: "id", Value: expectedUser.ID},
				{Key: "login", Value: expectedUser.Login},
				{Key: "password", Value: expectedUser.Password},
				{Key: "logins", Value: expectedUser.Logins},
				{Key: "cards", Value: expectedUser.BankCards},
				{Key: "texts", Value: expectedUser.Texts},
				{Key: "binaries", Value: expectedUser.Binaries},
			},
		))

		user, err := repo.ReadUserByID(context.Background(), expectedUser.ID)
		require.NoError(t, err)
		require.Equal(t, expectedUser, user)
	})

	mt.Run("no user", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}

		uid := uuid.New()
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		_, err := repo.ReadUserByID(context.Background(), uid)
		require.Error(t, err)
	})

	mt.Run("nil uuid", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}

		uid := uuid.Nil
		_, err := repo.ReadUserByID(context.Background(), uid)
		require.Error(t, err)
	})
}

func TestCreateItem(t *testing.T) {
	cfg := config.ServerConfig{}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}

		uid := uuid.New()
		item := &models.LoginPasswordItem{
			Login:    "test",
			Password: "test",
			Meta: map[string]string{
				"one":   "two",
				"three": "four",
				"five":  "six",
			},
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.CreateItem(context.Background(), item, "logins", uid)
		require.NoError(t, err)
	})

	mt.Run("insert err", func(mt *mtest.T) {
		repo := &repository{
			cfg:    cfg,
			logger: logger,
			client: nil,
			users:  mt.Coll,
		}

		uid := uuid.New()
		item := &models.LoginPasswordItem{
			Login:    "test",
			Password: "test",
			Meta: map[string]string{
				"one":   "two",
				"three": "four",
				"five":  "six",
			},
		}

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		err := repo.CreateItem(context.Background(), item, "logins", uid)
		require.Error(t, err)
	})
}
