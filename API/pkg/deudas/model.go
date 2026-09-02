package deudas

type PostDeuda struct {
	Monto  int    `json:"monto"`
	Motivo string `json:"motivo"`
	Tipo   string `json:"tipo"`
}

type PutDeudas struct {
	Codigo int `json:"codigo"`
}

type Deuda struct {
	Codigo        int    `json:"codigo"`
	Monto         int    `json:"monto"`
	FechaCreacion string `json:"fecha_de_creacion"`
	FechaPago     string `json:"fecha_de_pago"`
	Tipo          string `json:"tipo"`
	Motivo        string `json:"motivo"`
	Estado        bool   `json:"estado"`
}

type GetDeudas struct {
	Id     int      `json:"id"`
	Nombre string   `json:"nombre"`
	Deudas []*Deuda `json:"deudas"`
}
