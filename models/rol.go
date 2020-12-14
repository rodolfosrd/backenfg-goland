package models

import (
	"gorm.io/gorm"
)

type Rol struct {
	gorm.Model
	NOMBRE string `json:"nombre"`
	ESTADO int    `json:"estado"`
}
