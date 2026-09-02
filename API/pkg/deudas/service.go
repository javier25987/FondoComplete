package deudas

import (
	"API-FONDO/global"
	"errors"
	"fmt"
)

func (s *centralStruct) serviceCrearDeuda(id int, deuda *PostDeuda) error {
	estado, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return fmt.Errorf("No se ppudo obtener el estado del usuario %d >> %v", id, err)
	}
	if !estado {
		return fmt.Errorf("El usuario %d no esta activo", id)
	}

	if deuda.Monto <= 0 {
		return errors.New("El monto debe ser mayor a cero")
	}
	if deuda.Motivo == "" {
		return errors.New("El motivo no puede estar vacio")
	}
	if deuda.Tipo == "" {
		return errors.New("El tipo no puede estar vacio")
	}
	if deuda.Tipo != "MULTA POR CABEZA" && deuda.Tipo != "MULTA EXTRA" {
		return errors.New("El tipo de deuda no es valido use 'MULTA POR CABEZA' o 'MULTA EXTRA'")
	}

	return s.storeCrearDeuda(id, deuda)
}

func (s *centralStruct) servicePagarDeuda(id, codigo int) error {
	estado, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return fmt.Errorf("No se ppudo obtener el estado del usuario %d >> %v", id, err)
	}
	if !estado {
		return fmt.Errorf("El usuario %d no esta activo", id)
	}

	propietario, err := s.storeObtenerPropietario(codigo)
	if err != nil {
		return fmt.Errorf("Revisar: %v", err)
	}
	if id != propietario {
		return fmt.Errorf("El usuario %d no es el propietario de la deuda", id)
	}

	estado_deuda, err := s.storeObtenerEstado(codigo)
	if err != nil {
		return fmt.Errorf("Revisar: %v", err)
	}
	if !estado_deuda {
		return errors.New("La deuda que se desea pagar ya esta paga")
	}

	return s.storePagarDeuda(codigo)
}
