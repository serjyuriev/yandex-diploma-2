package handlers

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/app/mocks"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/models"
	g "github.com/serjyuriev/yandex-diploma-2/proto"
	"github.com/stretchr/testify/require"
)

func TestMakeRPC(t *testing.T) {
	os.Args = []string{"test", "-c", "A:\\dev\\yandex\\yandex-diploma-2\\dev_srv_config.yaml"}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	_, err := MakeRPC(logger)
	require.NoError(t, err)
}

func TestSignUpUser(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	ms := mocks.NewMockService(ctrl)

	t.Run("success", func(t *testing.T) {
		in := &g.SignUpUserRequest{
			User: &g.User{
				Login:    "test",
				Password: "somepwd",
			},
		}

		sign := ms.EXPECT().
			SignUpUser(context.Background(), gomock.Any()).
			Return("uid", nil)
		gomock.InOrder(sign)

		rpc := &RPC{
			logger: logger,
			svc:    ms,
		}
		out, err := rpc.SignUpUser(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "uid", out.UserID)
	})

	t.Run("sign err", func(t *testing.T) {
		in := &g.SignUpUserRequest{
			User: &g.User{
				Login:    "test",
				Password: "somepwd",
			},
		}

		sign := ms.EXPECT().
			SignUpUser(context.Background(), gomock.Any()).
			Return("", fmt.Errorf("some err"))
		gomock.InOrder(sign)

		rpc := &RPC{
			logger: logger,
			svc:    ms,
		}
		_, err := rpc.SignUpUser(context.Background(), in)
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
	ms := mocks.NewMockService(ctrl)

	t.Run("success", func(t *testing.T) {
		in := &g.LoginUserRequest{
			User: &g.User{
				Login:    "test",
				Password: "somepwd",
			},
		}

		login := ms.EXPECT().
			LoginUser(context.Background(), gomock.Any()).
			Return("uid", nil)
		gomock.InOrder(login)

		rpc := &RPC{
			logger: logger,
			svc:    ms,
		}
		out, err := rpc.LoginUser(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "uid", out.UserID)
	})

	t.Run("login err", func(t *testing.T) {
		in := &g.LoginUserRequest{
			User: &g.User{
				Login:    "test",
				Password: "somepwd",
			},
		}

		login := ms.EXPECT().
			LoginUser(context.Background(), gomock.Any()).
			Return("", fmt.Errorf("some err"))
		gomock.InOrder(login)

		rpc := &RPC{
			logger: logger,
			svc:    ms,
		}
		_, err := rpc.LoginUser(context.Background(), in)
		require.Error(t, err)
	})
}

func TestUpdateItems(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		in := &g.UpdateItemsRequest{
			UserID: uid.String(),
		}
		expectedUser := &g.User{
			Login: "test",
			Logins: []*g.LoginItem{
				{
					Login:    "one",
					Password: "two",
					Meta:     nil,
				},
			},
			Cards: []*g.BankCardItem{
				{
					Number:           "",
					Holder:           "",
					Expires:          "",
					CardSecurityCode: 123,
					Meta:             nil,
				},
			},
			Texts: []*g.TextItem{
				{
					Value: "some text",
					Meta:  nil,
				},
			},
			Binaries: []*g.BinaryItem{
				{
					Value: []byte("some bins"),
					Meta:  nil,
				},
			},
		}
		dbUser := &models.User{
			Login: "test",
			Logins: []*models.LoginPasswordItem{
				{
					Login:    "one",
					Password: "two",
					Meta:     nil,
				},
			},
			BankCards: []*models.BankCardItem{
				{
					Number:           "",
					Holder:           "",
					Expires:          "",
					CardSecurityCode: 123,
					Meta:             nil,
				},
			},
			Texts: []*models.TextItem{
				{
					Value: "some text",
					Meta:  nil,
				},
			},
			Binaries: []*models.BinaryItem{
				{
					Value: []byte("some bins"),
					Meta:  nil,
				},
			},
		}

		read := mr.EXPECT().
			ReadUserByID(
				context.Background(),
				gomock.Eq(uid),
			).Return(dbUser, nil)
		gomock.InOrder(read)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		out, err := rpc.UpdateItems(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "", out.Error)
		require.Equal(t, expectedUser, out.User)
	})

	t.Run("read err", func(t *testing.T) {
		uid := uuid.New()
		in := &g.UpdateItemsRequest{
			UserID: uid.String(),
		}

		read := mr.EXPECT().
			ReadUserByID(
				context.Background(),
				gomock.Eq(uid),
			).Return(nil, fmt.Errorf("some err"))
		gomock.InOrder(read)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.UpdateItems(context.Background(), in)
		require.Error(t, err)
	})

	t.Run("wrong uid", func(t *testing.T) {
		in := &g.UpdateItemsRequest{
			UserID: "asdkjnfklasj lkmndjkl fna",
		}

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.UpdateItems(context.Background(), in)
		require.Error(t, err)
	})
}

