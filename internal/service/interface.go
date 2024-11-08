package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/slipneff/auto-yt/internal/models"
	"github.com/slipneff/auto-yt/pkg/clients/youtube"
)

type Service struct {
	storage  Storage
	ytClient *youtube.Client
}
type Storage interface {
	CreateAccount(ctx context.Context, input *models.Account) (*models.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetAccounts(ctx context.Context) ([]*models.Account, error)
	BatchCreateAccounts(ctx context.Context, accounts []*models.Account) error
	GetAccountsWithToken(ctx context.Context) ([]*models.Account, error)
	GetAccountsWithoutSecret(ctx context.Context) ([]*models.Account, error)
	ConfirmAuth(ctx context.Context, id uuid.UUID) error
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
}

func NewService(storage Storage, ytClient *youtube.Client) *Service {
	return &Service{
		storage:  storage,
		ytClient: ytClient,
	}
}
