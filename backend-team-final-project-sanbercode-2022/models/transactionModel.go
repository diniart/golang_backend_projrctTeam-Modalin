package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Transaction struct {
		ID             uuid.UUID `json:"id" gorm:"primary_key"`
		Debit          int       `json:"debit" gorm:"type:int"`
		Credit         int       `json:"credit" gorm:"type:int"`
		Sender         uuid.UUID `json:"sender"`
		ShoppingCartID uuid.UUID `json:"shopping_cart_id"`
		ProjectID      uuid.UUID `json:"project_id"`
		InstallmentID  uuid.UUID `json:"installment_id"`

		Status string `json:"status" gorm:"type:varchar(255)"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		// Relation
		UserID uuid.UUID `json:"user_id" gorm:"index"`
		User   User      `json:"user" gorm:"foreignKey:UserID"`
	}

	TransactionInput struct {
		Debit          int       `json:"debit"`
		Credit         int       `json:"credit"`
		Sender         uuid.UUID `json:"sender"`
		Status         string    `json:"status"`
		ShoppingCartID uuid.UUID `json:"shopping_cart_id"`
		ProjectID      uuid.UUID `json:"project_id"`
	}
)
