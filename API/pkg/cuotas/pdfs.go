package cuotas

import (
	"API-FONDO/global"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const CHEQUE string = `
	#set page(width: 8cm, height: 15cm)
	#set text(size: 13pt)

	#align(center)[
		= #raw("Fondo San Javier")
		#raw("cheque de pago №:") #raw("%d")\
		#raw("pago de cuotas")
	]

	#raw("Nombre:") #h(1fr) #raw("%v")\
	#raw("Puesto:") #h(1fr) #raw("%v")\
	#raw("Puestos:") #h(1fr) #raw("%d")
	#line()
	#raw("Multas Pagas:") #h(1fr) #raw("%v")\
	#raw("Valor Por Multa:") #h(1fr) #raw("%s")\
	#raw("Total:") #h(1fr) #raw("%s")\
	#line()
	#raw("Cuotas Pagas:") #h(1fr) #raw("%v")\
	#raw("Valor Por Cuota:") #h(1fr) #raw("%s")\
	#raw("Total:") #h(1fr) #raw("%s")\
	#line()
	#raw("Metodo De Pago:") #h(1fr) #raw("%v")\
	#raw("TOTAL PAGADO:") #h(1fr) #raw("%s")\
	#line()
	#raw("Fecha:") #h(1fr) #raw("%v")\
`

func (s *centralStruct) CrearUnChequeParaGet(codigo int) (string, error) {
	// crear el nombre del archivo temporal
	nombreArchivo := time.Now().Format("0102150405") + ".pdf"

	// obtener los datos del cheque
	datosCheque, err := s.ObtenerDatosChequeCuotas(codigo)
	if err != nil {
		return "", err
	}

	// formatear el cheque
	chequeFormateado := fmt.Sprintf(
		CHEQUE,
		datosCheque.Codigo,
		datosCheque.Nombre,
		datosCheque.Id,
		datosCheque.Puestos,
		datosCheque.MultasPagas,
		global.FormatInt(datosCheque.ValorMulta),
		global.FormatInt(datosCheque.TotalMultas),
		datosCheque.CuotasPagas,
		global.FormatInt(datosCheque.ValorCuota),
		global.FormatInt(datosCheque.TotalCuotas),
		datosCheque.MetodoPago,
		global.FormatInt(datosCheque.TotalPagado),
		datosCheque.FechaYHora,
	)

	basePath, _ := os.Getwd()
	rutaArchivo := filepath.Join(basePath, "tmp", nombreArchivo)

	if err := os.MkdirAll(filepath.Join(basePath, "tmp"), 0755); err != nil {
		return "", fmt.Errorf("no se pudo crear tmp: %w", err)
	}

	cmd := exec.Command("typst", "compile", "-", rutaArchivo)
	cmd.Stdin = strings.NewReader(chequeFormateado)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("no se pudo compilar el cheque (%s): %v\nSalida: %s", rutaArchivo, err, string(output))
	}

	return nombreArchivo, nil
}
