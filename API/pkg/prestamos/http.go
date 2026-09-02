package prestamos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *centralStruct) RutaPrestamos_(w http.ResponseWriter, r *http.Request) {
	Id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "La ruta de busqueda no cuenta con un ID valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.serviceObtenerDatosUsuario(Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var abonar *AbonarPrestamo
		if err := json.NewDecoder(r.Body).Decode(&abonar); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.serviceAbonarPrestamo(Id, abonar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodPost:
		var cargar *FormatoCrearPrestamo
		if err := json.NewDecoder(r.Body).Decode(&cargar); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		codigo, err := h.serviceCrearPrestamo(Id, cargar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{"codigo": codigo})

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaPrestamosCheque_(w http.ResponseWriter, r *http.Request) {
	Codigo, err := strconv.Atoi(r.PathValue("codigo"))
	if err != nil {
		http.Error(w, "La ruta de busqueda no cuenta con un Codigo valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		archivoCheque, err := h.CrearUnChequeParaGet(Codigo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rutaCheque := fmt.Sprintf("./tmp/%s", archivoCheque)
		http.ServeFile(w, r, rutaCheque)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}
