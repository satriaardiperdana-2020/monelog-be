package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
	"github.com/satriaardiperdana-2020/monelog/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"net/http"
)

type CategoryHandler struct {
	Queries *postgresql.Queries
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, request api.CreateCategoryRequestObject) (api.CreateCategoryResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	cat, err := h.Queries.CreateCategory(ctx, postgresql.CreateCategoryParams{
		UserID: userID,
		Name:   request.Body.Name,
		Type:   string(request.Body.Type),
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.CreateCategory201JSONResponse(toApiCategory(cat)), nil
}

func (h *CategoryHandler) GetCategories(ctx context.Context, request api.GetCategoriesRequestObject) (api.GetCategoriesResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	cats, err := h.Queries.GetCategoriesByUser(ctx, userID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	result := make([]api.Category, len(cats))
	for i, c := range cats {
		result[i] = toApiCategory(c)
	}
	return api.GetCategories200JSONResponse(result), nil
}

func (h *CategoryHandler) UpdateCategory(ctx context.Context, request api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	// Check ownership
	_, err := h.Queries.GetCategoryByIdAndUser(ctx, postgresql.GetCategoryByIdAndUserParams{
		ID:     request.Id,
		UserID: userID,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	nameText := pgtype.Text{String: request.Body.Name, Valid: true}
	typeText := pgtype.Text{String: string(request.Body.Type), Valid: true}

	cat, err := h.Queries.UpdateCategory(ctx, postgresql.UpdateCategoryParams{
		ID:     request.Id,
		UserID: userID,
		Name:   nameText,
		Type:   typeText,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.UpdateCategory200JSONResponse(toApiCategory(cat)), nil
}

func (h *CategoryHandler) DeleteCategory(ctx context.Context, request api.DeleteCategoryRequestObject) (api.DeleteCategoryResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	err := h.Queries.DeleteCategory(ctx, postgresql.DeleteCategoryParams{
		ID:     request.Id,
		UserID: userID,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.DeleteCategory204Response{}, nil
}

// Helper conversion
func toApiCategory(c postgresql.Category) api.Category {
	return api.Category{
		Id:        c.ID,
		UserId:    c.UserID,
		Name:      c.Name,
		Type:      api.CategoryType(c.Type),
		CreatedAt: &c.CreatedAt,
	}
}
