package models

import "github.com/google/uuid"

type (
	Category struct {
		ID       uuid.UUID `json:"id" gorm:"primary_key"`
		Category string    `json:"category"  gorm:"type:varchar(255);not null"`

		// Relation
		Projects []Projects
	}
)
