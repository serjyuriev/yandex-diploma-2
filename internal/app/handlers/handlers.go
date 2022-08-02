package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/serjyuriev/yandex-diploma-2/internal/app/repository"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/service"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	g "github.com/serjyuriev/yandex-diploma-2/proto"
)

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
		Login:     in.User.Login,
		Password:  in.User.Password,
		Logins:    make([]*models.LoginPasswordItem, 0),
		BankCards: make([]*models.BankCardItem, 0),
		Texts:     make([]*models.TextItem, 0),
		Binaries:  make([]*models.BinaryItem, 0),
	}
	res := new(g.SignUpUserResponse)

	r.logger.Debug().Str("user", in.User.Login).Msg("passing user's info to service layer")
	userID, err := r.svc.SignUpUser(ctx, user)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", user.Login).
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

// LoginUser logins existing user.
func (r *RPC) LoginUser(ctx context.Context, in *g.LoginUserRequest) (*g.LoginUserResponse, error) {
	r.logger.Info().Str("user", in.User.Login).Msg("received user login request")
	user := &models.User{
		Login:    in.User.Login,
		Password: in.User.Password,
	}
	res := new(g.LoginUserResponse)

	r.logger.Debug().Str("user", in.User.Login).Msg("passing user's info to service layer")
	userID, err := r.svc.LoginUser(ctx, user)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", user.Login).
			Msg("unable to login user")
		res.UserID = ""
		res.Error = err.Error()
		return res, err
	}

	r.logger.Info().Str("user", in.User.Login).Msg("user was successfully logged in")
	res.UserID = userID
	res.Error = ""
	return res, nil
}

// AddLoginItem adds new login entry in the user's vault.
func (r *RPC) AddLoginItem(ctx context.Context, in *g.AddLoginItemRequest) (*g.AddLoginItemResponse, error) {
	r.logger.Info().Str("user", in.UserID).Msg("received new login item")
	login := &models.LoginPasswordItem{
		Login:    in.Item.Login,
		Password: in.Item.Password,
		Meta:     in.Item.Meta,
	}
	res := new(g.AddLoginItemResponse)

	r.logger.Debug().Str("user", in.UserID).Msg("parsing user uuid")
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to parse user uuid")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Debug().Str("user", in.UserID).Msg("passing new login item to data layer")
	if err := r.repo.CreateItem(ctx, login, "logins", userID); err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to create new login item")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Info().Str("user", in.UserID).Msg("login item was successfully added")
	res.Error = ""
	return res, nil
}

// AddBankCardItem adds new bank card entry in the user's vault.
func (r *RPC) AddBankCardItem(ctx context.Context, in *g.AddBankCardItemRequest) (*g.AddBankCardItemResponse, error) {
	r.logger.Info().Str("user", in.UserID).Msg("received new bank card item")
	card := &models.BankCardItem{
		Number:           in.Item.Number,
		Holder:           in.Item.Holder,
		Expires:          in.Item.Expires,
		CardSecurityCode: int(in.Item.CardSecurityCode),
		Meta:             in.Item.Meta,
	}
	res := new(g.AddBankCardItemResponse)

	r.logger.Debug().Str("user", in.UserID).Msg("parsing user uuid")
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to parse user uuid")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Debug().Str("user", in.UserID).Msg("passing new bank card item to data layer")
	if err := r.repo.CreateItem(ctx, card, "cards", userID); err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to create new bank card item")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Info().Str("user", in.UserID).Msg("bank card item was successfully added")
	res.Error = ""
	return res, nil
}

// AddTextItem adds new text entry in the user's vault.
func (r *RPC) AddTextItem(ctx context.Context, in *g.AddTextItemRequest) (*g.AddTextItemResponse, error) {
	r.logger.Info().Str("user", in.UserID).Msg("received new text item")
	text := &models.TextItem{
		Value: in.Item.Value,
		Meta:  in.Item.Meta,
	}
	res := new(g.AddTextItemResponse)

	r.logger.Debug().Str("user", in.UserID).Msg("parsing user uuid")
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to parse user uuid")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Debug().Str("user", in.UserID).Msg("passing new text item to data layer")
	if err := r.repo.CreateItem(ctx, text, "texts", userID); err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to create new text item")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Info().Str("user", in.UserID).Msg("text item was successfully added")
	res.Error = ""
	return res, nil
}

// AddBinaryItem adds new binary entry in the user's vault.
func (r *RPC) AddBinaryItem(ctx context.Context, in *g.AddBinaryItemRequest) (*g.AddBinaryItemResponse, error) {
	r.logger.Info().Str("user", in.UserID).Msg("received new binary item")
	bin := &models.BinaryItem{
		Value: in.Item.Value,
		Meta:  in.Item.Meta,
	}
	res := new(g.AddBinaryItemResponse)

	r.logger.Debug().Str("user", in.UserID).Msg("parsing user uuid")
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to parse user uuid")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Debug().Str("user", in.UserID).Msg("passing new binary item to data layer")
	if err := r.repo.CreateItem(ctx, bin, "binaries", userID); err != nil {
		r.logger.
			Err(err).
			Caller().
			Str("user", in.UserID).
			Msg("unable to create new binary item")
		res.Error = err.Error()
		return res, err
	}

	r.logger.Info().Str("user", in.UserID).Msg("binary item was successfully added")
	res.Error = ""
	return res, nil
}
