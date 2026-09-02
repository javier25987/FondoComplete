package anotaciones

import "errors"

type Service struct {
	_store *Store
}

func NewService(store *Store) *Service {
	return &Service{
		_store: store,
	}
}

func (s *Service) ObtenerAnotaciones(id int) ([]*Anotacion, error) {
	return s._store.ObtenerAnotaciones(id)
}

func (s *Service) CreateAnotacion(id int, nota *PostAnotacion) error {
	if nota.Nota == "" {
		return errors.New("La nota no puede estar vacia")
	}
	if nota.Tipo != "GENERAL" && nota.Tipo != "MONETARIA" {
		return errors.New("El tipo de anotacion no es valido")
	}
	if nota.Tipo == "MONETARIA" && nota.Monto == 0 {
		return errors.New("Una anotacion MONETARIA no puede tener monto igual a 0")
	}
	return s._store.CrearAnotacion(id, nota)
}
