package models

import (
	"gorm.io/gorm"
)

type Academic_Period struct {
	gorm.Model
	NAME   string `json:"name"`
	STATUS string `json:"status"`
}
