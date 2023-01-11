package models

import (
	"time"

	"github.com/google/uuid"
)

type Projects struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" gorm:"type:varchar(255)" validate:"omitempty,min=1,max=255"`
	Margin      int       `json:"margin"`
	Duration    int       `json:"duration"`
	Periode     int       `json:"periode"`
	Description string    `json:"description" gorm:"type:text"`
	Quantity    int       `json:"quantity"`
	Price       int       `json:"price"`
	DueDate     time.Time `json:"dueDate"`
	Status      string    `json:"status" gorm:"default:pending"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relation
	CategoryID uuid.UUID `json:"category_id" gorm:"index"`
	Category   Category  `json:"category" gorm:"foreignKey:CategoryID"`

	UserID uuid.UUID `json:"user_id" gorm:"index"`
	User   User      `json:"user" gorm:"foreignKey:UserID"`

	Images      []Images      `json:"images"`
	CartItems   []CartItems   `json:"cart_items"`
	Installment []Installment `json:"installment"`
}
