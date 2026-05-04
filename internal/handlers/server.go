package handlers

import (
	"context"
	"github.com/satriaardiperdana-2020/monelog/internal/api"
)

type Server struct {
	Transaction *TransactionHandler
	Category    *CategoryHandler
	User        *UserHandler
}

func (s *Server) CreateTransaction(ctx context.Context, request api.CreateTransactionRequestObject) (api.CreateTransactionResponseObject, error) {
	return s.Transaction.CreateTransaction(ctx, request)
}
func (s *Server) GetTransactions(ctx context.Context, request api.GetTransactionsRequestObject) (api.GetTransactionsResponseObject, error) {
	return s.Transaction.GetTransactions(ctx, request)
}
func (s *Server) UpdateTransaction(ctx context.Context, request api.UpdateTransactionRequestObject) (api.UpdateTransactionResponseObject, error) {
	return s.Transaction.UpdateTransaction(ctx, request)
}
func (s *Server) SoftDeleteTransaction(ctx context.Context, request api.SoftDeleteTransactionRequestObject) (api.SoftDeleteTransactionResponseObject, error) {
	return s.Transaction.SoftDeleteTransaction(ctx, request)
}
func (s *Server) CreateCategory(ctx context.Context, request api.CreateCategoryRequestObject) (api.CreateCategoryResponseObject, error) {
	return s.Category.CreateCategory(ctx, request)
}
func (s *Server) GetCategories(ctx context.Context, request api.GetCategoriesRequestObject) (api.GetCategoriesResponseObject, error) {
	return s.Category.GetCategories(ctx, request)
}
func (s *Server) UpdateCategory(ctx context.Context, request api.UpdateCategoryRequestObject) (api.UpdateCategoryResponseObject, error) {
	return s.Category.UpdateCategory(ctx, request)
}
func (s *Server) DeleteCategory(ctx context.Context, request api.DeleteCategoryRequestObject) (api.DeleteCategoryResponseObject, error) {
	return s.Category.DeleteCategory(ctx, request)
}
func (s *Server) SoftDeleteUser(ctx context.Context, request api.SoftDeleteUserRequestObject) (api.SoftDeleteUserResponseObject, error) {
	return s.User.SoftDeleteUser(ctx, request)
}
