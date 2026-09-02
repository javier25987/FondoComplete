package verusuaios

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *centralStruct) RutaVerUsuarios(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		llave := r.URL.Query().Get("nombre")

		tabla, err := h.storeBuscarUsuarios(llave)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tabla)
	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaVerAcuerdo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tabla, err := h.storeFiltrarAcuerdo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tabla)
	default:
		http.Error(w, "Metodo no encotrado", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaVerBoletas(w http.ResponseWriter, r *http.Request) {
	rifa, err := strconv.Atoi(r.PathValue("rifa"))
	if err != nil {
		http.Error(w, "La ruta no cuenta con una rifa valida como numero", http.StatusBadRequest)
		return
	}

	numero := r.URL.Query().Get("boleta")
	if numero == "" {
		http.Error(w, "La boleta a buscar esta vacia", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tabla, err := h.storeBuscarBoletas(rifa, numero)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tabla)
	default:
		http.Error(w, "Metodo no encotrado", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaVerPrestamo(w http.ResponseWriter, r *http.Request) {
	codigo, err := strconv.Atoi(r.PathValue("codigo"))
	if err != nil {
		http.Error(w, "La ruta no cuenta con un codigo valido como numero", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tabla, err := h.storeObtenerPrestamo(codigo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tabla)
	default:
		http.Error(w, "Metodo no encotrado", http.StatusMethodNotAllowed)
	}
}
