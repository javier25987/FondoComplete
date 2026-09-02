package cuotas

type Semana struct {
	/*
		En esta estructura los nombres estan escritos con mayusculas ya que estan
		destinados a ser los nombres de una tabla en el frontend
	*/
	Semana      int    `json:"Semana"`
	Fecha       string `json:"Fecha"`
	Cuota       string `json:"Cuota"`
	MultaActiva string `json:"Multas Activas"`
	MultaPaga   string `json:"Multas Pagas"`
}

type User struct {
	Id           int           `json:"id"`
	Nombre       string        `json:"nombre"`
	Puestos      int           `json:"puestos"`
	Telefono     string        `json:"telefono"`
	Tabla        []*Semana     `json:"tabla"`
	ValorCuota   int           `json:"valor_cuota"`
	ValorMulta   int           `json:"valor_multa"`
	Cheques      []*ChequeJson `json:"cheques"`
	ChequesUser  []int         `json:"cheques_user"`
	CuotasAPagar int           `json:"cuotas_a_pagar"`
	MultasAPagar int           `json:"multas_a_pagar"`
}

type FormatoPago struct {
	Cuotas int    `json:"cuotas"`
	Multas int    `json:"multas"`
	Metodo string `json:"metodo"`
}

type ChequeJson struct {
	Codigo      int    `json:"codigo"`
	CuotasPagas int    `json:"cuotas_pagas"`
	MultasPagas int    `json:"multas_pagas"`
	MetodoPago  string `json:"metodo_pago"`
	FechaYHora  string `json:"fecha_y_hora"`
}

type ChequePdf struct {
	Nombre      string
	Id          int
	Codigo      int
	Puestos     int
	CuotasPagas int
	ValorCuota  int
	TotalCuotas int
	MultasPagas int
	ValorMulta  int
	TotalMultas int
	MetodoPago  string
	TotalPagado int
	FechaYHora  string
}
