package cuotas

import (
	"API-FONDO/global"
	"errors"
)

func (s *centralStruct) serviceObtenerUsuario(id int) (*User, error) {
	if id < 1 || id > s.Ajustes.Usuarios.Count {
		return nil, errors.New("El ID proporcionado no pertenece a ningun usuario del sistema")
	}

	err := s.storeRectificarAsuntos(id)
	if err != nil {
		return nil, err
	}

	user, err := s.storeObtenerDatosUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *centralStruct) servicePagarSemanas(id int, formato *FormatoPago) (int, error) {
	if formato.Cuotas < 0 || formato.Multas < 0 {
		return -1, errors.New("El número de cuotas o multas a pagar no puede ser negativo")
	}

	if formato.Cuotas == 0 && formato.Multas == 0 {
		return -1, errors.New("Debe pagar al menos una cuota o multa")
	}

	if formato.Metodo == "" {
		return -1, errors.New("Debe especificar un método de pago")
	}

	if formato.Metodo != "EFECTIVO" && formato.Metodo != "TRANSFERENCIA" {
		return -1, errors.New("Método de pago no reconocido use 'EFECTIVO' o 'TRANSFERENCIA'")
	}

	is_active, err := global.UserIsActive(s.Conexion, id)
	if err != nil {
		return -1, errors.New("No se pudo verificar el estado del usuario")
	}

	if !is_active {
		return -1, errors.New("El usuario no esta activo")
	}

	return s.storePagarSemanas(id, formato)
}

// func (s *Service) CrearChequeUser(codigo int) error {
// 	return s._store.CrearChequeUser(codigo)
// }
