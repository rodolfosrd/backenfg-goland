package models

type Docente struct {
	IDPERSONA int    `json:"idpersona"`
	CODIGO    string `json:"codigo"`
	PERSONA   Person
}
