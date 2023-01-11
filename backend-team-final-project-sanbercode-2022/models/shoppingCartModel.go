package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	ShoppingCart struct {
		ID            uuid.UUID `json:"id" gorm:"primary_key"`
		PaymentStatus string    `json:"payment_status"  gorm:"type:varchar(255);not null"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`

		// Relation
		CartItems []CartItems

		UserID uuid.UUID `json:"user_id" gorm:"index"`
		User   User      `json:"user" gorm:"foreignKey:UserID"`
	}
)
