package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog-be/internal/api"
	"github.com/satriaardiperdana-2020/monelog-be/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
)

type UserHandler struct {
	Queries *postgresql.Queries
}

func (h *UserHandler) SoftDeleteUser(ctx context.Context, req api.SoftDeleteUserRequestObject) (api.SoftDeleteUserResponseObject, error) {
	// Get the authenticated user ID from the context (set by JWT middleware)
	currentUserID := ctx.Value(middleware.UserIDKey).(int64)

	// Ensure the user can only delete their own account
	if req.Id != currentUserID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "You can only delete your own account")
	}

	// Perform soft delete
	err := h.Queries.SoftDeleteUser(ctx, req.Id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return api.SoftDeleteUser204Response{}, nil
}
