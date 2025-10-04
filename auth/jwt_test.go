package auth

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateRandomSecret() []byte {
	secret := make([]byte, 256/8)
	if _, err := rand.Read(secret); err != nil {
		panic(err)
	}
	return secret
}

func TestCreateAccessToken(t *testing.T) {
	accessTtl := 20
	tf := NewTokenFactory(accessTtl, 40, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	parsedAccessToken, err := tf.ParseAccessToken(accessTokenStr)
	assert.NoError(t, err)
	assert.Equal(t, accessToken.Claims, parsedAccessToken.Claims.(*AccessTokenClaims))
}

func TestValidateExpiredAccessToken(t *testing.T) {
	tf := NewTokenFactory( /*accessTtl=*/ 0 /*refreshTtl=*/, 0, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	_, err = tf.ParseAccessToken(accessTokenStr)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestRefreshAccessToken(t *testing.T) {
	accessTtl := 1
	tf := NewTokenFactory(accessTtl /*refreshTtl=*/, 20, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	refreshToken := tf.CreateRefreshToken(1)
	refreshTokenStr, err := tf.SignedString(refreshToken)
	assert.NoError(t, err)
	time.Sleep(time.Second) // wait for access token to expire
	newAccessTokenStr, err := tf.RefreshAccessToken(accessTokenStr, refreshTokenStr)
	assert.NoError(t, err)
	newAccessToken, err := tf.ParseAccessToken(newAccessTokenStr)
	assert.NoError(t, err)
	issuedAt, err := newAccessToken.Claims.GetIssuedAt()
	assert.NoError(t, err)
	expiresAt, err := newAccessToken.Claims.GetExpirationTime()
	assert.NoError(t, err)
	assert.Equal(t, issuedAt.Time.Add(time.Second*time.Duration(accessTtl)), expiresAt.Time)
}
func TestRefreshAccessTokenWithDifferentUserId(t *testing.T) {
	accessTtl := 1
	tf := NewTokenFactory(accessTtl /*refreshTtl=*/, 20, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	refreshToken := tf.CreateRefreshToken(2)
	refreshTokenStr, err := tf.SignedString(refreshToken)
	assert.NoError(t, err)
	time.Sleep(time.Second) // wait for access token to expire
	newAccessTokenStr, err := tf.RefreshAccessToken(accessTokenStr, refreshTokenStr)
	assert.ErrorIs(t, err, ErrDistinctUserIds)
	assert.Empty(t, newAccessTokenStr)
}
func TestRefreshAccessTokenWithExpiredRefreshToken(t *testing.T) {
	tf := NewTokenFactory(10, 1, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	refreshToken := tf.CreateRefreshToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	refreshTokenStr, err := tf.SignedString(refreshToken)
	assert.NoError(t, err)
	time.Sleep(time.Second) // wait for refresh token to expire
	newAccessToken, err := tf.RefreshAccessToken(accessTokenStr, refreshTokenStr)
	assert.ErrorIs(t, err, ErrExpiredRefreshToken)
	assert.Empty(t, newAccessToken)
}

func TestRefreshAccessTokenWithInvalidAccessTokenSignature(t *testing.T) {
	tf := NewTokenFactory(10, 1, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	refreshToken := tf.CreateRefreshToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	refreshTokenStr, err := tf.SignedString(refreshToken)
	accessTokenStr += "1" // mess up signature at the end of the jwt
	assert.NoError(t, err)
	newAccessToken, err := tf.RefreshAccessToken(accessTokenStr, refreshTokenStr)
	assert.ErrorIs(t, err, ErrInvalidSignature)
	assert.Empty(t, newAccessToken)
}

func TestRefreshAccessTokenWithInvalidRefreshTokenSignature(t *testing.T) {
	tf := NewTokenFactory(10, 1, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	refreshToken := tf.CreateRefreshToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	refreshTokenStr, err := tf.SignedString(refreshToken)
	refreshTokenStr += "1" // mess up signature at the end of the jwt
	assert.NoError(t, err)
	newAccessToken, err := tf.RefreshAccessToken(accessTokenStr, refreshTokenStr)
	assert.ErrorIs(t, err, ErrInvalidSignature)
	assert.Empty(t, newAccessToken)
}

func TestRefreshAccessTokenWithDistinctUserIds(t *testing.T) {
	tf := NewTokenFactory(10, 1, generateRandomSecret())
	accessToken := tf.CreateAccessToken(1)
	refreshToken := tf.CreateRefreshToken(1)
	accessTokenStr, err := tf.SignedString(accessToken)
	assert.NoError(t, err)
	refreshTokenStr, err := tf.SignedString(refreshToken)
	accessTokenStr += "1" // mess up signature at the end of the jwt
	assert.NoError(t, err)
	newAccessToken, err := tf.RefreshAccessToken(accessTokenStr, refreshTokenStr)
	assert.Error(t, err)
	assert.Empty(t, newAccessToken)
}
