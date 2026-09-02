package prestamos

import (
	"API-FONDO/ajustes"
	"database/sql"
)

type centralStruct struct {
	Conexion *sql.DB
	Ajustes  *ajustes.Ajustes
}

func NewCentralStruct(ajustes *ajustes.Ajustes, conexion *sql.DB) *centralStruct {
	return &centralStruct{
		Conexion: conexion,
		Ajustes:  ajustes,
	}
}
