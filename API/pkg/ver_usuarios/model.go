package verusuaios

import "API-FONDO/pkg/prestamos"

type Usuario struct {
	Id       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
	Puestos  string `json:"puestos"`
	Estado   string `json:"estado"`
}

type Acuerdo struct {
	Id               int    `json:"Id"`
	Nombre           string `json:"Nombre"`
	CapitalPago      int    `json:"Pago"`
	CapitalRetirado  int    `json:"Retirado"`
	NecesarioRetirar int    `json:"Necesario Retirar"`
}

type Boleta struct {
	Id      int    `json:"Id"`
	Nombre  string `json:"Nombre"`
	Numeros string `json:"Boleta"`
}

type GetPrestamo struct {
	Id       int                 `json:"id"`
	Nombre   string              `json:"nombre"`
	Prestamo *prestamos.Prestamo `json:"prestamo"`
}
