package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
	"github.com/satriaardiperdana-2020/monelog/internal/repository/postgresql"
	"net/http"
	"time"
)

type TransactionHandler struct {
	Queries *db.Queries
}

func (h *TransactionHandler) CreateTransaction(ctx echo.Context, request api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	userID := ctx.Get("user_id").(int64)
	date, _ := time.Parse("2006-01-02", request.Date)
	tx, err := h.Queries.CreateTransaction(ctx.Request().Context(), db.CreateTransactionParams{
		UserID:     userID,
		CategoryID: request.Body.CategoryId,
		Amount:     request.Body.Amount,
		Note:       request.Body.Note,
		Date:       date,
	})
	if err != nil {
		return nil, ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return api.CreateTransaction201JSONResponse(tx), nil
}

func (h *TransactionHandler) GetTransactions(ctx echo.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	userID := ctx.Get("user_id").(int64)
	from, _ := time.Parse("2006-01-02", request.Params.From)
	to, _ := time.Parse("2006-01-02", request.Params.To)
	transactions, err := h.Queries.GetTransactionsByUserAndDateRange(ctx.Request().Context(), db.GetTransactionsByUserAndDateRangeParams{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return nil, ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return api.GetTransactions200JSONResponse(transactions), nil
}
