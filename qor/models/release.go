package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Release struct {
	gorm.Model

	Comment string
	Date    *time.Time
	Log     string
}

func (r *Release) BeforeCreate() error {
	now := time.Now()
	r.Date = &now
	r.Log = "TODO..."
	return nil
}
