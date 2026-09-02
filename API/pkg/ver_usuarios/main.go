package verusuaios

import (
	"API-FONDO/ajustes"
	"database/sql"
)

type centralStruct struct {
	Ajustes  *ajustes.Ajustes
	Conexion *sql.DB
}

func NewCentralStruct(ajustes *ajustes.Ajustes, conexion *sql.DB) *centralStruct {
	return &centralStruct{
		Conexion: conexion,
		Ajustes:  ajustes,
	}
}
