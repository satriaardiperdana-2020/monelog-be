package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
	"github.com/satriaardiperdana-2020/monelog/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"net/http"
)

type UserHandler struct {
	Queries *postgresql.Queries
}

func (h *UserHandler) SoftDeleteUser(ctx context.Context, request api.SoftDeleteUserRequestObject) (api.SoftDeleteUserResponseObject, error) {
	currentUserID := ctx.Value(middleware.UserIDKey).(int64)
	if request.Id != currentUserID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "You can only delete your own account")
	}
	err := h.Queries.SoftDeleteUser(ctx, request.Id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.SoftDeleteUser204Response{}, nil
}
