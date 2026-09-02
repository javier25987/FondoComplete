package anotaciones

type Anotacion struct {
	Codigo int    `json:"codigo"`
	IdUser int    `json:"id_user"`
	Nota   string `json:"nota"`
	Monto  int    `json:"monto"`
	Tipo   string `json:"tipo"`
	Fecha  string `json:"fecha"`
}

type PostAnotacion struct {
	Nota  string `json:"nota"`
	Monto int    `json:"monto"`
	Tipo  string `json:"tipo"`
}
