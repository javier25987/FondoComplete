package cuotas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *centralStruct) RutaCuotas_(w http.ResponseWriter, r *http.Request) {
	Id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "La ruta de busqueda no cuenta con un ID valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.serviceObtenerUsuario(Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var formato FormatoPago
		if err := json.NewDecoder(r.Body).Decode(&formato); err != nil {
			http.Error(w, "No se pudo interpretar el formato de pago", http.StatusBadRequest)
			return
		}
		codigo, err := h.servicePagarSemanas(Id, &formato)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"codigo_cheque": codigo})

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaCuotasCheque_(w http.ResponseWriter, r *http.Request) {
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
