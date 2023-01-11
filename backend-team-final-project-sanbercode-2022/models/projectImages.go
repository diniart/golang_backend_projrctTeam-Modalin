package models

import "github.com/google/uuid"

type (
	Images struct {
		ID       uuid.UUID `json:"id" gorm:"primary_key"`
		ImageURL string    `json:"images_url"`

		// Relation
		ProjectsID uuid.UUID `json:"projects_id" gorm:"index"`
		Projects   Projects  `json:"projects" gorm:"foreignKey:ProjectsID"`
	}
)
