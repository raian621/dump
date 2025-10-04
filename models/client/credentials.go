package client

import (
	"github.com/raian621/dump/models/storage"
	"github.com/raian621/dump/util"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Credentials) ToStorageModel() *storage.Credentials {
	return &storage.Credentials{
		Username: c.Username,
		Passhash: util.HashPassword(c.Password),
	}
}
