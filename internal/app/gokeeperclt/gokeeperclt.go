package gokeeperclt

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/serjyuriev/yandex-diploma-2/internal/pkg/config"
	g "github.com/serjyuriev/yandex-diploma-2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client holds app's client-side related objects.
type Client struct {
	cfg    config.ClientConfig
	rpc    g.GokeeperClient
	logger zerolog.Logger
	mode   *mode
	user   *g.User
}

// mode stores all flag values.
type mode struct {
	SignUp         bool
	GetLoginItems  bool
	GetCardItems   bool
	GetTextItems   bool
	GetBinaryItems bool
	AddLoginItem   bool
	AddCardItem    bool
	AddTextItem    bool
	AddBinaryItem  bool
}

// NewClient initializes app's client.
func NewClient() (*Client, error) {
	mode := &mode{}
	user := &g.User{}
	flag.BoolVar(&mode.SignUp, "signup", false, "sign up as new user")
	flag.StringVar(&user.Login, "login", "", "user login")
	flag.StringVar(&user.Password, "password", "", "user password")
	flag.BoolVar(&mode.GetLoginItems, "lp", false, "get login-password items")
	flag.BoolVar(&mode.GetCardItems, "bc", false, "get bank card items")
	flag.BoolVar(&mode.GetTextItems, "text", false, "get text items")
	flag.BoolVar(&mode.GetBinaryItems, "bins", false, "get binary items")
	flag.BoolVar(&mode.AddLoginItem, "alp", false, "add login-password item")
	flag.BoolVar(&mode.AddCardItem, "abc", false, "add bank card item")
	flag.BoolVar(&mode.AddTextItem, "atext", false, "add text item")
	flag.BoolVar(&mode.AddBinaryItem, "abins", false, "add binary item")

	cfg := config.GetClientConfig()

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02-01-2006 15:04:05 MST",
	}
	var level zerolog.Level
	if cfg.IsDebug {
		level = zerolog.DebugLevel
	} else {
		level = zerolog.ErrorLevel
	}
	logger := zerolog.New(output).With().Timestamp().Logger().Level(level)

	logger.Debug().Msg("initializing go-keeper client")

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

	logger.Debug().Msg("parsing flags")
	flag.Parse()

	logger.Info().Msg("go-keeper client was successfully initialized")
	return &Client{
		cfg:    cfg,
		rpc:    rpcClient,
		logger: logger,
		mode:   mode,
		user:   user,
	}, nil
}

func (c *Client) Run() error {
	if c.user.Login == "" || c.user.Password == "" {
		return fmt.Errorf("login and/or password cannot be empty")
	}
	if c.mode.SignUp {
		userID, err := c.signUpUser(context.Background(), c.user.Login, c.user.Password)
		if err != nil {
			c.logger.
				Err(err).
				Caller().
				Msg("unable to sign up user")
			return err
		}
		fmt.Printf("successfully signed up, your user id is %s\n", userID)
	} else {
		userID, err := c.loginUser(context.Background(), c.user.Login, c.user.Password)
		if err != nil {
			c.logger.
				Err(err).
				Caller().
				Msg("unable to log user in")
			return err
		}
		fmt.Printf("successfully logged in, your user id is %s\n", userID)
		if err = c.updateItems(context.Background(), userID); err != nil {
			c.logger.
				Err(err).
				Caller().
				Msg("unable to update user's items")
			return err
		}
		fmt.Println("updated your items")
		if c.mode.GetLoginItems {
			c.displayLoginItems()
		}
		if c.mode.GetCardItems {
			c.displayCardItems()
		}
		if c.mode.GetTextItems {
			c.displayTextItems()
		}
		if c.mode.GetBinaryItems {
			c.displayBinaryItems()
		}
		if c.mode.AddLoginItem {
			item, err := c.getLoginItemFromUser()
			if err != nil {
				c.logger.Err(err).Caller().Msg("unable to get login item from user")
				return err
			}
			if err = c.addLoginItem(context.Background(), item, userID); err != nil {
				c.logger.Err(err).Caller().Msg("unable to add new login item")
			}
		}
		if c.mode.AddCardItem {
			item, err := c.getCardItemFromUser()
			if err != nil {
				c.logger.Err(err).Caller().Msg("unable to get card item from user")
				return err
			}
			if err = c.addCardItem(context.Background(), item, userID); err != nil {
				c.logger.Err(err).Caller().Msg("unable to add new card item")
			}
		}
		if c.mode.AddTextItem {
			item, err := c.getTextItemFromUser()
			if err != nil {
				c.logger.Err(err).Caller().Msg("unable to get text item from user")
				return err
			}
			if err = c.addTextItem(context.Background(), item, userID); err != nil {
				c.logger.Err(err).Caller().Msg("unable to add new text item")
			}
		}
		if c.mode.AddBinaryItem {
			item, err := c.getBinaryItemFromUser()
			if err != nil {
				c.logger.Err(err).Caller().Msg("unable to get binary item from user")
				return err
			}
			if err = c.addBinaryItem(context.Background(), item, userID); err != nil {
				c.logger.Err(err).Caller().Msg("unable to add new binary item")
			}
		}
	}
	return nil
}

