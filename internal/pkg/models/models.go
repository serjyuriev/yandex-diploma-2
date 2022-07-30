package models

import (
	"github.com/google/uuid"
)

// User holds information about app's user.
type User struct {
	ID        uuid.UUID            `json:"id" bson:"id"`
	Login     string               `json:"login" bson:"login"`
	Password  string               `json:"password" bson:"password"`
	Logins    []*LoginPasswordItem `json:"logins" bson:"logins"`
	BankCards []*BankCardItem      `json:"bank_cards" bson:"bank_cards"`
	Texts     []*TextItem          `json:"texts" bson:"texts"`
	Binaries  []*BinaryItem        `json:"binaries" bson:"binaries"`
}

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

// BinaryItem holds arbitrary binary information.
type BinaryItem struct {
	Value []byte
	Meta  map[string]string
}
