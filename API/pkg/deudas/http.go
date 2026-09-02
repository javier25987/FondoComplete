package deudas

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *centralStruct) RutaDeudas_(w http.ResponseWriter, r *http.Request) {
	Id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "La ruta de busqueda no cuenta con un ID valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var deuda *PostDeuda
		if err := json.NewDecoder(r.Body).Decode(&deuda); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.serviceCrearDeuda(Id, deuda); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodPut:
		codigo, err := strconv.Atoi(r.URL.Query().Get("codigo"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.servicePagarDeuda(Id, codigo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		get_deuda, err := h.storeObtenerDeudas(Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(get_deuda)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
		return
	}
}
