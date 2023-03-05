package domain

import (
	"context"
)

type PromotionUsecase interface {
	Excute(ctx context.Context, promoID int) error
}
