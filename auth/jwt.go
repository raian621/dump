package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredAccessToken   = errors.New("access token expired")
	ErrExpiredRefreshToken  = errors.New("refresh token expired")
	ErrInvalidSignature     = errors.New("token signature invalid")
	ErrDistinctUserIds      = errors.New("refresh and access token user IDs do not match")
	ErrDecodingAccessToken  = errors.New("error decoding access token")
	ErrDecodingRefreshToken = errors.New("error decoding refresh token")
)

type TokenFactory struct {
	accessTtl  int
	refreshTtl int
	secret     []byte
}

type AccessTokenClaims struct {
	UserId int32 `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserId int32 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenFactory(accessTtl int, refreshTtl int, secret []byte) *TokenFactory {
	return &TokenFactory{accessTtl, refreshTtl, secret}
}

func (f TokenFactory) CreateAccessToken(userId int32) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &AccessTokenClaims{
		UserId:           userId,
		RegisteredClaims: createRegisteredClaims(f.accessTtl),
	})
}

func (f TokenFactory) CreateRefreshToken(userId int32) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshTokenClaims{
		UserId:           userId,
		RegisteredClaims: createRegisteredClaims(f.refreshTtl),
	})
}

func createRegisteredClaims(ttl int) jwt.RegisteredClaims {
	issuedTime := time.Now()
	expireTime := issuedTime.Add(time.Second * time.Duration(ttl))
	return jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(issuedTime),
		ExpiresAt: jwt.NewNumericDate(expireTime),
	}
}

func (f TokenFactory) SignedString(token *jwt.Token) (string, error) {
	return token.SignedString(f.secret)
}

func (f TokenFactory) ParseAccessToken(tokenString string) (*jwt.Token, error) {
	return f.parseToken(tokenString, &AccessTokenClaims{})
}

func (f TokenFactory) ParseRefreshToken(tokenString string) (*jwt.Token, error) {
	return f.parseToken(tokenString, &RefreshTokenClaims{})
}

func (f TokenFactory) parseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return f.secret, nil
	})
}

func (f TokenFactory) RefreshAccessToken(accessTokenStr, refreshTokenStr string) (string, error) {
	accessToken, err := f.ParseAccessToken(accessTokenStr)
	// Ignore token expired errors
	if err != nil && !strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
		if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
			return "", ErrInvalidSignature
		}
		return "", err
	}

	refreshToken, err := f.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
			return "", ErrInvalidSignature
		} else if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			return "", ErrExpiredRefreshToken
		}
		return "", err
	}

	accessTokenClaims, ok := accessToken.Claims.(*AccessTokenClaims)
	if !ok {
		return "", ErrDecodingAccessToken
	}
	refreshTokenClaims, ok := refreshToken.Claims.(*RefreshTokenClaims)
	if !ok {
		return "", ErrDecodingRefreshToken
	}

	if accessTokenClaims.UserId != refreshTokenClaims.UserId {
		return "", ErrDistinctUserIds
	}

	issuedTime := time.Now()
	expireTime := issuedTime.Add(time.Second * time.Duration(f.accessTtl))
	accessTokenClaims.IssuedAt = jwt.NewNumericDate(issuedTime)
	accessTokenClaims.ExpiresAt = jwt.NewNumericDate(expireTime)

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	return f.SignedString(newAccessToken)
}
