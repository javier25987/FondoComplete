package rifas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *centralStruct) RutaAjustesBoletas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rifas, err := h.storeObtenerRifasCreadas()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rifas)

	case http.MethodPut:
		var formato_nueva_rifa FormatoDatosDeRifa
		if err := json.NewDecoder(r.Body).Decode(&formato_nueva_rifa); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nuevo_codigo, err := h.serviceCrearUnaRifa(&formato_nueva_rifa)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mensaje := fmt.Sprintf("La rifa ha sido creada. Codigo de esta: %d", nuevo_codigo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"mensaje": mensaje,
		})

	default:
		http.Error(w, "El metodo requerido no existe", http.StatusMethodNotAllowed)
	}
}

func (h *centralStruct) RutaInsertar(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var boletas []Boleta
		if err := json.NewDecoder(r.Body).Decode(&boletas); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.serviceInsertarBoletas(boletas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodGet:
		boletas, err := h.storeBoletasLibres()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(boletas)

	case http.MethodDelete:
		boletaStr := r.URL.Query().Get("codigo")
		boletaInt, err := strconv.Atoi(boletaStr)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("No se puede formatear '%s' como un numero >> %q", boletaStr, err),
				http.StatusBadRequest,
			)
			return
		}

		err = h.storeEliminarBoleta(boletaInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
		return
	}
}

func (h *centralStruct) RutaRifas(w http.ResponseWriter, r *http.Request) {
	Id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "La ruta de busqueda no cuenta con un ID valido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.storeObtenerDatosDelUsuario(Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

	case http.MethodPost:
		var boletas BoletasAEntregar
		if err := json.NewDecoder(r.Body).Decode(&boletas); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.serviceEntregarBoletas(Id, boletas.Boletas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodPut:
		var formatoMonto FormatoMontoAPagar
		if err := json.NewDecoder(r.Body).Decode(&formatoMonto); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.servicePagarDeudaPorBoletas(Id, formatoMonto.Monto); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
		return
	}
}
