package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthHandler struct {
	Queries   *db.Queries
	JWTSecret []byte
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req api.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user, err := h.Queries.CreateUser(c.Request().Context(), db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashed),
		Name:         req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(h.JWTSecret)
	return c.JSON(http.StatusCreated, api.AuthResponse{Token: tokenString, User: user})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req api.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	user, err := h.Queries.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(h.JWTSecret)
	return c.JSON(http.StatusOK, api.AuthResponse{Token: tokenString, User: user})
}
