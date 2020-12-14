package models

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	NAME   string `json:"name"`
	STATUS string `json:"status"`
}
