package handlers

import (
	"monelog/hello/internal/api"
	"time"

	"github.com/labstack/echo/v4"
)

// Pastikan struct ini mengimplementasikan api.StrictServerInterface
type TransactionHandler struct {
	// Tambahkan dependencies di sini, misalnya: queries *db.Queries
}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

// Implementasi fungsi untuk POST /transactions
func (h *TransactionHandler) CreateTransaction(ctx echo.Context, request api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	// --- Logika bisnis Anda di sini ---
	// 1. Ambil user_id dari context JWT
	// 2. Simpan ke DB menggunakan sqlc
	// 3. Return response yang sesuai

	// Contoh response sukses
	transaction := api.Transaction{
		Id:         1,
		UserId:     1, // Ambil dari JWT
		CategoryId: request.CategoryId,
		Amount:     request.Amount,
		Note:       request.Note,
		Date:       request.Date,
		CreatedAt:  time.Now(),
	}
	return api.CreateTransaction201JSONResponse(transaction), nil
}

// Implementasi fungsi untuk GET /transactions
func (h *TransactionHandler) GetTransactions(ctx echo.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	// --- Logika bisnis Anda di sini ---
	// 1. Ambil user_id dari context JWT
	// 2. Ambil parameter request.Params.From dan request.Params.To
	// 3. Query DB dengan sqlc
	// 4. Return response yang sesuai

	// Contoh response sukses
	transactions := []api.Transaction{}
	return api.GetTransactions200JSONResponse(transactions), nil
}
