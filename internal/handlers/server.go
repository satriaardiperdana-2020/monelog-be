package handlers

import (
	"context"

	"github.com/satriaardiperdana-2020/monelog-be/internal/api"
	"github.com/satriaardiperdana-2020/monelog-be/internal/repository/postgresql"
)

// Server implements the generated StrictServerInterface by delegating to specific handlers.
type Server struct {
	Queries     *postgresql.Queries
	JWTSecret   []byte
	Auth        *AuthHandler
	Category    *CategoryHandler
	Transaction *TransactionHandler
	User        *UserHandler
}

// ---------- Auth ----------
func (s *Server) Register(ctx context.Context, req api.RegisterRequestObject) (api.RegisterResponseObject, error) {
	return s.Auth.Register(ctx, req)
}

func (s *Server) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	return s.Auth.Login(ctx, req)
}

// ---------- Categories ----------
func (s *Server) CreateCategory(ctx context.Context, req api.CreateCategoryRequestObject) (api.CreateCategoryResponseObject, error) {
	return s.Category.CreateCategory(ctx, req)
}

func (s *Server) GetCategories(ctx context.Context, req api.GetCategoriesRequestObject) (api.GetCategoriesResponseObject, error) {
	return s.Category.GetCategories(ctx, req)
}

func (s *Server) UpdateCategory(ctx context.Context, req api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	return s.Category.UpdateCategory(ctx, req)
}

func (s *Server) DeleteCategory(ctx context.Context, req api.DeleteCategoryRequestObject) (api.DeleteCategoryResponseObject, error) {
	return s.Category.DeleteCategory(ctx, req)
}

// ---------- Transactions ----------
func (s *Server) CreateTransaction(ctx context.Context, req api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	return s.Transaction.CreateTransaction(ctx, req)
}

/*func (s *Server) GetTransactions(ctx context.Context, req api.GetTransactionsBetweenDatesWithNoteRequestObject) (api.GetTransactionsBetweenDatesWithNoteResponseObject, error) {
	return s.Transaction.GetTransactions(ctx, req)
}*/

func (s *Server) UpdateTransaction(ctx context.Context, req api.UpdateTransactionRequestObject) (api.UpdateTransactionResponseObject, error) {
	return s.Transaction.UpdateTransaction(ctx, req)
}

func (s *Server) SoftDeleteTransaction(ctx context.Context, req api.SoftDeleteTransactionRequestObject) (api.SoftDeleteTransactionResponseObject, error) {
	return s.Transaction.SoftDeleteTransaction(ctx, req)
}

// ---------- Reports (delegated to TransactionHandler) ----------
func (s *Server) GetMainPageSummary(ctx context.Context, req api.GetMainPageSummaryRequestObject) (api.GetMainPageSummaryResponseObject, error) {
	return s.Transaction.GetMainPageSummary(ctx, req)
}

func (s *Server) GetTransactionDetailsByDate(ctx context.Context, req api.GetTransactionDetailsByDateRequestObject) (api.GetTransactionDetailsByDateResponseObject, error) {
	return s.Transaction.GetTransactionDetailsByDate(ctx, req)
}

func (s *Server) GetLast7DaysDetail(ctx context.Context, req api.GetLast7DaysDetailRequestObject) (api.GetLast7DaysDetailResponseObject, error) {
	return s.Transaction.GetLast7DaysDetail(ctx, req)
}

func (s *Server) GetLast30DaysDetail(ctx context.Context, req api.GetLast30DaysDetailRequestObject) (api.GetLast30DaysDetailResponseObject, error) {
	return s.Transaction.GetLast30DaysDetail(ctx, req)
}

func (s *Server) GetTransactionsBetweenDatesWithNote(ctx context.Context, req api.GetTransactionsBetweenDatesWithNoteRequestObject) (api.GetTransactionsBetweenDatesWithNoteResponseObject, error) {
	return s.Transaction.GetTransactionsBetweenDatesWithNote(ctx, req)
}

// ---------- User (soft delete) ----------
func (s *Server) SoftDeleteUser(ctx context.Context, req api.SoftDeleteUserRequestObject) (api.SoftDeleteUserResponseObject, error) {
	return s.User.SoftDeleteUser(ctx, req)
}
