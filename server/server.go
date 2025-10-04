package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/raian621/dump/auth"
	"github.com/raian621/dump/database"
	"github.com/raian621/dump/models/client"
	"github.com/raian621/dump/util"
)

type Server struct {
	e  *echo.Echo
	db *pgxpool.Pool
	tf *auth.TokenFactory
}

func (s *Server) Start(address string) error {
	return s.e.Start(address)
}

func (s *Server) Hello(c echo.Context) error {
	return c.HTML(200, "<h1>Hello</h1>")
}

func (s *Server) CreateUserWithCredentials(c echo.Context) error {
	creds := &client.Credentials{}
	if err := json.NewDecoder(c.Request().Body).Decode(creds); err != nil {
		c.Logger().Warn("Failed to decode user credentials: ", err)
		return c.String(500, "Failed to decode user credentials")
	}
	if exists, err := database.UsernameExists(s.db, creds.Username); err != nil {
		c.Logger().Error("Unexpected error while checking for username: ", err)
		return c.String(500, "Unexpected error")
	} else if exists {
		return c.String(http.StatusConflict, "Username already exists")
	}
	if err := database.CreateUserWithCredentials(s.db, creds.ToStorageModel()); err != nil {
		c.Logger().Error("Failed to create user: ", err)
		return c.String(500, "Failed to create user")
	}
	return nil
}

// Sign a user in with credentials (username and password)
func (s *Server) SignInWithCredentials(c echo.Context) error {
	creds := &client.Credentials{}
	if err := json.NewDecoder(c.Request().Body).Decode(creds); err != nil {
		c.Logger().Warn("Failed to decode user credentials: ", err)
		return c.String(
			http.StatusUnprocessableEntity, "Failed to decode user credentials")
	}

	passhash, err := database.GetPasshashForUsername(s.db, creds.Username)
	if err == sql.ErrNoRows {
		c.Logger().Error("Couldn't find username: ", creds.Username)
		return c.String(http.StatusBadRequest, "Incorrect username or password")
	} else if err != nil {
		c.Logger().Error("Unexpected error occurred: ", err)
		return c.String(http.StatusInternalServerError, "Unexpected error occurred")
	}

	if !util.ValidatePassword(creds.Password, passhash) {
		c.Logger().Error("Incorrect password for user: ", creds.Username)
		return c.String(http.StatusNotFound, "Incorrect username or password")
	}

	userId, err := database.GetUserIdFromUsername(s.db, creds.Username)
	accessToken := s.tf.CreateAccessToken(userId)
	refreshToken := s.tf.CreateRefreshToken(userId)

	accessTokenStr, err := s.tf.SignedString(accessToken)
	if err != nil {
		c.Logger().Error("Unexpected error occurred: ", err)
		return c.String(500, "Unexpected error occurred")
	}

	refreshTokenStr, err := s.tf.SignedString(refreshToken)
	if err != nil {
		c.Logger().Error("Unexpected error occurred: ", err)
		return c.String(500, "Unexpected error occurred")
	}

	return c.JSON(200, client.AuthPayload{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	})
}

func (s *Server) RefreshAccessToken(c echo.Context) error {
	tokens := &client.AuthPayload{}
	if err := json.NewDecoder(c.Request().Body).Decode(tokens); err != nil {
		c.Logger().Errorf("Failed to parse refresh access token payload: err=%v, payload=%v\n")
		return c.String(500, "Unexpected error occurred")
	}

	newAccessTokenStr, err := s.tf.RefreshAccessToken(tokens.AccessToken, tokens.RefreshToken)
	if err != nil {
		c.Logger().Error("Failed to refresh access token: ", err)
		return c.String(http.StatusUnauthorized, "Failed to refresh token")
	}

	return c.JSON(200, client.AuthPayload{
		AccessToken: newAccessTokenStr,
	})
}

func (s *Server) CreateVault(c echo.Context) error {
	return nil
}

func New() *Server {
	s := &Server{e: echo.New()}
	s.e.Use(middleware.Logger())
	return s
}

func (s *Server) AddDatabaseClient(db *pgxpool.Pool) {
	s.db = db
}

func (s *Server) AddTokenFactory(tf *auth.TokenFactory) {
	s.tf = tf
}

func (s *Server) AddHandlers() {
	s.e.GET("/hello", s.Hello)
	s.e.POST("/users/create", s.CreateUserWithCredentials)
	s.e.POST("/users/signin/credentials", s.SignInWithCredentials)
	s.e.POST("/users/signin/refresh", s.RefreshAccessToken)
	s.e.POST("/vaults/create", s.CreateVault, auth.AuthMiddleware(s.tf))
}
