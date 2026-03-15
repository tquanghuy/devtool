package config

import (
	"github.com/zalando/go-keyring"
)

const serviceName = "devtool"

// SetPassword saves a password for a given user in the system keyring.
func SetPassword(user string, password string) error {
	return keyring.Set(serviceName, user, password)
}

// GetPassword retrieves a password for a given user from the system keyring.
func GetPassword(user string) (string, error) {
	return keyring.Get(serviceName, user)
}
