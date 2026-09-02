package ajustes

import "time"

type Cuotas struct {
	ValorCuota       int  `toml:"valor_cuota" json:"valor_cuota"`
	ValorMulta       int  `toml:"valor_multa" json:"valor_multa"`
	CobrarMultas     bool `toml:"cobrar_multas" json:"cobrar_multas"`
	BloquearUsuarios bool `toml:"bloquear_usuarios" json:"bloquear_usuarios"`
}

type Intereses struct {
	InteresMTope  int `toml:"interes_M_tope" json:"interes_M_tope"`
	InteresmTope  int `toml:"interes_m_tope" json:"interes_m_tope"`
	TopeIntereses int `toml:"tope_intereses" json:"tope_intereses"`
}

type Prestamos struct {
	CapitalUsable int       `toml:"capital_usable" json:"capital_usable"`
	FechaDeCierre time.Time `toml:"fecha_de_cierre" json:"fecha_de_cierre"`
}

type Rifas struct {
	RifaAct int `toml:"rifa_act" json:"rifa_act"`
}

type Usuarios struct {
	Count int `toml:"count" json:"count"`
}

type BaseDatos struct {
	Db string `toml:"db_name" json:"db_name"`
}

type Ajustes struct {
	Cuotas      Cuotas    `toml:"cuotas" json:"cuotas"`
	Intereses   Intereses `toml:"intereses" json:"intereses"`
	Prestamos   Prestamos `toml:"prestamos" json:"prestamos"`
	Rifas       Rifas     `toml:"rifas" json:"rifas"`
	Usuarios    Usuarios  `toml:"usuarios" json:"usuarios"`
	BaseDeDatos BaseDatos `toml:"db" json:"db"`
}
