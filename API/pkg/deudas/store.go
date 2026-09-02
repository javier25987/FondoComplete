package deudas

import (
	"fmt"
	"time"
)

func (s *centralStruct) storeCrearDeuda(id int, deuda *PostDeuda) error {
	fecha := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	q := `INSERT INTO deudas (monto, motivo, tipo, id_user, fecha_de_pago) VALUES (?,?,?,?,?)`

	if _, err := s.Conexion.Exec(q, deuda.Monto, deuda.Motivo, deuda.Tipo, id, fecha); err != nil {
		return fmt.Errorf("No se pudo crear la deuda >> %v", err)
	}

	return nil
}

func (s *centralStruct) storeObtenerPropietario(codigo int) (int, error) {
	q := `SELECT id_user FROM deudas WHERE codigo = ?`

	var id_user int
	if err := s.Conexion.QueryRow(q, codigo).Scan(&id_user); err != nil {
		return -1, fmt.Errorf("No se pudo obtener el propietario de la deuda %d >> %v", codigo, err)
	}

	return id_user, nil
}

func (s *centralStruct) storeObtenerEstado(codigo int) (bool, error) {
	q := `SELECT estado FROM deudas WHERE codigo = ?`

	var estado bool
	if err := s.Conexion.QueryRow(q, codigo).Scan(&estado); err != nil {
		return false, fmt.Errorf("No se pudo obtener el estado de la deuda %d >> %v", codigo, err)
	}

	return estado, nil
}

func (s *centralStruct) storePagarDeuda(codigo int) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("No de pudo iniciar la transaccion >> %v", err)
	}
	defer tx.Rollback()

	fecha := time.Now()

	q := `UPDATE deudas SET fecha_de_pago = ?, estado = FALSE WHERE codigo = ?`

	if _, err := tx.Exec(q, fecha, codigo); err != nil {
		return fmt.Errorf("No se pudo pagar la deuda >> %v", err)
	}

	var (
		tipo    string
		monto   int
		id_user int
	)

	q = `SELECT tipo, monto, id_user FROM deudas WHERE codigo = ?`
	if err := tx.QueryRow(q, codigo).Scan(&tipo, &monto, &id_user); err != nil {
		return fmt.Errorf("No se pudo obtener la infomacion de la deuda")
	}

	switch tipo {
	case "MULTA POR CABEZA":
		q = `UPDATE usuarios SET multas_por_cabeza = multas_por_cabeza + ? WHERE id = ?`
	case "MULTA EXTRA":
		q = `UPDATE usuarios SET multas_extras = multas_extras + ? WHERE id = ?`
	default:
		return fmt.Errorf("Tipo de deuda no reconocido")
	}

	if _, err := tx.Exec(q, monto, id_user); err != nil {
		return fmt.Errorf("No se pudo sumar el monto de la deuda a las multas >> %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("No se pudo cerrar la transaccion >> %v", err)
	}

	return nil
}

func (s *centralStruct) storeObtenerDeudas(id int) (*GetDeudas, error) {
	q := `SELECT codigo, monto, fecha_de_creacion, fecha_de_pago, tipo, motivo, estado FROM deudas WHERE id_user = ?`

	user := &GetDeudas{Id: id}

	rows, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("No se pudo inicial la consulta de los datos >> %v", err)
	}
	defer rows.Close()

	var (
		fecha_de_creacion time.Time
		fecha_de_pago     time.Time
	)

	for rows.Next() {
		deuda := &Deuda{}

		if err := rows.Scan(&deuda.Codigo, &deuda.Monto, &fecha_de_creacion, &fecha_de_pago, &deuda.Tipo, &deuda.Motivo, &deuda.Estado); err != nil {
			return nil, fmt.Errorf("No se pudo obtener los datos de la deuda >> %v", err)
		}

		deuda.FechaCreacion = fecha_de_creacion.Format("2006/01/02 15:04")
		deuda.FechaPago = fecha_de_pago.Format("2006/01/02 15:04")

		user.Deudas = append(user.Deudas, deuda)
	}

	q = `SELECT nombre FROM usuarios WHERE id = ?`
	if err := s.Conexion.QueryRow(q, id).Scan(&user.Nombre); err != nil {
		return nil, fmt.Errorf("No se pudo obtener el nombre del usuario >> %v", err)
	}

	return user, nil
}
