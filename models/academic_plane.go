package models

import (
	"gorm.io/gorm"
)

type Academic_Plane struct {
	gorm.Model
	NAME     string `json:"name"`
	FECHA    string `json:"fecha"`
	CURSO_ID int    `json:"curso_id"`
	LEVEL_ID int    `json:"level_id"`
	STATUS   string `json:"status"`
}
