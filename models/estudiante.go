package models

type Estudiante struct {
	IDPERSONA int    `json:"idpersona"`
	CODIGO    string `json:"codigo"`
	PERSONA   Person
	// use CompanyRefer as foreign key
}
