package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
	"net/http"
	"strings"
)

//type contextKey string

// const UserIDKey contextKey = "user_id"
const UserIDKey = "user_id" // string, not custom type

func JWTAuth(secret []byte, queries *postgresql.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// Skip token validation for auth endpoints
			if path == "/api/v1/auth/register" || path == "/api/v1/auth/login" {
				return next(c)
			}

			// Validate token for all other endpoints
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
			}

			tokenString := parts[1]
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Check blacklist
			jti, ok := claims["jti"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: missing jti")
			}
			blacklisted, err := queries.IsTokenBlacklisted(c.Request().Context(), jti)
			if err != nil {
				// Log error but allow? Deny access to be safe.
				return echo.NewHTTPError(http.StatusInternalServerError, "Authentication error")
			}
			if blacklisted {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has been revoked")
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}
			userID := int64(userIDFloat)

			// Store in Echo context (backward compatibility)
			c.Set("user_id", userID)
			c.Set("jti", jti)

			// Store in request context (for strict server handlers)
			ctx := context.WithValue(c.Request().Context(), "user_id", userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
