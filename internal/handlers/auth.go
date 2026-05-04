package handlers

import (
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Queries   *postgresql.Queries
	JWTSecret []byte
}

// Request & response structs (tidak tergantung OpenAPI)
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        int64     `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	user, err := h.Queries.CreateUser(c.Request().Context(), postgresql.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashed),
		Name:         req.Name,
		Picture:      pgtype.Text{String: "", Valid: true},
	})
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	response := AuthResponse{
		Token: tokenString,
		User: struct {
			ID        int64     `json:"id"`
			Email     string    `json:"email"`
			Name      string    `json:"name"` //
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
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
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	response := AuthResponse{
		Token: tokenString,
		User: struct {
			ID        int64     `json:"id"`
			Email     string    `json:"email"`
			Name      string    `json:"name"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	}
	return c.JSON(http.StatusOK, response)
}
