package usecase

import (
	"context"
	"fmt"

	"github.com/mandarinkb/go-unit-test-with-mockgen/domain"
)

var (
	ErrorAccount          = fmt.Errorf("error account")
	ErrorPromotionSetting = fmt.Errorf("error promotion setting")
)

type promotionUsecase struct {
	accountRepo domain.AccountRepository
	promoRepo   domain.PromotionSettingRepository
}

func NewPromotionUsecase(accountRepo domain.AccountRepository, promoRepo domain.PromotionSettingRepository) domain.PromotionUsecase {
	return &promotionUsecase{
		accountRepo: accountRepo,
		promoRepo:   promoRepo,
	}
}
func (u *promotionUsecase) Excute(ctx context.Context, promoID int) error {

	_, err := u.promoRepo.Get(promoID)
	if err != nil {
		return ErrorPromotionSetting
	}
	_, err = u.accountRepo.GetInfo(1)
	if err != nil {
		return ErrorAccount
	}

	return nil
}
