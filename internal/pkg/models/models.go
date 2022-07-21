package models

// LoginPasswordItem holds information about
// single login-password entry.
type LoginPasswordItem struct {
	Login    string
	Password string
	Meta     map[string]string
}

// BankCardItem holds bank card related information.
type BankCardItem struct {
	Number           string
	Holder           string
	Expires          string
	CardSecurityCode int
	Meta             map[string]string
}

// TextItem holds arbitrary text information.
type TextItem struct {
	Value string
	Meta  map[string]string
}

