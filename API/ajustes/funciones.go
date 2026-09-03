package ajustes

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const ConfigPath = "./setings.toml"

func AbrirAjustes() (*Ajustes, error) {
	ajustes := Ajustes{}

	if _, err := toml.DecodeFile(ConfigPath, &ajustes); err != nil {
		return nil, fmt.Errorf("No se pueden abrir los ajustes >> %q", err)
	}

	return &ajustes, nil
}

func (a *Ajustes) IncrementarUsuarios() error {
	nuevosAjustes := *a
	nuevosAjustes.Usuarios.Count++

	if err := a.GuardarAjustes(nuevosAjustes); err != nil {
		return err
	}

	return nil
}

func (a *Ajustes) GuardarAjustes(nuevosAjustes Ajustes) error {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)

	if err := encoder.Encode(nuevosAjustes); err != nil {
		return fmt.Errorf("No se pueden formatear los nuevos ajustes >> %q", err)
	}

	if err := os.WriteFile(ConfigPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("Error al guardar los ajustes >> %q", err)
	}

	*a = nuevosAjustes

	return nil
}

func CrearAjustes() {
	ajustes := Ajustes{
		Cuotas: Cuotas{
			ValorCuota:       10_000,
			ValorMulta:       3_000,
			CobrarMultas:     false,
			BloquearUsuarios: false,
		},
		Intereses: Intereses{
			InteresMTope:  2,
			InteresmTope:  3,
			TopeIntereses: 20_000_000,
		},
		Prestamos: Prestamos{
			CapitalUsable: 75,
			FechaDeCierre: time.Date(2029, time.May, 21, 23, 59, 0, 0, time.UTC),
		},
		Rifas: Rifas{
			RifaAct: 1,
		},
		Usuarios: Usuarios{
			Count: 0,
		},
		BaseDeDatos: BaseDatos{
			Db: "Fondo.db",
		},
	}

	ajustes.GuardarAjustes(ajustes)
}
