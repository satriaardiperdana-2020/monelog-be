package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	runtime_types "github.com/oapi-codegen/runtime/types"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
	"github.com/satriaardiperdana-2020/monelog/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"net/http"
)

type TransactionHandler struct {
	Queries *postgresql.Queries
}

// Helper: konversi postgresql.Transaction -> api.Transaction
func toApiTransaction(tx postgresql.Transaction) api.Transaction {
	// Konversi Date (pgtype.Date) -> runtime_types.Date
	var dateVal runtime_types.Date
	if tx.Date.Valid {
		dateVal = runtime_types.Date{Time: tx.Date.Time}
	}

	// Konversi Note (pgtype.Text) -> *string
	var notePtr *string
	if tx.Note.Valid && tx.Note.String != "" {
		notePtr = &tx.Note.String
	}

	// Konversi CreatedAt (pgtype.Timestamptz) -> *time.Time
	/*var createdAtPtr *time.Time
	if tx.CreatedAt.Valid {
		createdAtPtr = &tx.CreatedAt.Time
	}
	*/
	return api.Transaction{
		Id:         tx.ID,
		UserId:     tx.UserID,
		CategoryId: tx.CategoryID,
		Amount:     tx.Amount,
		Note:       notePtr,
		Date:       dateVal, // ✅ tipe openapi_types.Date
		CreatedAt:  &tx.CreatedAt,
	}
}

func (h *TransactionHandler) CreateTransaction(ctx context.Context, request api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	//userID := ctx.Get("user_id").(int64)
	userIDVal := ctx.Value(middleware.UserIDKey)
	if userIDVal == nil {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Missing user_id"}
	}
	// request.Body.Date bertipe string (dari OpenAPI)

	userID := userIDVal.(int64)
	date := request.Body.Date.Time

	/*if err != nil {
		return nil, ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, use YYYY-MM-DD"})
	}*/

	var noteText pgtype.Text
	if request.Body.Note != nil && *request.Body.Note != "" {
		noteText = pgtype.Text{String: *request.Body.Note, Valid: true}
	} else {
		noteText = pgtype.Text{Valid: false}
	}

	tx, err := h.Queries.CreateTransaction(ctx, postgresql.CreateTransactionParams{
		UserID:     userID,
		CategoryID: request.Body.CategoryId,
		Amount:     request.Body.Amount,
		Note:       noteText,
		Date:       pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	apiTx := toApiTransaction(tx)
	return api.CreateTransaction201JSONResponse(apiTx), nil
}

func (h *TransactionHandler) GetTransactions(ctx context.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	userIDVal := ctx.Value(middleware.UserIDKey)

	if userIDVal == nil {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Missing user_id"}
	}
	userID := userIDVal.(int64)

	// FIX: use .Time field, not parsing from string
	fromTime := request.Params.From.Time
	toTime := request.Params.To.Time

	transactions, err := h.Queries.GetTransactionsByUserAndDateRange(ctx, postgresql.GetTransactionsByUserAndDateRangeParams{
		UserID: userID,
		Date:   pgtype.Date{Time: fromTime, Valid: true},
		Date_2: pgtype.Date{Time: toTime, Valid: true},
	})
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	// Konversi slice ke api.Transaction
	apiTransactions := make([]api.Transaction, len(transactions))
	for i, tx := range transactions {
		apiTransactions[i] = toApiTransaction(tx)
	}
	return api.GetTransactions200JSONResponse(apiTransactions), nil
}
