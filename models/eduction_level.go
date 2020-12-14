package models

import (
	"gorm.io/gorm"
)

type Education_level struct {
	gorm.Model
	NAME   string `json:"name"`
	STATUS string `json:"status"`
}
