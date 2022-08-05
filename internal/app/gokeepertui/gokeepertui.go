package gokeepertui

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/marcusolsson/tui-go"
	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	g "github.com/serjyuriev/yandex-diploma-2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logo = `                                     
               _                       
   ___ ___ ___| |_ ___ ___ ___ ___ ___ 
  | . | . |___| '_| -_| -_| . | -_|  _|
  |_  |___|   |_,_|___|___|  _|___|_|  
  |___|                   |_|          `

// Client holds app's client-side related objects.
type Client struct {
	cfg    config.ClientConfig
	rpc    g.GokeeperClient
	ui     tui.UI
	logger zerolog.Logger
	userID string
	user   *g.User
}

// NewClient initializes app's client.
func NewClient() (*Client, error) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	logger.Debug().Msg("initializing go-keeper client")

	logger.Debug().Msg("getting app's configuration")
	cfg := config.GetClientConfig()

	logger.Debug().Msg("creating gRPC client")
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.
			Err(err).
			Caller().
			Msg("unable to connect to go-keeper server")
		return nil, err
	}
	rpcClient := g.NewGokeeperClient(conn)

	logger.Debug().Msg("drawing TUI")
	clt := new(Client)
	window, err := clt.DrawAuthWindow()
	if err != nil {
		logger.
			Err(err).
			Caller().
			Msg("unable to initialize go-keeper tui")
		return nil, err
	}

	clt.cfg = cfg
	clt.rpc = rpcClient
	clt.ui = window
	clt.logger = logger

	logger.Info().Msg("go-keeper client was successfully initialized")
	return clt, nil
}

// Start launches app's client.
func (c *Client) Start() error {
	if err := c.ui.Run(); err != nil {
		panic(err)
	}
	return nil
}

func (c *Client) SignUpUser(ctx context.Context, login, password string) (string, error) {
	user := &g.User{
		Login:    login,
		Password: password,
	}
	resp, err := c.rpc.SignUpUser(ctx, &g.SignUpUserRequest{User: user})
	if err != nil {
		c.logger.
			Err(err).
			Caller().
			Msg("unable to sign user up")
		return "", err
	}
	if resp.Error != "" {
		c.logger.
			Error().
			Caller().
			Msg(resp.Error)
		return "", errors.New(resp.Error)
	}
	return resp.UserID, nil
}

func (c *Client) LoginUser(ctx context.Context, login, password string) (string, error) {
	user := &g.User{
		Login:    login,
		Password: password,
	}
	resp, err := c.rpc.LoginUser(ctx, &g.LoginUserRequest{User: user})
	if err != nil {
		c.logger.
			Err(err).
			Caller().
			Msg("unable to login user")
		return "", err
	}
	if resp.Error != "" {
		c.logger.
			Error().
			Caller().
			Msg(resp.Error)
		return "", errors.New(resp.Error)
	}
	return resp.UserID, nil
}

// func (c *Client) Sync() {
// 	stream, err := c.rpc.Sync(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}
// 	go func() {
// 		for {
// 			in, err := stream.Recv()
// 			if err == io.EOF {
// 				return
// 			}
// 			if err != nil {
// 				panic(err)
// 			}
// 			c.user = in.User
// 		}
// 	}()
// 	if err = stream.Send(&g.SyncRequest{UserID: c.userID}); err != nil {
// 		panic(err)
// 	}
// 	stream.CloseSend()
// }

func (c *Client) AddLoginItem(ctx context.Context, login, password, userID string, meta map[string]string) error {
	loginItem := &g.LoginItem{
		Login:    login,
		Password: password,
		Meta:     meta,
	}
	resp, err := c.rpc.AddLoginItem(ctx, &g.AddLoginItemRequest{
		Item:   loginItem,
		UserID: userID,
	})
	if err != nil {
		c.logger.
			Err(err).
			Caller().
			Msg("unable to add login item")
		return err
	}
	if resp.Error != "" {
		c.logger.
			Error().
			Caller().
			Msg(resp.Error)
		return errors.New(resp.Error)
	}
	return nil
}

func (c *Client) DrawAuthWindow() (tui.UI, error) {
	user := tui.NewEntry()
	user.SetFocused(true)

	password := tui.NewEntry()
	password.SetEchoMode(tui.EchoModePassword)

	form := tui.NewGrid(0, 0)
	form.AppendRow(tui.NewLabel("User"), tui.NewLabel("Password"))
	form.AppendRow(user, password)

	status := tui.NewStatusBar("Not logged in.")

	login := tui.NewButton("[Login]")
	login.OnActivated(func(b *tui.Button) {
		id, err := c.LoginUser(context.TODO(), user.Text(), password.Text())
		if err != nil {
			status.SetText(fmt.Sprintf("Unable to login: %s", err.Error()))
			return
		}
		status.SetText(fmt.Sprintf("Logged in, userID = %s", id))
		c.userID = id
		c.rpc.UpdateItems(context.Background(), &g.UpdateItemsRequest{UserID: c.userID})
		c.DrawMenuWindow()
	})

	register := tui.NewButton("[Sign Up]")
	register.OnActivated(func(*tui.Button) {
		id, err := c.SignUpUser(context.TODO(), user.Text(), password.Text())
		if err != nil {
			status.SetText(fmt.Sprintf("Unable to sign up: %s", err.Error()))
			return
		}
		status.SetText(fmt.Sprintf("Signed up, userID = %s", id))
	})

	addLogin := tui.NewButton("[Add Login]")
	addLogin.OnActivated(func(*tui.Button) {
		if err := c.AddLoginItem(context.TODO(), user.Text(), password.Text(), c.userID, map[string]string{
			"testing": "onetwothree",
		}); err != nil {
			status.SetText(fmt.Sprintf("Unable to add new item: %s", err.Error()))
			return
		}
		status.SetText("Login item was added")
	})

	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(1, 0, login),
		tui.NewPadder(1, 0, register),
	)

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(logo)),
		tui.NewPadder(12, 0, tui.NewLabel("Welcome to go-keeper! Login or sign up.")),
		tui.NewPadder(1, 1, form),
		buttons,
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(
		content,
		status,
	)

	tui.DefaultFocusChain.Set(user, password, login, register)

	ui, err := tui.New(root)
	if err != nil {
		c.logger.
			Err(err).
			Caller().
			Msg("unable to create login window")
		return nil, err
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	return ui, nil
}

func (c *Client) DrawMenuWindow() error {
	c.ui.Quit()
	menu := tui.NewTable(0, 0)
	menu.SetColumnStretch(0, 1)

	logins := tui.NewButton("[Logins]")
	logins.OnActivated(func(b *tui.Button) {
		c.DrawLoginItemsWindow()
	})
	cards := tui.NewButton("[Cards]")
	cards.OnActivated(func(b *tui.Button) {
	})
	texts := tui.NewButton("[Texts]")
	texts.OnActivated(func(b *tui.Button) {
	})
	binaries := tui.NewButton("[Binaries]")
	binaries.OnActivated(func(b *tui.Button) {
	})

	menu.AppendRow(logins)
	menu.AppendRow(cards)
	menu.AppendRow(texts)
	menu.AppendRow(binaries)
	menu.SetFocused(true)

	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(logo)),
		tui.NewPadder(12, 0, tui.NewLabel("Welcome to go-keeper! Select item type.")),
		tui.NewPadder(1, 1, menu),
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(content)

	ui, err := tui.New(root)
	if err != nil {
		c.logger.
			Err(err).
			Caller().
			Msg("unable to create menu window")
		return err
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	c.ui = ui
	c.Start()

	return nil
}

func (c *Client) DrawLoginItemsWindow() error {
	c.ui.Quit()

	inbox := tui.NewTable(0, 0)
	inbox.SetColumnStretch(0, 0)
	inbox.SetFocused(true)

	for _, l := range c.user.Logins {
		inbox.AppendRow(
			tui.NewLabel(l.Login),
		)
	}

	var (
		login    = tui.NewLabel("")
		password = tui.NewLabel("")
		meta     = tui.NewLabel("")
	)

	info := tui.NewGrid(0, 0)
	info.AppendRow(tui.NewLabel("Login:"), login)
	info.AppendRow(tui.NewLabel("Password:"), password)
	info.AppendRow(tui.NewLabel("Meta:"), meta)

	body := tui.NewLabel("")
	body.SetSizePolicy(tui.Preferred, tui.Expanding)

	mail := tui.NewVBox(info, body)
	mail.SetSizePolicy(tui.Preferred, tui.Expanding)

	inbox.OnSelectionChanged(func(t *tui.Table) {
		l := c.user.Logins[t.Selected()]
		login.SetText(l.Login)
		password.SetText(l.Password)
		var metaString string
		for k, v := range l.Meta {
			metaString += fmt.Sprintf("%s: %s\n", k, v)
		}
		meta.SetText(metaString)
	})

	// Select first mail on startup.
	inbox.Select(0)

	root := tui.NewVBox(inbox, tui.NewLabel(""), mail)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	c.ui = ui
	c.Start()

	return nil
}
