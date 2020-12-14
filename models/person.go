package models

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	NOMBRE    string `json:"nombre"`
	APPATERNO string `json:"appaterno"`
	APMATERNO string `json:"apmaterno"`
	DNI       string `json:"dni"`
	DIRECCION string `json:"direccion"`
	TELEFONO  string `json:"telefono"`
	FECHA     string `json:"fecha"`
	USUARIO   string `json:"usuario"`
	PASSWORD  string `json:"password"`
}
