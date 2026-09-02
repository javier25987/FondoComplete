package anotaciones

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

func (h *Handlers) Anotaciones_(w http.ResponseWriter, r *http.Request) {
	Id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "La ruta de busqueda no cuenta con un ID valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		anotaciones, err := h._service.ObtenerAnotaciones(Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(anotaciones)

	case http.MethodPost:
		var nueva_anotacion PostAnotacion
		if err := json.NewDecoder(r.Body).Decode(&nueva_anotacion); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h._service.CreateAnotacion(Id, &nueva_anotacion); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}
