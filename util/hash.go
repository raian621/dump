package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	ARGON2ID_TIME        = 1
	ARGON2ID_MEMORY      = 64 * 1024
	ARGON2ID_THREADS     = 4
	ARGON2ID_KEY_LENGTH  = 32
	ARGON2ID_SALT_LENGTH = 16
)

func HashPassword(password string) string {
	return argon2idHash(password, generateRandomSalt())
}

func argon2idHash(password string, salt []byte) string {
	key := argon2.IDKey(
		[]byte(password), salt, ARGON2ID_TIME, ARGON2ID_MEMORY, ARGON2ID_THREADS,
		ARGON2ID_KEY_LENGTH)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, ARGON2ID_MEMORY,
		ARGON2ID_TIME, ARGON2ID_THREADS, base64.URLEncoding.EncodeToString(salt),
		base64.URLEncoding.EncodeToString(key))
}

func generateRandomSalt() []byte {
	salt := make([]byte, ARGON2ID_SALT_LENGTH)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	return salt
}

func ValidatePassword(password string, hash string) bool {
	fields := strings.Split(hash, "$")
	saltBase64 := fields[len(fields)-2]
	salt := make([]byte, ARGON2ID_SALT_LENGTH)
	if _, err := base64.URLEncoding.Decode(salt, []byte(saltBase64)); err != nil {
		panic(err)
	}
	hash2 := argon2idHash(password, salt)
	return hash2 == hash
}
