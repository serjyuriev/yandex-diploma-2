package repository

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Database.Address))
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
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	r.logger.Debug().Str("user", user.Login).Msg("getting users collection")
	collection := r.client.Database(r.cfg.Database.Name).Collection("users")

	r.logger.Debug().Str("user", user.Login).Msg("marshalling user's info to bson")
	doc, err := bson.Marshal(user)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", user.Login).
			Msg("unable to marshal user info to bson")
		return err
	}

	r.logger.Debug().Str("user", user.Login).Msg("inserting new user to the database")
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", user.Login).
			Msg("unable to insert new user to the database")
		return err
	}

	r.logger.Debug().Str("user", user.Login).Msgf(
		"new user was inserted to the database with %v",
		result.InsertedID,
	)
	return nil
}
