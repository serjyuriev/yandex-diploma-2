package repository

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var errNotImplemented = errors.New("not implemented yet")

// Repository holds objects for data layer implementation.
type Repository struct {
	cfg    config.ServerConfig
	client *mongo.Client
	logger zerolog.Logger
}

// NewRepository initializes connection to mongo db.
func NewRepository(logger zerolog.Logger) (*Repository, error) {
	logger.Debug().Str("module", "repo").Msg("getting app's configuration")
	cfg := config.GetServerConfig()

	logger.Debug().Str("module", "repo").Msg("initializing database connection")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.DataSourceName))
	if err != nil {
		logger.
			Err(err).
			Caller().
			Msg("unable to initialize data layer")
		return nil, err
	}

	logger.Info().Msg("data layer was successfully initialized")
	return &Repository{
		cfg:    cfg,
		client: client,
		logger: logger,
	}, nil
}

// CreateUser adds new user entry to the database.
func (r *Repository) CreateUser(ctx context.Context, user models.User) error {
	return errNotImplemented
}
