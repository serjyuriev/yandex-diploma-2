// Package models provides items and user models for gokeeper app.
package models

import (
	"github.com/google/uuid"
)

// User holds information about app's user.
type User struct {
	ID        uuid.UUID            `bson:"id"`
	Login     string               `bson:"login"`
	Password  string               `bson:"password"`
	Logins    []*LoginPasswordItem `bson:"logins"`
	BankCards []*BankCardItem      `bson:"cards"`
	Texts     []*TextItem          `bson:"texts"`
	Binaries  []*BinaryItem        `bson:"binaries"`
}

// LoginPasswordItem holds information about
// single login-password entry.
type LoginPasswordItem struct {
	Login    string            `bson:"login"`
	Password string            `bson:"password"`
	Meta     map[string]string `bson:"meta"`
}

// BankCardItem holds bank card related information.
type BankCardItem struct {
	Number           string            `bson:"number"`
	Holder           string            `bson:"holder"`
	Expires          string            `bson:"expires"`
	CardSecurityCode int               `bson:"csc"`
	Meta             map[string]string `bson:"meta"`
}

// TextItem holds arbitrary text information.
type TextItem struct {
	Value string            `bson:"value"`
	Meta  map[string]string `bson:"meta"`
}

// BinaryItem holds arbitrary binary information.
type BinaryItem struct {
	Value []byte            `bson:"value"`
	Meta  map[string]string `bson:"meta"`
}
