package analizar

import (
	"fmt"
	"net/http"
)

func (h *centralStruct) RutaAnalizarBoletas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		archivoBoletas, err := h.EstadoDeBoletas()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rutaCheque := fmt.Sprintf("./tmp/%s", archivoBoletas)
		http.ServeFile(w, r, rutaCheque)

	default:
		http.Error(w, "Metodo no encontrado", http.StatusMethodNotAllowed)
	}
}
