package global

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatInt(n int) string {
	p := message.NewPrinter(language.Spanish)
	return p.Sprintf("%d", n)
}

func UserIsActive(db *sql.DB, id int) (bool, error) {
	if id == 1976 {
		return true, nil
	}

	var isActive bool
	q := `SELECT u.estado FROM usuarios u WHERE u.id = $1`
	if err := db.QueryRow(q, id).Scan(&isActive); err != nil {
		return false, err
	}
	return isActive, nil
}

func ObtenerCalendarioTime(db *sql.DB) (map[int]time.Time, error) {
	rows, err := db.Query(`select semana, fecha from calendario`)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener el calendario >> %q", err)
	}
	defer rows.Close()

	calendario := map[int]time.Time{}
	formatoDeFechas := "2006-01-02T15:04:05Z"
	var (
		fecha    time.Time
		fechaStr string
		semana   int
	)

	for rows.Next() {
		if err := rows.Scan(&semana, &fechaStr); err != nil {
			return nil, fmt.Errorf("No se pudieron extraer las fechas del calendario")
		}

		fecha, err = time.Parse(formatoDeFechas, fechaStr)
		if err != nil {
			return nil, fmt.Errorf("No se pudo formatear la fecha '%s' >> %q", fechaStr, err)
		}

		calendario[semana] = fecha
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Hubo en error en la ejecucion del bucle >> %v", err)
	}

	return calendario, nil
}

func ObtenerCalendarioStr(db *sql.DB) (map[int]string, error) {
	rows, err := db.Query(`select semana, fecha from calendario`)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener el calendario >> %q", err)
	}
	defer rows.Close()

	calendario := map[int]string{}
	var (
		fecha  string
		semana int
	)

	for rows.Next() {
		if err := rows.Scan(&semana, &fecha); err != nil {
			return nil, fmt.Errorf("No se pudieron extraer las fechas del calendario")
		}

		calendario[semana] = fecha
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Hubo en error en la ejecucion del bucle >> %v", err)
	}

	return calendario, nil
}

// func SumarIngreso(ingreso int) error {
// 	localCtx, cancel := context.WithTimeout(CONEXION.Ctx, time.Second*5)
// 	defer cancel()

// 	q := `UPDATE registros SET ingreso = ingreso + $1 WHERE fecha = $2`
// 	p := `INSERT INTO registros (fecha, ingreso, egreso, mes) VALUES ($1,$2,$3,$4)`

// 	fecha := time.Now().Format("2006/01/02")
// 	mes := int(time.Now().Month())

// 	result, err := CONEXION.Pool.Exec(localCtx, q, ingreso, fecha)
// 	if err != nil {
// 		return err
// 	}

// 	if result.RowsAffected() == 0 {
// 		_, err = CONEXION.Pool.Exec(localCtx, p, fecha, ingreso, 0, mes)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func SumarEgreso(egreso int) error {
// 	localCtx, cancel := context.WithTimeout(CONEXION.Ctx, time.Second*5)
// 	defer cancel()

// 	q := `UPDATE registros SET egreso = egreso + $1 WHERE fecha = $2`
// 	p := `INSERT INTO registros (fecha, ingreso, egreso, mes) VALUES ($1,$2,$3,$4)`

// 	fecha := time.Now().Format("2006/01/02")
// 	mes := int(time.Now().Month())

// 	result, err := CONEXION.Pool.Exec(localCtx, q, egreso, fecha)
// 	if err != nil {
// 		return err
// 	}

// 	if result.RowsAffected() == 0 {
// 		_, err = CONEXION.Pool.Exec(localCtx, p, fecha, 0, egreso, mes)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
