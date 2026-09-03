package rifas

import (
	"API-FONDO/global"
	"errors"
	"fmt"
)

func (s *centralStruct) serviceCrearUnaRifa(frmt *FormatoDatosDeRifa) (int64, error) {
	// confirmacciones de que todos los valores sean correctos
	if frmt.NumeroDeBoletas < 50 {
		return -1, fmt.Errorf("rifas >> crear-una-rifa service >> %d son muy pocas boletas", frmt.NumeroDeBoletas)
	}
	if len(frmt.Premios) == 0 {
		return -1, fmt.Errorf("rifas >> crear-una-rifa service >> tiene que haber al menos un premio")
	}
	if frmt.CostoDeBoletas < 1 {
		return -1, fmt.Errorf("rifas >> crear-una-rifa service >> los costos de boletas no pueden ser nulos o negativos")
	}
	if frmt.CostosDeAdministracion < 1 {
		return -1, fmt.Errorf("rifas >> crear-una-rifa service >> los costos de administracion no pueden ser nulos o negativos")
	}

	return s.storeCrearUnaRifa(frmt)
}

func (s *centralStruct) serviceInsertarBoletas(boletas []Boleta) error {
	/*
		aca me voy a marcar un gran yo creo y lo unico que voy a certificar es que al menos se envie una boleta
		la verdad que todas las boletas tengan el formato correcto 000-000-000-000 es algo que se lo dejo al
		frontend
	*/

	if len(boletas) < 1 {
		return fmt.Errorf("rifas >> insertar boletas >> no hay boletas a registrar")
	}

	return s.storeInsertarBoletas(boletas)
}

func (s *centralStruct) serviceEntregarBoletas(id int, boletas []int) error {
	/*
		FALTA LA PROTECCION PARA NO REENTREGAR BOLETAS YA ENTREGADAS
	*/

	estado, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return fmt.Errorf("No se pudo rectificar el estado del usuario >> %v", err)
	}

	if !estado {
		return fmt.Errorf("El usuario %d no esta activo", id)
	}

	if len(boletas) == 0 {
		return fmt.Errorf("No hay boletas a entregar")
	}

	return s.storeEntregarBoletas(id, boletas)
}

func (s *centralStruct) servicePagarDeudaPorBoletas(id, monto int) error {
	if monto <= 0 {
		return errors.New("No se puede pagar menos a cero o cero en si mismo")
	}

	estado, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return fmt.Errorf("No se pudo rectificar el estado del usuario >> %v", err)
	}

	if !estado {
		return fmt.Errorf("El usuario %d no esta activo", id)
	}

	deuda, err := s.storeObtenerDeudasPorBoletas(id)
	if err != nil {
		return err
	}

	if monto > deuda {
		return errors.New("No se puede pagar mas de lo que se debe")
	}

	return s.storePagarDeudaPorBoletas(id, monto)
}
