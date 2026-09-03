package analizar

type FormatoDeEntregas struct {
	Codigo int
	Fecha  string
}

type FormatoDeFechaYBoletas struct {
	Fecha   string   `json:"fecha"`
	Boletas []string `json:"boletas"`
}

type FormatoDeUsuarioYBoletas struct {
	Puesto   int                      `json:"puesto"`
	Nombre   string                   `json:"nombre"`
	Entregas []FormatoDeFechaYBoletas `json:"entregas"`
}

type FormatoDeConsultaCompletaDeBoletas struct {
	RifaAct       int                        `json:"rifa_actual"`
	Usuarios      []FormatoDeUsuarioYBoletas `json:"usuarios"`
	BoletasLibres []string                   `json:"boletas_libres"`
}
