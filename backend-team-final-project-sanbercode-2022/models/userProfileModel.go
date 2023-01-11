package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key; unique"`
	Fullname          string    `json:"fullname" gorm:"type:varchar(255)" validate:"omitempty,min=1,max=255"`
	Phone             string    `json:"phone" gorm:"type:varchar(20)" validate:"omitempty,min=5,max=20"`
	KTP               int       `json:"ktp"`
	Address           string    `json:"address" gorm:"type:text"`
	City              string    `json:"city" gorm:"type:varchar(255)" validate:"omitempty,min=1,max=255"`
	Province          string    `json:"province" gorm:"type:varchar(255)" validate:"omitempty,min=1,max=255"`
	Gender            string    `json:"gender" gorm:"type:varchar(255)" validate:"omitempty,min=1,max=255"`
	ProfileURL        string    `json:"profile_url" gorm:"type:text"`
	BankAccountNumber int       `json:"bank_account_number"`
	BankName          string    `json:"bank_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
