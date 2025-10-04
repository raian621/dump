package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "MySuperSecurePassword"
	hash := HashPassword(password)
	assert.Equal(t, len(hash), 100)
}

func TestValidateHash(t *testing.T) {
	password := "MySuperSecurePassword"
	hash := HashPassword(password)
	assert.True(t, ValidatePassword(password, hash))
	assert.False(t, ValidatePassword("NotMySuperSecurePassword", hash))
}
