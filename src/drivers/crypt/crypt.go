package crypt

import (
	"encoding/hex"

	"github.com/wakuwaku3/account-book.api/src/application"
	"golang.org/x/crypto/scrypt"
)

type (
	c struct {
		env application.Env
	}
)

// NewCrypt は インスタンスを生成します
func NewCrypt(env application.Env) application.Crypt {
	return &c{env: env}
}
func (c *c) Hash(text *string) *string {
	salt := c.env.GetPasswordHashedKey()
	converted, _ := scrypt.Key([]byte(*text), *salt, 32768, 8, 1, 32)
	hashed := hex.EncodeToString(converted[:])
	return &hashed
}