func (c *Client) signUpUser(ctx context.Context, login, password string) (string, error) {
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

func (c *Client) loginUser(ctx context.Context, login, password string) (string, error) {
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

func (c *Client) updateItems(ctx context.Context, userID string) error {
	resp, err := c.rpc.UpdateItems(ctx, &g.UpdateItemsRequest{UserID: userID})
	if err != nil {
		c.logger.
			Err(err).
			Caller().
			Msg("unable to perform grpc request")
		return err
	}
	if resp.Error != "" {
		c.logger.
			Error().
			Caller().
			Msg(resp.Error)
		return errors.New(resp.Error)
	}
	c.user = resp.User
	return nil
}

func (c *Client) addLoginItem(ctx context.Context, item *g.LoginItem, userID string) error {
	req := &g.AddLoginItemRequest{
		Item:   item,
		UserID: userID,
	}
	resp, err := c.rpc.AddLoginItem(ctx, req)
	if err != nil {
		c.logger.Err(err).Caller().Msg("unable to perform rpc request")
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

func (c *Client) addCardItem(ctx context.Context, item *g.BankCardItem, userID string) error {
	req := &g.AddBankCardItemRequest{
		Item:   item,
		UserID: userID,
	}
	resp, err := c.rpc.AddBankCardItem(ctx, req)
	if err != nil {
		c.logger.Err(err).Caller().Msg("unable to perform rpc request")
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

func (c *Client) addTextItem(ctx context.Context, item *g.TextItem, userID string) error {
	req := &g.AddTextItemRequest{
		Item:   item,
		UserID: userID,
	}
	resp, err := c.rpc.AddTextItem(ctx, req)
	if err != nil {
		c.logger.Err(err).Caller().Msg("unable to perform rpc request")
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

func (c *Client) addBinaryItem(ctx context.Context, item *g.BinaryItem, userID string) error {
	req := &g.AddBinaryItemRequest{
		Item:   item,
		UserID: userID,
	}
	resp, err := c.rpc.AddBinaryItem(ctx, req)
	if err != nil {
		c.logger.Err(err).Caller().Msg("unable to perform rpc request")
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

func (c *Client) getLoginItemFromUser() (*g.LoginItem, error) {
	sc := bufio.NewScanner(os.Stdin)
	item := &g.LoginItem{}
	fmt.Println("Login:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Login = sc.Text()

	fmt.Println("Password:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Password = sc.Text()

	fmt.Println("Meta (leave field empty to stop):")
	fmt.Println()
	item.Meta = make(map[string]string)
	for {
		fmt.Println("Key:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		key := sc.Text()
		if key == "" {
			break
		}

		fmt.Println("Value:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		val := sc.Text()
		if val == "" {
			break
		}
		item.Meta[key] = val
	}
	return item, nil
}

func (c *Client) getCardItemFromUser() (*g.BankCardItem, error) {
	sc := bufio.NewScanner(os.Stdin)
	item := &g.BankCardItem{}
	fmt.Println("Holder:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Holder = sc.Text()

	fmt.Println("Number:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Number = sc.Text()

	fmt.Println("Expires:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Expires = sc.Text()

	fmt.Println("Security code:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	code, err := strconv.Atoi(sc.Text())
	if err != nil {
		c.logger.Err(err).Caller().Msg("unable to parse security code")
		return nil, err
	}
	item.CardSecurityCode = int32(code)

	fmt.Println("Meta (leave field empty to stop):")
	fmt.Println()
	item.Meta = make(map[string]string)
	for {
		fmt.Println("Key:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		key := sc.Text()
		if key == "" {
			break
		}

		fmt.Println("Value:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		val := sc.Text()
		if val == "" {
			break
		}
		item.Meta[key] = val
	}
	return item, nil
}

func (c *Client) getTextItemFromUser() (*g.TextItem, error) {
	sc := bufio.NewScanner(os.Stdin)
	item := &g.TextItem{}
	fmt.Println("Text:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Value = sc.Text()

	fmt.Println("Meta (leave field empty to stop):")
	fmt.Println()
	item.Meta = make(map[string]string)
	for {
		fmt.Println("Key:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		key := sc.Text()
		if key == "" {
			break
		}

		fmt.Println("Value:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		val := sc.Text()
		if val == "" {
			break
		}
		item.Meta[key] = val
	}
	return item, nil
}

func (c *Client) getBinaryItemFromUser() (*g.BinaryItem, error) {
	sc := bufio.NewScanner(os.Stdin)
	item := &g.BinaryItem{}
	fmt.Println("Binary:")
	sc.Scan()
	if sc.Err() != nil {
		c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
		return nil, sc.Err()
	}
	item.Value = []byte(sc.Text())

	fmt.Println("Meta (leave field empty to stop):")
	fmt.Println()
	item.Meta = make(map[string]string)
	for {
		fmt.Println("Key:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		key := sc.Text()
		if key == "" {
			break
		}

		fmt.Println("Value:")
		sc.Scan()
		if sc.Err() != nil {
			c.logger.Err(sc.Err()).Caller().Msg("unable to scan user input")
			return nil, sc.Err()
		}
		val := sc.Text()
		if val == "" {
			break
		}
		item.Meta[key] = val
	}
	return item, nil
}

func (c *Client) displayLoginItems() {
	fmt.Println("\n---------------- LOGINS ----------------")
	if len(c.user.Logins) == 0 {
		fmt.Println("there are no login items yet")
		return
	}
	for _, item := range c.user.Logins {
		fmt.Printf("Login: %s\n", item.Login)
		fmt.Printf("Password: %s\n", item.Password)
		fmt.Println("Meta:")
		for k, v := range item.Meta {
			fmt.Printf("\t%s: %v\n", k, v)
		}
		fmt.Println("----------------------------------------")
	}
	fmt.Println()
}

func (c *Client) displayCardItems() {
	fmt.Println("\n---------------- CARDS ----------------")
	if len(c.user.Cards) == 0 {
		fmt.Println("there are no card items yet")
		return
	}
	for _, item := range c.user.Cards {
		fmt.Printf("Holder: %s\n", item.Holder)
		fmt.Printf("Number: %s\n", item.Number)
		fmt.Printf("Expires: %s\n", item.Expires)
		fmt.Printf("Security code: %d\n", item.CardSecurityCode)
		fmt.Println("Meta:")
		for k, v := range item.Meta {
			fmt.Printf("\t%s: %v\n", k, v)
		}
		fmt.Println("---------------------------------------")
	}
	fmt.Println()
}

func (c *Client) displayTextItems() {
	fmt.Println("\n---------------- TEXTS ----------------")
	if len(c.user.Texts) == 0 {
		fmt.Println("there are no text items yet")
		return
	}
	for _, item := range c.user.Texts {
		fmt.Printf("Text: %s\n", item.Value)
		fmt.Println("Meta:")
		for k, v := range item.Meta {
			fmt.Printf("\t%s: %v\n", k, v)
		}
		fmt.Println("---------------------------------------")
	}
	fmt.Println()
}

func (c *Client) displayBinaryItems() {
	fmt.Println("\n---------------- BINARIES ----------------")
	if len(c.user.Binaries) == 0 {
		fmt.Println("there are no binary items yet")
		return
	}
	for _, item := range c.user.Binaries {
		fmt.Printf("Binary data: %s\n", item.Value)
		fmt.Println("Meta:")
		for k, v := range item.Meta {
			fmt.Printf("\t%s: %v\n", k, v)
		}
		fmt.Println("------------------------------------------")
	}
	fmt.Println()
}