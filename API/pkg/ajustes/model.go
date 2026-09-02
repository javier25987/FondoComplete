package ajustes_http

type ToGetAjuste struct {
	Ajuste string `json:"ajuste"`
}

type FormatAjuste struct {
	Nombre string `json:"ajuste"` // se llama Nombre ya que no quiero mesclar el campo y la estructura
	// de todas manera cuando parsee a json se va allamar ajuste
	ValorN int    `json:"nvalor"`
	ValorT string `json:"tvalor"`
}

type AjustesResponse struct {
	Ajustes        []*FormatAjuste `json:"ajustes"`
	NombresAjustes []string        `json:"nombres_ajustes"`
}

/*
	Formato del calendario
*/

type FormatNewCalendar struct {
	FechaInicial string   `json:"fecha_inicial"`
	FechasDobles []string `json:"fechas_dobles"`
	HoraDeCierre string   `json:"hora_de_cierre"`
}

type Semana struct {
	Numero int    `json:"Semana"`
	Fecha  string `json:"Fecha"`
}
