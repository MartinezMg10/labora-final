package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID        string `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	Name      string `json:"name" gorm:"type:char(50);not null;`
	StartDate string `json:start_date"`
	EndDate   string `json:"end_date"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	Deleted   string `json:"-"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) (err error) {

	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
