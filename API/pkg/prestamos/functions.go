package prestamos

import (
	"time"
)

func CalendarioMeses(fechaInicio, fechaFinal time.Time) []time.Time {
	mesAct := int(fechaInicio.Month())

	if fechaInicio.After(fechaFinal) || fechaInicio.Equal(fechaFinal) {
		return []time.Time{fechaFinal}
	}

	meses := []time.Time{}

	var (
		mesFor   time.Time
		newMes   time.Time
		mesBucle int
	)

	for i := 1; i <= 24; i++ {
		mesFor = fechaInicio
		newMes = mesFor.AddDate(0, i, 0)

		mesBucle = (mesAct + i) % 12
		if mesBucle == 0 {
			mesBucle = 12
		}

		for mesBucle != int(newMes.Month()) {
			mesFor = mesFor.AddDate(0, 0, -1)
			newMes = mesFor.AddDate(0, i, 0)
		}

		if newMes.After(fechaFinal) || newMes.Equal(fechaFinal) {
			break
		}

		meses = append(meses, newMes)
	}

	if len(meses) == 0 {
		meses = append(meses, fechaFinal)
	}

	return meses
}
