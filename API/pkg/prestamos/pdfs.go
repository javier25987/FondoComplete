package prestamos

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
		#raw("abono a prestamo")
	]
	#raw("Nombre:") #h(1fr) #raw(%q)\
	#raw("Puesto:") #h(1fr) #raw("%d")\
	#line()
	#raw("Monto Abonado:") #h(1fr) #raw(%q)\
	#raw("Abono a:") #h(1fr) #raw(%q)\
	#line()
	#raw("Estado de la Deuda:")\
	#raw(%q) #h(1fr) $->$ #h(1fr) #raw(%q)\
	#raw("Estado de los Intereses:")\
	#raw(%q) #h(1fr) $->$ #h(1fr) #raw(%q)\
	#line()
	#raw("Codigo del Prestamo:") #h(1fr) #raw("%d")\
	#line()
	#raw("Fecha:") #h(1fr) #raw(%q)
`

func (s *centralStruct) CrearUnChequeParaGet(codigo int) (string, error) {
	// crear el nombre del archivo temporal
	nombreArchivo := time.Now().Format("0102150405") + ".pdf"

	// obtener los datos del cheque
	datosCheque, err := s.storeObtenerDatosChequePrestamos(codigo)
	if err != nil {
		return "", err
	}

	// formatear el cheque
	chequeFormateado := fmt.Sprintf(
		CHEQUE,
		datosCheque.Codigo,
		datosCheque.Nombre,
		datosCheque.Id,
		global.FormatInt(datosCheque.MontoAbonado),
		datosCheque.AbonoA,
		global.FormatInt(datosCheque.DeudaAntes),
		global.FormatInt(datosCheque.DeudaDespues),
		global.FormatInt(datosCheque.InteresesAntes),
		global.FormatInt(datosCheque.InteresesDespues),
		datosCheque.CodigoPrestamo,
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
