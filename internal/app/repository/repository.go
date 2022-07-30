package repository

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoUser = errors.New("there is no such user in the database")
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

// ReadUserByLogin searches the database for a user
// with provided login, returning found user or ErrNoUser.
func (r *Repository) ReadUserByLogin(ctx context.Context, login string) (*models.User, error) {
	r.logger.Debug().Str("user", login).Msg("getting users collection")
	collection := r.client.Database(r.cfg.Database.Name).Collection("users")

	r.logger.Debug().Str("user", login).Msg("preparing filter")
	filter := bson.D{{Key: "login", Value: login}}

	r.logger.Debug().Str("user", login).Msg("searching for user in the database")
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			r.logger.Debug().Str("user", login).Msg("no such user in the database")
			return nil, ErrNoUser
		}
		r.logger.
			Err(result.Err()).
			Caller().
			Str("user", login).
			Msg("unable to perform read operation in the database")
		return nil, result.Err()
	}

	r.logger.Debug().Str("user", login).Msg("processing query result")
	var user models.User
	if err := result.Decode(&user); err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", login).
			Msg("unable to decode query result")
		return nil, err
	}

	r.logger.Debug().Str("user", login).Msg("user was found in the database")
	return &user, nil
}
