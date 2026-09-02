package registros

type GetMes struct {
	Mes int `json:"mes"`
}

type TablaMeses struct {
	Fecha      []string `json:"Fecha"`
	Ingreso    []int    `json:"Ingreso"`
	Egreso     []int    `json:"Egreso"`
	Diferencia []int    `json:"Diferencia"`
	Muestra    []string `json:"Muestra"`
}
