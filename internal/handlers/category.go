package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/satriaardiperdana-2020/monelog-be/internal/api"
	"github.com/satriaardiperdana-2020/monelog-be/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
	"net/http"
)

type CategoryHandler struct {
	Queries *postgresql.Queries
}

// ---------- Categories ----------
func (s *CategoryHandler) CreateCategory(ctx context.Context, req api.CreateCategoryRequestObject) (api.CreateCategoryResponseObject, error) {
	userIDVal := ctx.Value("user_id")
	log.Info(" iduserIDVal cek :", userIDVal)
	//userIDVal, ok := ctx.Value("user_id").(int64)
	if userIDVal == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID type")
	}
	// Convert string to pgtype.Text if needed (your sqlc model uses string? Actually from earlier schema it's TEXT, but sqlc may generate string)
	log.Info("user id cek :", userID)
	log.Info("req cek :", req)
	cat, err := s.Queries.CreateCategory(ctx, postgresql.CreateCategoryParams{
		UserID: userID,
		Name:   req.Body.Name,
		Type:   string(req.Body.Type),
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.CreateCategory201JSONResponse(toApiCategory(cat)), nil
}

func (s *CategoryHandler) GetCategories(ctx context.Context, req api.GetCategoriesRequestObject) (api.GetCategoriesResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	cats, err := s.Queries.GetCategoriesByUser(ctx, userID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	result := make([]api.Category, len(cats))
	for i, c := range cats {
		result[i] = toApiCategory(c)
	}
	return api.GetCategories200JSONResponse(result), nil
}

func (s *CategoryHandler) UpdateCategory(ctx context.Context, req api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	// First, check if the category exists and belongs to the user
	_, err := s.Queries.GetCategoryByIdAndUser(ctx, postgresql.GetCategoryByIdAndUserParams{
		ID:     req.Id,
		UserID: userID,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}

	nameText := pgtype.Text{String: req.Body.Name, Valid: true}
	typeText := pgtype.Text{String: string(req.Body.Type), Valid: true}
	// Update the category
	cat, err := s.Queries.UpdateCategory(ctx, postgresql.UpdateCategoryParams{
		ID:     req.Id,
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
		Id:        &c.ID,
		UserId:    &c.UserID,
		Name:      &c.Name,
		Type:      (*api.CategoryType)(&c.Type),
		CreatedAt: &c.CreatedAt,
	}
}
