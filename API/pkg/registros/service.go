package registros

type Service struct {
	_store *Store
}

func NewService(store *Store) *Service {
	return &Service{
		_store: store,
	}
}

func (s *Service) ObtenerRegistro(mes int) (*TablaMeses, error) {
	return s._store.GetMes(mes)
}
