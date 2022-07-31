package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	tn          = time.Now()
	id          = uuid.New()
	MockAccount = Account{
		Email:          "test@gmail.com",
		CitizenID:      "999999999",
		TitleTH:        "man",
		FirstnameTH:    "thongchai",
		LastnameTH:     "teangtham",
		DateOfBirth:    "12/12/1986",
		Gender:         "male",
		Nationality:    "thai",
		Education:      1,
		EducationText:  "Master's Degree",
		ContactAddress: Address{AddressNo: "99/2 ม2", Subdistrict: "ห้วยกระเจา", District: "ห้วยกระเจา", Province: "กาญจนบุรี", Zipcode: "71170"},

		CitizenImage:  "",
		FaceImage:     "",
		ImageDate:     &tn,
		MobileNo:      "0898989899",
		WalletID:      fmt.Sprint(id),
		EarnPromotion: AccountEarnPromotion{},
	}
)
