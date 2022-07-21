package models

// LoginPasswordItem holds information about
// single login-password entry.
type LoginPasswordItem struct {
	Login    string
	Password string
	Meta     map[string]string
}

