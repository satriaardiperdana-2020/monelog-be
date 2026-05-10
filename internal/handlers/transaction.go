package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	runtime_types "github.com/oapi-codegen/runtime/types"
	"github.com/satriaardiperdana-2020/monelog-be/internal/api"
	"github.com/satriaardiperdana-2020/monelog-be/internal/middleware"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
	"net/http"
)

type TransactionHandler struct {
	Queries *postgresql.Queries
}

// CreateTransaction inserts a new transaction.
func (s *TransactionHandler) CreateTransaction(ctx context.Context, req api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	//date, _ := time.Parse("2006-01-02", req.Body.Date)
	date := req.Body.Date.Time
	var noteText pgtype.Text
	if req.Body.Note != nil && *req.Body.Note != "" {
		noteText = pgtype.Text{String: *req.Body.Note, Valid: true}
	}
	tx, err := s.Queries.CreateTransaction(ctx, postgresql.CreateTransactionParams{
		UserID:     userID,
		CategoryID: req.Body.CategoryId,
		Amount:     req.Body.Amount,
		Note:       noteText,
		Date:       pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.CreateTransaction201JSONResponse(toApiTransaction(tx)), nil
}

// view range between
// For GET /transactions?from=...&to=...
/*func (s *TransactionHandler) GetTransactions(ctx context.Context, req api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	// FIX: use .Time field, not parsing from string
	fromTime := req.Params.From.Time
	toTime := req.Params.To.Time
	txs, err := s.Queries.GetMainPageSummaryTransactions(ctx, postgresql.GetTransactionsBetweenDatesWithNoteParams{
		UserID: userID,
		Date:   pgtype.Date{Time: fromTime, Valid: true},
		Date_2: pgtype.Date{Time: toTime, Valid: true},
	})
	if err != nil {
		return nil, echo.NewHTTPError(500, err.Error())
	}
	result := make([]api.TransactionDetail, len(txs))
	for i, tx := range txs {
		dateStr := runtime_types.Date{Time: tx.Date.Time}
		noteStr := tx.Note.String
		result[i] = api.TransactionDetail{
			Date:         &dateStr,
			Note:         &noteStr,
			TotalIncome:  &tx.TotalIncome,
			TotalExpense: &tx.TotalExpense,
			Balance:      &tx.Balance,
		}
	}
	return api.GetTransactions200JSONResponse(result), nil
}*/

// UpdateTransaction modifies an existing transaction.
func (h *TransactionHandler) UpdateTransaction(ctx context.Context, request api.UpdateTransactionRequestObject) (api.UpdateTransactionResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	date := request.Body.Date.Time
	var noteText pgtype.Text
	if request.Body.Note != nil && *request.Body.Note != "" {
		noteText = pgtype.Text{String: *request.Body.Note, Valid: true}
	}

	tx, err := h.Queries.UpdateTransaction(ctx, postgresql.UpdateTransactionParams{
		ID:         request.Id,
		UserID:     userID,
		CategoryID: request.Body.CategoryId,
		Amount:     request.Body.Amount,
		Note:       noteText,
		Date:       pgtype.Date{Time: date, Valid: true},
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.UpdateTransaction200JSONResponse(toApiTransaction(tx)), nil
}

func (h *TransactionHandler) SoftDeleteTransaction(ctx context.Context, request api.SoftDeleteTransactionRequestObject) (api.SoftDeleteTransactionResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)

	// Call sqlc method (returns error only)
	err := h.Queries.SoftDeleteTransaction(ctx, postgresql.SoftDeleteTransactionParams{
		ID:     request.Id,
		UserID: userID,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return api.SoftDeleteTransaction204Response{}, nil
}

// ----------------------------------------------------------------------------
// Report endpoints (used directly by Echo routes)
// ----------------------------------------------------------------------------
// view dihalaman utama
// GetMainPageSummary returns daily income/expense/balance, last 10 days with data.
/*func (h *TransactionHandler) GetMainPageSummary(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	rows, err := h.Queries.GetMainPageSummaryTransactions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rows)
}
*/

func (s *TransactionHandler) GetMainPageSummary(ctx context.Context, req api.GetMainPageSummaryRequestObject) (api.GetMainPageSummaryResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	rows, err := s.Queries.GetMainPageSummaryTransactions(ctx, userID)

	if err != nil {
		return nil, echo.NewHTTPError(500, err.Error())
	}
	res := make([]api.DailySummary, len(rows))
	for i, r := range rows {
		dateVal := runtime_types.Date{Time: r.Date.Time}
		res[i] = api.DailySummary{
			Date:         &dateVal,
			TotalIncome:  &r.TotalIncome,
			TotalExpense: &r.TotalExpense,
			Balance:      &r.Balance,
		}
	}
	return api.GetMainPageSummary200JSONResponse(res), nil
}

func (s *TransactionHandler) GetTransactionDetailsByDate(ctx context.Context, req api.GetTransactionDetailsByDateRequestObject) (api.GetTransactionDetailsByDateResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	rows, err := s.Queries.GetTransactionDetailsByDateLimit10(ctx, userID)

	if err != nil {
		return nil, echo.NewHTTPError(500, err.Error())
	}
	res := make([]api.TransactionDetail, len(rows))
	for i, r := range rows {
		dateVal := runtime_types.Date{Time: r.Date.Time}
		noteStr := r.Note.String
		res[i] = api.TransactionDetail{
			Date:         &dateVal,
			CategoryName: &r.CategoryName,
			CategoryType: &r.CategoryType,
			Note:         &noteStr,
			TotalIncome:  &r.TotalIncome,
			TotalExpense: &r.TotalExpense,
			Balance:      &r.Balance,
		}
	}
	return api.GetTransactionDetailsByDate200JSONResponse(res), nil
}

func (s *TransactionHandler) GetLast7DaysDetail(ctx context.Context, req api.GetLast7DaysDetailRequestObject) (api.GetLast7DaysDetailResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	rows, err := s.Queries.GetLast7DaysDetail(ctx, userID)
	if err != nil {
		return nil, echo.NewHTTPError(500, err.Error())
	}
	res := make([]api.TransactionDetail, len(rows))
	for i, r := range rows {
		dateVal := runtime_types.Date{Time: r.Date.Time}
		noteStr := r.Note.String
		res[i] = api.TransactionDetail{
			Date:         &dateVal,
			Note:         &noteStr,
			TotalIncome:  &r.TotalIncome,
			TotalExpense: &r.TotalExpense,
			Balance:      &r.Balance,
		}
	}
	return api.GetLast7DaysDetail200JSONResponse(res), nil
}

// GetLast30DaysDetail returns all transactions (grouped by date and note) for the last 30 days.
func (s *TransactionHandler) GetLast30DaysDetail(ctx context.Context, req api.GetLast30DaysDetailRequestObject) (api.GetLast30DaysDetailResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	rows, err := s.Queries.GetLast30DaysDetail(ctx, userID)
	if err != nil {
		return nil, echo.NewHTTPError(500, err.Error())
	}
	res := make([]api.TransactionDetail, len(rows))
	for i, r := range rows {
		dateVal := runtime_types.Date{Time: r.Date.Time}
		noteStr := r.Note.String
		res[i] = api.TransactionDetail{
			Date:         &dateVal,
			Note:         &noteStr,
			TotalIncome:  &r.TotalIncome,
			TotalExpense: &r.TotalExpense,
			Balance:      &r.Balance,
		}
	}
	return api.GetLast30DaysDetail200JSONResponse(res), nil
}

// GetTransactionsBetweenDatesWithNote returns all transactions (grouped by date and note) within a given range.
/*func (h *TransactionHandler) GetTransactionsBetweenDatesWithNote(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid 'from' date, use YYYY-MM-DD"})
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid 'to' date, use YYYY-MM-DD"})
	}

	rows, err := h.Queries.GetTransactionsBetweenDatesWithNote(c.Request().Context(), postgresql.GetTransactionsBetweenDatesWithNoteParams{
		UserID: userID,
		From:   from,
		To:     to,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rows)
}*/
//SAma range between
func (s *TransactionHandler) GetTransactionsBetweenDatesWithNote(ctx context.Context, req api.GetTransactionsBetweenDatesWithNoteRequestObject) (api.GetTransactionsBetweenDatesWithNoteResponseObject, error) {
	userID := ctx.Value(middleware.UserIDKey).(int64)
	fromTime := req.Params.From.Time
	toTime := req.Params.To.Time

	rows, err := s.Queries.GetTransactionsBetweenDatesWithNote(ctx, postgresql.GetTransactionsBetweenDatesWithNoteParams{
		UserID: userID,
		Date:   pgtype.Date{Time: fromTime, Valid: true},
		Date_2: pgtype.Date{Time: toTime, Valid: true},
	})
	if err != nil {
		return nil, echo.NewHTTPError(500, err.Error())
	}
	res := make([]api.TransactionDetail, len(rows))
	for i, r := range rows {
		dateStr := runtime_types.Date{Time: r.Date.Time}
		noteStr := r.Note.String
		res[i] = api.TransactionDetail{
			Date:         &dateStr,
			Note:         &noteStr,
			TotalIncome:  &r.TotalIncome,
			TotalExpense: &r.TotalExpense,
			Balance:      &r.Balance,
		}
	}
	return api.GetTransactionsBetweenDatesWithNote200JSONResponse(res), nil
}

// ----------------------------------------------------------------------------
// Helper: convert database transaction to API transaction
// ----------------------------------------------------------------------------
func toApiTransaction(tx postgresql.Transaction) api.Transaction {
	var dateVal runtime_types.Date
	if tx.Date.Valid {
		dateVal = runtime_types.Date{Time: tx.Date.Time}
	}
	var notePtr *string
	if tx.Note.Valid && tx.Note.String != "" {
		notePtr = &tx.Note.String
	}
	return api.Transaction{
		Id:         &tx.ID,
		UserId:     &tx.UserID,
		CategoryId: &tx.CategoryID,
		Amount:     &tx.Amount,
		Note:       notePtr,
		Date:       &dateVal,
		CreatedAt:  &tx.CreatedAt,
	}
}

// toApiTransactionFromSummary converts a daily summary row to api.Transaction
/*func toApiTransactionFromSummary(tx postgresql.GetMainPageSummaryTransactionsRow, userID int64) api.Transaction {
	var dateVal runtime_types.Date
	if tx.Date.Valid {
		dateVal = runtime_types.Date{Time: tx.Date.Time}
	}
	// Placeholder values because summary row has no transaction ID, category, note, etc.
	/*id := int64(0)
	cid := int64(0)
	note := ""
	// Use TotalIncome as the amount (or TotalExpense depending on context, but we'll use income)
	amount := row.TotalIncome

	return api.Transaction{
		Id:         &tx.ID,
		UserId:     &userID,
		CategoryId: tx.CategoryID,
		Amount:     &tx.Amount,
		Note:       notePtr,
		Date:       &dateVal,
		CreatedAt:  &tx.CreatedAt,
	}
}
*/
