package rifas

type FormatoDatosDeRifa struct {
	Codigo int `json:"codigo"`

	NumeroDeBoletas        int      `json:"numero_de_boletas"`
	Premios                []Premio `json:"premios"`
	CostoDeBoletas         int      `json:"costo_de_boletas"`
	CostosDeAdministracion int      `json:"costos_de_administracion"`

	GananciaPorBoleta int `json:"ganancia_por_boleta"`
}

type Premio struct {
	Numero int `json:"numero"`
	Premio int `json:"premio"`
}

type Boleta struct {
	Codigo int    `json:"codigo"`
	Boleta string `json:"boleta"`
}

type FormatoGet struct {
	Puesto             int      `json:"puesto"`
	Nombre             string   `json:"nombre"`
	BoletasDisponibles []int    `json:"boletas_disponibles"`
	DeudaPorBoletas    int      `json:"deuda_por_boletas"`
	BoletasEntregadas  []string `json:"boletas_entregadas"`
}

type BoletasAEntregar struct {
	Boletas []int `json:"boletas"`
}

type FormatoMontoAPagar struct {
	Monto int `json:"monto"`
}
