package modificar

import (
	"encoding/json"
	"net/http"
)

func (h *centralStruct) RutaModificar(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tabla, err := h.serviceObtenerUsuarios()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tabla)

	case http.MethodPost:
		var nuevo_usuario FormatoCrearUsuario
		if err := json.NewDecoder(r.Body).Decode(&nuevo_usuario); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := h.serviceCrearUsuario(&nuevo_usuario)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"nuevo_id": id})

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaModificarFechas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		var formato FormatoModificarFechas
		if err := json.NewDecoder(r.Body).Decode(&formato); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.storeCambiarFechasDeUnPrestamo(formato.Codigo, formato.Fecha); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
