package prestamos

type FormatoCrearPrestamo struct {
	Monto    int       `json:"monto"`
	Motivo   string    `json:"motivo"`
	Fiadores []*Fiador `json:"fiadores"`
}

type AbonarPrestamo struct {
	Monto  int    `json:"monto"`
	Codigo int    `json:"codigo"`
	Tipo   string `json:"tipo"`
	Fecha  string `json:"fecha"`
}

type Fiador struct {
	Id_user int `json:"id_user"`
	Deuda   int `json:"deuda"`
}

type Prestamo struct {
	Codigo          int       `json:"codigo"`
	Interes         int       `json:"interes"`
	InteresVencido  int       `json:"interes_vencido"`
	InteresGenerado int       `json:"interes_generado"`
	Deuda           int       `json:"deuda"`
	Monto           int       `json:"monto"`
	DeudaActual     int       `json:"deuda_actual"`
	Fiadores        []*Fiador `json:"fiadores"`
	FechasPago      []string  `json:"fechas_de_pago"`
	Estado          string    `json:"estado"`
	Motivo          string    `json:"motivo"`
	Anotaciones     []string  `json:"anotaciones"`
}

type DeudaFiador struct {
	Cuanto     int `json:"cuanto"`
	EnPrestamo int `json:"en_prestamo"`
}

type DeudaPrestamo struct {
	Prestamo int `json:"prestamo"`
	Deuda    int `json:"deuda"`
}

type ConsultaCapital struct {
	CapitalPago            int              `json:"capital_pago"`
	CapitalDisponible      int              `json:"capital_disponible"`
	DeudasPorFiador        []*DeudaFiador   `json:"deudas_por_fiador"`
	TotalDeudasPorFiador   int              `json:"total_deudas_por_fiador"`
	DeudasEnPrestamos      []*DeudaPrestamo `json:"deudas_en_prestamos"`
	TotalDeudasEnPrestamos int              `json:"total_deudas_en_prestamos"`
	TotalDisponible        int              `json:"total_disponible"`
}

type ChequeJson struct {
	Codigo         int    `json:"codigo"`
	CodigoPrestamo int    `json:"codigo_prestamo"`
	Monto          int    `json:"monto"`
	DeudaAntes     int    `json:"deuda_antes"`
	DeudaDespues   int    `json:"deuda_despues"`
	Fecha          string `json:"fecha"`
}

type GetUsuario struct {
	Id                int              `json:"id"`
	Nombre            string           `json:"nombre"`
	PrestamosUser     []int            `json:"prestamos_user"`
	Prestamos         []*Prestamo      `json:"prestamos"`
	CapitalDisponible *ConsultaCapital `json:"capital_disponible"`
	ChequesUser       []int            `json:"cheques_user"`
	Cheques           []*ChequeJson    `json:"cheques"`
}

type ChequePdf struct {
	Nombre           string
	Id               int
	Codigo           int
	CodigoPrestamo   int
	MontoAbonado     int
	AbonoA           string
	DeudaAntes       int
	DeudaDespues     int
	InteresesAntes   int
	InteresesDespues int
	FechaYHora       string
}
