package domain

import "github.com/mandarinkb/go-unit-test-with-mockgen/model"

type PromotionSettingRepository interface {
	Get(id int) (*model.PromotionSetting, error)
}
