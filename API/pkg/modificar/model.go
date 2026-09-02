package modificar

type FormatoCrearUsuario struct {
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
	Puestos  int    `json:"puestos"`
}

type Usuario struct {
	Id       int    `json:"ID"`
	Nombre   string `json:"Nombre"`
	Telefono string `json:"Telefono"`
	Puestos  int    `json:"Puestos"`
}

type FormatoModificarFechas struct {
	Codigo int    `json:"codigo"`
	Fecha  string `json:"fecha"`
}
