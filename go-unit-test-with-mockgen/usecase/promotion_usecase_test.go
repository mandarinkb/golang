package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_domain "github.com/mandarinkb/go-unit-test-with-mockgen/domain/mock"
	"github.com/mandarinkb/go-unit-test-with-mockgen/model"
	"github.com/stretchr/testify/assert"
)

type promotionUsecaseDependency struct {
	t               *testing.T
	mockAccountRepo *mock_domain.MockAccountRepository
	mockPromoRepo   *mock_domain.MockPromotionSettingRepository
}

func TestExcute(t *testing.T) {
	testCases := []struct {
		name       string
		err        error
		buildStubs func(d *promotionUsecaseDependency)
	}{
		{
			name: "case_success",
			err:  nil,
			buildStubs: func(d *promotionUsecaseDependency) {
				d.CaseSuccess()
			},
		},
		{
			name: "case_promotion_setting_error",
			err:  ErrorPromotionSetting,
			buildStubs: func(d *promotionUsecaseDependency) {
				d.CasePromotionError()
			},
		},
		{
			name: "case_account_error",
			err:  ErrorAccount,
			buildStubs: func(d *promotionUsecaseDependency) {
				d.CaseAccountError()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			d := &promotionUsecaseDependency{
				t,
				mock_domain.NewMockAccountRepository(ctrl),
				mock_domain.NewMockPromotionSettingRepository(ctrl),
			}
			tc.buildStubs(d)
			promoUsecase := NewPromotionUsecase(d.mockAccountRepo, d.mockPromoRepo)
			err := promoUsecase.Excute(context.Background(), 1)
			assert.Equal(t, tc.err, err)
		})

	}
}

func (d *promotionUsecaseDependency) CaseSuccess() {
	d.mockPromoRepo.EXPECT().Get(gomock.Any()).Return(&model.PromotionSetting{}, nil)
	d.mockAccountRepo.EXPECT().GetInfo(gomock.Any()).Return(&model.Account{}, nil)
}

func (d *promotionUsecaseDependency) CasePromotionError() {
	d.mockPromoRepo.EXPECT().Get(gomock.Any()).Return(nil, errors.New("any error"))
}

func (d *promotionUsecaseDependency) CaseAccountError() {
	d.mockPromoRepo.EXPECT().Get(gomock.Any()).Return(&model.PromotionSetting{}, nil)
	d.mockAccountRepo.EXPECT().GetInfo(gomock.Any()).Return(nil, errors.New("any error"))
}
