package ajustes_http

import (
	"API-FONDO/global"
	"encoding/json"
	"net/http"
)

// func (h *centralStruct) Ajuste(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		ajuste := r.URL.Query().Get("ajuste")

// 		ajuste_obtenido, err := h._service.ObtenerAjuste(ajuste)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(ajuste_obtenido)

// 	case http.MethodPut:
// 		var ajuste Ajuste

// 		if err := json.NewDecoder(r.Body).Decode(&ajuste); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		if err := h._service.ModificarAjuste(&ajuste); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)

// 	case http.MethodPost:
// 		var ajuste Ajuste

// 		if err := json.NewDecoder(r.Body).Decode(&ajuste); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		if err := h._service.CrearAjuste(&ajuste); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)

// 	default:
// 		http.Error(w, "Metodo no encontrado", http.StatusBadRequest)

// 	}
// }

func (h *centralStruct) RutaAjustes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h.Ajustes)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusBadRequest)
	}
}

func (h *centralStruct) RutaAjusteCalendario(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		calendarioMap, err := global.ObtenerCalendarioStr(h.Conexion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		calendario := []Semana{}
		for i := range 51 {
			if value, ok := calendarioMap[i]; ok {
				calendario = append(calendario, Semana{Numero: i, Fecha: value[:10]})
			}
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(calendario)

	case http.MethodPost:
		var nuevasFechas FormatNewCalendar
		if err := json.NewDecoder(r.Body).Decode(&nuevasFechas); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.storeCrearCalendario(nuevasFechas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
