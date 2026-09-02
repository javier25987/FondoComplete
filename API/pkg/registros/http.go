package registros

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handlers struct {
	_service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{
		_service: service,
	}
}

func (h *Handlers) Registros_(w http.ResponseWriter, r *http.Request) {
	Mes, err := strconv.Atoi(r.PathValue("mes"))
	if err != nil {
		http.Error(w, "El mes proporcionado en la ruta no es valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tabla, err := h._service.ObtenerRegistro(Mes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tabla)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}
