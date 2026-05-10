package handlers

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog-be/internal/api"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type AuthHandler struct {
	Queries   *postgresql.Queries
	JWTSecret []byte
}

/*// Request & response structs (tidak tergantung OpenAPI)
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
*/
/*type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        int64     `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user"`
}*/

func (s *AuthHandler) Register(ctx context.Context, req api.RegisterRequestObject) (api.RegisterResponseObject, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Body.Password), bcrypt.DefaultCost)
	user, err := s.Queries.CreateUser(ctx, postgresql.CreateUserParams{
		Email:        string(req.Body.Email),
		PasswordHash: string(hashed),
		Name:         req.Body.Name,
		Picture:      pgtype.Text{String: "", Valid: true},
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusConflict, "Email already exists")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"jti":     uuid.New().String(), // unique ID
	})
	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Create api.User struct (fields may be pointers)
	created := user.CreatedAt
	apiUser := api.User{
		Id:        &user.ID,
		Email:     &user.Email, // pointer to string
		Name:      &user.Name,  // pointer to string
		CreatedAt: &created,    // pointer to time.Time
	}

	// Return response with pointer to apiUser
	return api.Register201JSONResponse(api.AuthResponse{
		Token: &tokenString,
		User:  &apiUser, // ✅ pointer to api.User
	}), nil
}

func (s *AuthHandler) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	user, err := s.Queries.GetUserByEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Body.Password)); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"jti":     uuid.New().String(), // unique ID
	})
	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	userID := user.ID
	userEmail := user.Email
	userName := user.Name
	userCreated := user.CreatedAt

	apiUser := api.User{
		Id:        &userID,
		Email:     &userEmail,
		Name:      &userName,
		CreatedAt: &userCreated,
	}

	return api.Login200JSONResponse(api.AuthResponse{
		Token: &tokenString,
		User:  &apiUser,
	}), nil
}

func (h *AuthHandler) Logout(c echo.Context) error {
	jti, ok := c.Get("jti").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No valid token context"})
	}
	// Optional: get expiration from the token (you can parse it again or store in context)
	// The token is still valid; we need its expiration time to set in blacklist.
	// Better: In middleware, also store expiration.
	// We'll read from claims again (simple).

	// Parse token string from Authorization header again (simpler)
	authHeader := c.Request().Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid authorization header"})
	}

	tokenString := parts[1]
	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid token"})
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		exp = 0
	}
	expiresAt := time.Unix(int64(exp), 0)

	err = h.Queries.AddTokenToBlacklist(c.Request().Context(), postgresql.AddTokenToBlacklistParams{
		Jti:       jti,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to revoke token"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
