package ajustes_http

import (
	"fmt"
	"time"
)

/*
Las fechas en este formato deben ser dadadas en el formato "AÑO-MES-DIA",
la hora debe ser dada por separado en el formato "HORA:MINUTOS"
*/
func fnCrearCalendario(fechaInicial, horaDeCorte string, fechasDobles []string) ([]time.Time, error) {
	formatoDeFechas := "2006-01-02 15:04"
	fechaStr := fechaInicial + " " + horaDeCorte

	formatoValidoDeFecha := "formato valido, fecha:'AÑO-MES-DIA', hora de corte:'HORA:MINUTOS'"

	fechaPivote, err := time.Parse(formatoDeFechas, fechaStr)
	if err != nil {
		return nil, fmt.Errorf(
			"La fecha '%s' no puede ser formateada. %s >> %q",
			fechaStr, formatoValidoDeFecha, err,
		)
	}

	mapFechasDobles := map[time.Time]struct{}{}

	var fechaTmp time.Time

	for _, value := range fechasDobles {
		fechaStr = value + " " + horaDeCorte

		fechaTmp, err = time.Parse(formatoDeFechas, fechaStr)
		if err != nil {
			return nil, fmt.Errorf(
				"La fecha '%s' no puede ser formateada. %s >> %q",
				fechaStr, formatoValidoDeFecha, err,
			)
		}

		mapFechasDobles[fechaTmp] = struct{}{}
	}

	fechas := []time.Time{fechaPivote}

	i := 1
	for len(fechas) < 50 {
		fechaTmp = fechaPivote.AddDate(0, 0, 7*i)
		fechas = append(fechas, fechaTmp)
		i += 1

		if _, ok := mapFechasDobles[fechaTmp]; ok {
			fechas = append(fechas, fechaTmp)
		}
	}

	fechas = fechas[:50]

	return fechas, nil
}
