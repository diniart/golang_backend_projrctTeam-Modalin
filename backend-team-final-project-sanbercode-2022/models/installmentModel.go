package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Installment struct {
		ID     uuid.UUID `json:"id" gorm:"primary_key"`
		Amount int       `json:"amount"`

		Status string `json:"status" gorm:"type:varchar(255)"`
		Type   string `json:"type" gorm:"type:varchar(255)"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		// Relation
		ProjectsID uuid.UUID `json:"projects_id" gorm:"index"`
		Projects   Projects  `json:"projects" gorm:"foreignKey:ProjectsID"`
	}

	InstallmentInput struct {
		Status string `json:"status"`
	}
)
