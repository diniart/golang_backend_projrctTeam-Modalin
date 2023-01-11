package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	CartItems struct {
		ID        uuid.UUID `json:"id" gorm:"primary_key"`
		Quantity  int       `json:"quantity"  gorm:"type:int;not null"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		// Relation
		ProjectsID uuid.UUID `json:"projects_id" gorm:"index"`
		Projects   Projects  `json:"projects" gorm:"foreignKey:ProjectsID"`

		ShoppingCartID uuid.UUID    `json:"shopping_cart_id" gorm:"index"`
		ShoppingCart   ShoppingCart `json:"shopping_cart" gorm:"foreignKey:ShoppingCartID"`
	}
)
