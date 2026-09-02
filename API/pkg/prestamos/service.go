package prestamos

import (
	"API-FONDO/global"
	"errors"
	"fmt"
)

func (s *centralStruct) serviceCrearPrestamo(id int, c *FormatoCrearPrestamo) (int64, error) {
	if c.Motivo != "PG" && c.Motivo != "ACUERDO" { // el motivo no es valido
		return -1, fmt.Errorf("'%s' No es un monto valido", c.Motivo)
	}

	for _, fiador := range c.Fiadores {
		if fiador.Deuda <= 0 {
			return -1, errors.New("La deuda de un fiador no puede ser menor o igual a cero")
		}
	}

	estado, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return -1, fmt.Errorf("Error al confimar si el usuario esta activo >> %v", err)
	}
	if !estado {
		return -1, errors.New("El usuario no esta activo")
	}

	for _, fiador := range c.Fiadores {
		if fiador.Id_user == id {
			return -1, errors.New("El usuario no puede ser su propio fiador")
		}
	}

	fiadores := make(map[int]struct{})
	for _, fiador := range c.Fiadores {
		if _, existe := fiadores[fiador.Id_user]; existe {
			return -1, errors.New("No pueden haber fiadores repetidos")
		}
		fiadores[fiador.Id_user] = struct{}{}
	}

	suma_deudas := 0
	for _, fiador := range c.Fiadores {
		suma_deudas += fiador.Deuda
	}

	capital_disponible, err := s.storeConsultarCapitalDisponibleUsr(id)
	if err != nil {
		return -1, fmt.Errorf("Error al obtener el capital del usuario >> %v", err)
	}

	if c.Monto <= 0 {
		return -1, errors.New("El monto no puede ser menor o igual a cero")
	}
	if suma_deudas > c.Monto {
		return -1, errors.New("Se esta pidiendo en deudas mas de lo que se necesita")
	}

	if capital_disponible > 0 {
		if c.Monto-suma_deudas > capital_disponible {
			return -1, fmt.Errorf("El dinero del usuario junto con los fiadores no alcanza se necesitan %d pesos mas", (c.Monto-suma_deudas)-capital_disponible)
		}
		if suma_deudas+capital_disponible < c.Monto {
			return -1, errors.New("No alcanza para solicitar el prestamo, solicite mas fiadores")
		}
	} else {
		if len(c.Fiadores) == 0 {
			return -1, errors.New("El usuario no tiene capital disponible, solicite fiadores")
		}
		if suma_deudas < c.Monto {
			return -1, errors.New("El dinero de los fiadores no alcanza")
		}
	}

	var (
		capital_fiador int
		estado_fiador  bool
	)

	for _, fiador := range c.Fiadores {
		capital_fiador, err = s.storeConsultarCapitalDisponibleUsr(fiador.Id_user)
		if err != nil {
			return -1, fmt.Errorf("Error al obtener el capital del fiador >> %v", err)
		}
		estado_fiador, err = global.UserIsActive(s.Conexion, fiador.Id_user)
		if err != nil {
			return -1, fmt.Errorf("Error al confirmar si el fiador esta activo >> %v", err)
		}

		if capital_fiador < fiador.Deuda {
			return -1, fmt.Errorf("El fiador %d no cuenta con capital para este prestamo", fiador.Id_user)
		}
		if !estado_fiador {
			return -1, fmt.Errorf("El usuario %d no esta activo", fiador.Id_user)
		}
	}

	return s.storeCargarPrestamo(id, c)
}

func (s *centralStruct) serviceAbonarPrestamo(id int, a *AbonarPrestamo) error {
	estadoPrestamo, err := s.storeObtenerEstadoDelPresamo(a.Codigo)
	if err != nil {
		return fmt.Errorf("Rectificar: %v", err)
	}
	if !estadoPrestamo {
		return fmt.Errorf("Rectificar: el prestamo no esta activo (fue pagado del todo o desactivado)")
	}

	if a.Monto <= 0 {
		return errors.New("Rectificar: el monto no puede ser menor o igual a cero")
	}

	deuda, intereses, err := s.storeObtenerDeudaEIntereses(a.Codigo)
	if err != nil {
		return fmt.Errorf("Rectificar: %v", err)
	}

	switch a.Tipo {
	case "DEUDA":
		if a.Monto > deuda {
			return errors.New("Rectificar: no se puede pagar mas de lo que se debe")
		}
	case "INTERESES":
		if a.Monto > intereses {
			return errors.New("Rectificar: no se puede pagar mas intereses de lo que se debe")
		}
	default:
		return fmt.Errorf("Rectificar: '%s' No es un tipo de pago valido", a.Tipo)
	}

	propietario, err := s.storeObtenerPropietario(a.Codigo)
	if err != nil {
		return fmt.Errorf("Rectificar: %v", err)
	}
	if id != propietario {
		return fmt.Errorf("Rectificar: %d No es el propietario del prestamo", id)
	}

	estado, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return fmt.Errorf("Rectificar: %v", err)
	}

	if !estado {
		return errors.New("Rectificar: el usuario no esta activo")
	}

	return s.storeAbonarPrestamo(id, a)
}

func (s *centralStruct) serviceObtenerDatosUsuario(id int) (*GetUsuario, error) {
	err := s.storeArreglarAsuntos(id)
	if err != nil {
		return nil, err
	}
	return s.storeGetDataUser(id)
}
