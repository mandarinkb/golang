package domain

import (
	"context"

	"github.com/mandarinkb/go-fiber-with-mongodb/model"
)

type AccountRepository interface {
	CreateOneAccount(ctx context.Context, account interface{}) error
	CreateManyAccount(ctx context.Context, accounts []interface{}) error
	FindOneByID(ctx context.Context, id string) (*model.Account, error)
	FindByWalletID(ctx context.Context, walletID string) ([]model.Account, error)
	FindOneByWalletID(ctx context.Context, walletID string) (*model.Account, error)
	UpdateOneByID(ctx context.Context, id string, updator interface{}) error
	UpdateOneByWalletID(ctx context.Context, walletID string, updator interface{}) error
	UpdateManyByWalletID(ctx context.Context, walletID string, updator interface{}) error
	ReplaceOneByWalletID(ctx context.Context, walletID string, replacement interface{}) error
	FindOneAndUpdateByWalletID(ctx context.Context, walletID string, updator interface{}) (*model.Account, error)
	FindOneAndReplace(ctx context.Context, email string, updator interface{}) (*model.Account, error)
	DeleteOneByWalletID(ctx context.Context, walletID string) error
	DeleteManyByWalletID(ctx context.Context, walletID string) error
	EarnPromotion(ctx context.Context, walletID string, key string, count int) (*model.Account, error)
}