func TestAddLoginItem(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddLoginItemRequest{
			Item: &g.LoginItem{
				Login:    "test",
				Password: "somepwd",
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"logins",
				gomock.Eq(uid),
			).Return(nil)
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		out, err := rpc.AddLoginItem(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "", out.Error)
	})

	t.Run("repo err", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddLoginItemRequest{
			Item: &g.LoginItem{
				Login:    "test",
				Password: "somepwd",
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"logins",
				gomock.Eq(uid),
			).Return(fmt.Errorf("some err"))
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddLoginItem(context.Background(), in)
		require.Error(t, err)
	})

	t.Run("wrong uuid", func(t *testing.T) {
		in := &g.AddLoginItemRequest{
			Item: &g.LoginItem{
				Login:    "test",
				Password: "somepwd",
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: "o413ngklj3 gk j324kg j4 ",
		}

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddLoginItem(context.Background(), in)
		require.Error(t, err)
	})
}

func TestAddBankCardItem(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddBankCardItemRequest{
			Item: &g.BankCardItem{
				Number:           "",
				Holder:           "",
				Expires:          "",
				CardSecurityCode: 12,
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"cards",
				gomock.Eq(uid),
			).Return(nil)
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		out, err := rpc.AddBankCardItem(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "", out.Error)
	})

	t.Run("repo err", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddBankCardItemRequest{
			Item: &g.BankCardItem{
				Number:           "",
				Holder:           "",
				Expires:          "",
				CardSecurityCode: 12,
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"cards",
				gomock.Eq(uid),
			).Return(fmt.Errorf("some err"))
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddBankCardItem(context.Background(), in)
		require.Error(t, err)
	})

	t.Run("wrong uuid", func(t *testing.T) {
		in := &g.AddBankCardItemRequest{
			Item: &g.BankCardItem{
				Number:           "",
				Holder:           "",
				Expires:          "",
				CardSecurityCode: 12,
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: "l 23knkl;jn t",
		}

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddBankCardItem(context.Background(), in)
		require.Error(t, err)
	})
}

func TestAddTextItem(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddTextItemRequest{
			Item: &g.TextItem{
				Value: "some text",
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"texts",
				gomock.Eq(uid),
			).Return(nil)
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		out, err := rpc.AddTextItem(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "", out.Error)
	})

	t.Run("repo err", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddTextItemRequest{
			Item: &g.TextItem{
				Value: "some text",
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"texts",
				gomock.Eq(uid),
			).Return(fmt.Errorf("some err"))
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddTextItem(context.Background(), in)
		require.Error(t, err)
	})

	t.Run("wrong uuid", func(t *testing.T) {
		in := &g.AddTextItemRequest{
			Item: &g.TextItem{
				Value: "some text",
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: "3ljnr k2jn3 ",
		}

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddTextItem(context.Background(), in)
		require.Error(t, err)
	})
}

func TestAddBinaryItem(t *testing.T) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	ctrl := gomock.NewController(t)
	mr := mocks.NewMockRepository(ctrl)

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddBinaryItemRequest{
			Item: &g.BinaryItem{
				Value: []byte("some text"),
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"binaries",
				gomock.Eq(uid),
			).Return(nil)
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		out, err := rpc.AddBinaryItem(context.Background(), in)
		require.NoError(t, err)
		require.Equal(t, "", out.Error)
	})

	t.Run("repo err", func(t *testing.T) {
		uid := uuid.New()
		in := &g.AddBinaryItemRequest{
			Item: &g.BinaryItem{
				Value: []byte("some text"),
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: uid.String(),
		}

		create := mr.EXPECT().
			CreateItem(
				context.Background(),
				gomock.Any(),
				"binaries",
				gomock.Eq(uid),
			).Return(fmt.Errorf("some err"))
		gomock.InOrder(create)

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddBinaryItem(context.Background(), in)
		require.Error(t, err)
	})

	t.Run("wrong uuid", func(t *testing.T) {
		in := &g.AddBinaryItemRequest{
			Item: &g.BinaryItem{
				Value: []byte("some text"),
				Meta: map[string]string{
					"one":   "two",
					"three": "four",
				},
			},
			UserID: "uid.String()",
		}

		rpc := &RPC{
			logger: logger,
			repo:   mr,
		}
		_, err := rpc.AddBinaryItem(context.Background(), in)
		require.Error(t, err)
	})
}
