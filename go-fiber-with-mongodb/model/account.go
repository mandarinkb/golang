package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID             primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Email          string               `json:"email" bson:"email"`
	CitizenID      string               `json:"citizen_id" bson:"citizen_id"`
	TitleTH        string               `json:"title_th" bson:"title_th"`
	FirstnameTH    string               `json:"first_name_th" bson:"first_name_th"`
	LastnameTH     string               `json:"last_name_th" bson:"last_name_th"`
	DateOfBirth    string               `json:"date_of_birth" bson:"date_of_birth"`
	Gender         string               `json:"gender" bson:"gender"`
	Nationality    string               `json:"nationality" bson:"nationality"`
	Education      int                  `json:"education" bson:"education"`
	EducationText  string               `json:"education_text" bson:"education_text"`
	ContactAddress Address              `json:"contact_address" bson:"contact_address"`
	CitizenImage   string               `json:"citizen_image" bson:"citizen_image"`
	FaceImage      string               `json:"face_image" bson:"face_image"`
	ImageDate      *time.Time           `json:"image_date" bson:"image_date"`
	MobileNo       string               `json:"mobile_no" bson:"mobile_no"`
	WalletID       string               `json:"wallet_id" bson:"wallet_id"`
	EarnPromotion  AccountEarnPromotion `json:"earn_promotion" bson:"earn_promotion,omitempty"`
}
type (
	Address struct {
		AddressNo   string `json:"address_no" bson:"address_no"`
		Subdistrict string `json:"subdistrict" bson:"subdistrict"`
		District    string `json:"district" bson:"district"`
		Province    string `json:"province" bson:"province"`
		Zipcode     string `json:"zipcode" bson:"zipcode"`
	}
	AccountEarnPromotion struct {
		Privilege map[string]interface{} `json:"privilege" bson:"privilege,omitempty"`
	}
)
