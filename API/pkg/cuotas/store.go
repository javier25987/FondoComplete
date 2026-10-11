package cuotas

import (
	"API-FONDO/global"
	"API-FONDO/pkg/logs"
	"fmt"
	"strconv"
	"time"
	// _ "time/tzdata"
)

func (s *centralStruct) storeObtenerDatosUser(id int) (*User, error) {
	user := &User{Id: id}

	// Obtener nommbre, puestos y telefono
	q := `SELECT nombre, puestos, telefono FROM usuarios WHERE id = ?`
	if err := s.Conexion.QueryRow(q, id).Scan(&user.Nombre, &user.Puestos, &user.Telefono); err != nil {
		return nil, fmt.Errorf("store >> Error al obtener datos de usuario : %v", err)
	}

	user.ValorCuota = s.Ajustes.Cuotas.ValorCuota
	user.ValorMulta = s.Ajustes.Cuotas.ValorMulta

	// Obtener tabla de cuotas
	tabla, cuotas_a_pagar, multas_a_pagar, err := s.ObtenerTablaCuotasYMultas(id)
	if err != nil {
		return nil, fmt.Errorf("store >> Error al obtener tabla de cuotas y multas : %v", err)
	}

	user.Tabla = tabla
	user.CuotasAPagar = cuotas_a_pagar
	user.MultasAPagar = multas_a_pagar

	q = `SELECT codigo, cuotas_pagas, multas_pagas, metodo_de_pago, fecha FROM cheques_de_cuotas WHERE id_user = ? ORDER BY codigo DESC`

	rows, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		codigo         int
		cuotas_pagas   int
		multas_pagas   int
		metodo_de_pago string
		fecha          time.Time
	)

	for rows.Next() {
		if err := rows.Scan(&codigo, &cuotas_pagas, &multas_pagas, &metodo_de_pago, &fecha); err != nil {
			return nil, err
		}

		user.Cheques = append(user.Cheques, &ChequeJson{
			Codigo:      codigo,
			CuotasPagas: cuotas_pagas,
			MultasPagas: multas_pagas,
			MetodoPago:  metodo_de_pago,
			FechaYHora:  fecha.Format("2006-01-02 15:04"),
		})

		user.ChequesUser = append(user.ChequesUser, codigo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error en el ciclo rows >> %v", err)
	}

	return user, nil
}

func (s *centralStruct) ObtenerTablaCuotasYMultas(id int) ([]*Semana, int, int, error) {
	// Obtenr Calandario
	calendario, err := global.ObtenerCalendarioStr(s.Conexion)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("No se puede obtener el calendario >> %q", err)
	}

	// Obtener Cuotas, Multas Y Bloqueos
	var (
		pagas  int
		deudas int
	)
	q := `SELECT u.cuotas_pagas, u.cuotas_deudas FROM usuarios u WHERE u.id = ?`
	if err := s.Conexion.QueryRow(q, id).Scan(&pagas, &deudas); err != nil {
		return nil, 0, 0, err
	}

	cuotas := make([]string, 50)
	multas_act := make([]string, 50)
	multas_pgs := make([]string, 50)

	// Hacer la columna de cuotas

	// obtener los bloqueos
	q = `SELECT cb.semana FROM cuotas_bloqueos cb WHERE cb.id_user = ?`
	rows1, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows1.Close()

	var (
		semana   int
		valor    int
		bloqueos int
	)

	for rows1.Next() {
		if err := rows1.Scan(&semana); err != nil {
			return nil, 0, 0, err
		}
		cuotas[semana] = "🔒 bloc"
	}

	if err := rows1.Err(); err != nil {
		return nil, -1, -1, err
	}

	cuotas_a_pagar := 50 - (pagas + bloqueos)

	idx := 0

	for pagas > 0 {
		if idx >= 50 {
			break
		}
		if cuotas[idx] == "" {
			cuotas[idx] = "✅ paga"
			pagas--
		}
		idx++
	}
	for deudas > 0 {
		if idx >= 50 {
			break
		}
		if cuotas[idx] == "" {
			cuotas[idx] = "🚨 deuda"
			deudas--
		}
		idx++
	}

	// Hacer la columna de cuotas pagadas

	var multas_a_pagar int

	q = `SELECT cma.semana, cma.valor FROM cuotas_multas_activas cma WHERE cma.id_user = ?`
	rows2, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows2.Close()

	for rows2.Next() {
		if err := rows2.Scan(&semana, &valor); err != nil {
			return nil, 0, 0, err
		}
		multas_act[semana] = strconv.Itoa(valor)
		multas_a_pagar += valor
	}

	if err := rows2.Err(); err != nil {
		return nil, -1, -1, err
	}

	q = `SELECT cmp.semana, cmp.valor FROM cuotas_multas_pagas cmp WHERE cmp.id_user = ?`
	rows3, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows3.Close()

	for rows3.Next() {
		if err := rows3.Scan(&semana, &valor); err != nil {
			return nil, 0, 0, err
		}
		multas_pgs[semana] = strconv.Itoa(valor)
	}

	if err := rows3.Err(); err != nil {
		return nil, -1, -1, err
	}

	salida := make([]*Semana, 50)

	for i := range 50 {
		salida[i] = &Semana{
			Semana:      i + 1,
			Fecha:       calendario[i+1][:10],
			Cuota:       cuotas[i],
			MultaActiva: multas_act[i],
			MultaPaga:   multas_pgs[i],
		}
	}

	return salida, cuotas_a_pagar, multas_a_pagar, nil
}

func (s *centralStruct) storePagarSemanas(id int, formato *FormatoPago) (int, error) {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return -1, fmt.Errorf("Pago de cuotas: no se pudo crear la transacción >> %v", err)
	}
	defer tx.Rollback()

	var puestos int
	q := `SELECT puestos FROM usuarios WHERE id = ?`
	if err := tx.QueryRow(q, id).Scan(&puestos); err != nil {
		return -1, fmt.Errorf("store >> Pago de cuotas: no se pudo obtener puestos >> %v", err)
	}

	var (
		total_cuotas int
		total_multas int
	)

	if formato.Cuotas >= 0 { // pagar cuotas
		q = `UPDATE usuarios SET cuotas_pagas = cuotas_pagas + ? WHERE id = ?`
		if _, err = tx.Exec(q, formato.Cuotas, id); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo actualizar cuotas pagas >> %v", err)
		}

		q = `SELECT cuotas_deudas FROM usuarios WHERE id = ?`
		var deudas int
		if err = tx.QueryRow(q, id).Scan(&deudas); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo obtener deudas >> %v", err)
		}

		if deudas > 0 {
			if formato.Cuotas >= deudas {
				deudas = 0
			} else {
				deudas -= formato.Cuotas
			}

			q = `UPDATE usuarios SET cuotas_deudas = ? WHERE id = ?`
			if _, err = tx.Exec(q, deudas, id); err != nil {
				return -1, fmt.Errorf("Pago de cuotas: no se pudo actualizar deudas >> %v", err)
			}
		}

		total_cuotas = formato.Cuotas * s.Ajustes.Cuotas.ValorCuota * puestos

		q = `UPDATE usuarios SET capital_pago = capital_pago + ? WHERE id = ?`
		if _, err := tx.Exec(q, total_cuotas, id); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo actualizar capital pago >> %v", err)
		}
	}

	if formato.Multas >= 0 { // pago de multas
		total_multas = formato.Multas * s.Ajustes.Cuotas.ValorMulta * puestos

		q = `SELECT semana, valor FROM cuotas_multas_activas WHERE id_user = ? ORDER BY semana`

		rows1, err := tx.Query(q, id)
		if err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo obtener multas activas >> %v", err)
		}
		defer rows1.Close()

		var (
			semana int
			valor  int

			semanas_a_pagar = formato.Multas

			nuevo_mutas_act = map[int]int{} // nuevos valores para multas activas
			set_multas_pgs  = map[int]int{} // valores a sumar a una semana en multas pagas
		)

		// Restamos todas las multas que pagamos

		for rows1.Next() {
			if err := rows1.Scan(&semana, &valor); err != nil {
				return -1, fmt.Errorf("Pago de cuotas: no se pudo escanear multas activas >> %v", err)
			}

			if semanas_a_pagar > valor {
				nuevo_mutas_act[semana] = 0
				set_multas_pgs[semana] = valor
				semanas_a_pagar -= valor
			} else {
				nuevo_mutas_act[semana] = valor - semanas_a_pagar
				set_multas_pgs[semana] = semanas_a_pagar
				break
			}
		}

		if err := rows1.Err(); err != nil {
			return -1, err
		}

		// Escribimos en multas activas las nuevas multas del usuario

		q = `UPDATE cuotas_multas_activas SET valor = ? WHERE id_user = ? AND semana = ?`

		for llave, valor := range nuevo_mutas_act {
			if _, err := tx.Exec(q, valor, id, llave); err != nil {
				return -1, fmt.Errorf("Pago de cuotas: no se pudo actualizar multas activas >> %v", err)
			} //
		}

		// Escribimos en multas pagas las multas que el usuario pago

		q = `
			insert into cuotas_multas_pagas(id_user, semana, valor) 
			values (?, ?, ?)
			on conflict(id_user, semana)
			do update set valor = valor + excluded.valor;
		`

		for llave, valor := range set_multas_pgs {
			_, err := tx.Exec(q, id, llave, valor)
			if err != nil {
				return -1, fmt.Errorf("Pago de cuotas: no se pudo actualizar multas pagas >> %v", err)
			}
		}

		// Sumamos al pago de multas semanal

		q = `UPDATE usuarios SET multas_semanales = multas_semanales + ? WHERE id = ?`
		if _, err := tx.Exec(q, total_multas, id); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo actualizar multas semanales >> %v", err)
		}

		if _, err := tx.Exec("DELETE FROM cuotas_multas_activas WHERE valor = 0"); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo limpiar multas activas >> %v", err)
		}

		if _, err := tx.Exec("DELETE FROM cuotas_multas_pagas WHERE valor = 0"); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo limpiar multas pagas >> %v", err)
		}
	}

	// Crear el cheque
	q = `INSERT INTO cheques_de_cuotas (id_user,cuotas_pagas,multas_pagas,metodo_de_pago) VALUES (?,?,?,?)`

	result, err := tx.Exec(q, id, formato.Cuotas, formato.Multas, formato.Metodo)
	if err != nil {
		return -1, fmt.Errorf("Pago de cuotas: no se pudo crear el cheque >> %v", err)
	}

	// Si es transferencia registrarla
	if formato.Metodo == "TRANSFERENCIA" {
		q = `INSERT INTO transferencias (id_user, monto) VALUES (?, ?)`

		if _, err = tx.Exec(q, id, total_cuotas+total_multas); err != nil {
			return -1, fmt.Errorf("Pago de cuotas: no se pudo registrar la transferencia >> %v", err)
		}
	}

	// Commit de la transacción
	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("Pago de cuotas: no se pudo commitear la transacción >> %v", err)
	}

	codigo, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Pago de cuotas: no se pudo obbtenr el codiigo del cheque >> %q", err)
	}

	// guardar el registro
	// global.SumarIngreso(total_cuotas + total_multas)

	return int(codigo), nil
}

func (s *centralStruct) storeRectificarAsuntos(id int) error {

	calendario, err := global.ObtenerCalendarioTime(s.Conexion)
	if err != nil {
		return fmt.Errorf("Rectificar asuntos: no se pudo obtener el calendario >> %v", err)
	}

	// creamos la transaccion
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("Rectificar asuntos: no se pudo crear la transacción >> %v", err)
	}
	defer tx.Rollback()

	// obbtenemos las semanas revisadas
	var semanas_revisadas int
	q := `SELECT semanas_revisadas FROM usuarios WHERE id = ?`
	if err := tx.QueryRow(q, id).Scan(&semanas_revisadas); err != nil {
		return fmt.Errorf("Rectificar asuntos: no se pudo obtener semanas revisadas >> %v", err)
	}

	if semanas_revisadas >= 50 { // no hacer nada si ya todas las semanas fueron revisadas
		logMsg := fmt.Sprintf("No fue necesario revisar los asuntos del usuario %d", id)
		logs.RegistrarLog(logMsg)
		return nil
	}

	// loc, err := time.LoadLocation("America/Bogota")
	// if err != nil {
	// 	return fmt.Errorf("store >> Error cargando zona horaria: %v", err)
	// }

	fecha_actual := time.Now() //.In(loc)

	var semanas_a_revisar int

	for i := range 50 {
		if calendario[i+1].After(fecha_actual) {
			break
		}
		semanas_a_revisar++
	}

	if semanas_a_revisar > semanas_revisadas {
		logMsg := fmt.Sprintf( // creamos el log de la operacion
			"Fue necesario arreglar los asuntos del usuario %d por %d semana(s), la fecha de comparacion es |%s|",
			id,
			semanas_a_revisar-semanas_revisadas,
			fecha_actual.Format("2006-01-02 15:04:05"),
		)
		logs.RegistrarLog(logMsg)

		semanas_array := [50]int{}

		/*
			1 pagas
			0 deudas
			-1 bloqueos
		*/

		q = `SELECT semana FROM cuotas_bloqueos WHERE id_user = ?`
		rows1, err := tx.Query(q, id)
		if err != nil {
			return fmt.Errorf("Rectificar asuntos: no se pudo obtener bloqueos >> %v", err)
		}
		defer rows1.Close()

		var idx int
		for rows1.Next() {
			if err := rows1.Scan(&idx); err != nil {
				return err
			}
			semanas_array[idx] = -1
		}

		if err := rows1.Err(); err != nil {
			return err
		}

		var semanas_pagas int
		q = `SELECT cuotas_pagas FROM usuarios WHERE id = ?`
		if err := tx.QueryRow(q, id).Scan(&semanas_pagas); err != nil {
			return fmt.Errorf("Rectificar asuntos: no se pudo obtener cuotas pagas >> %v", err)
		}

		idx = 0

		for semanas_pagas > 0 {
			if semanas_array[idx] == 0 {
				semanas_array[idx] = 1
				semanas_pagas--
			}
			idx++
		}

		deudas := 0

		for i := semanas_a_revisar - 1; i >= 0; i-- {
			if semanas_array[i] == 0 {
				deudas++
			}
		}

		gen := newGenNum(semanas_a_revisar - semanas_revisadas)
		add_multas := map[int]int{}

		for i := semanas_a_revisar - 1; i >= 0; i-- {
			if semanas_array[i] == 0 && s.Ajustes.Cuotas.CobrarMultas {
				add_multas[i] = gen.Next()
			}
		}

		// guardamos deudas revisiones y multas

		q = `UPDATE usuarios SET cuotas_deudas = ?, semanas_revisadas = ? WHERE id = ?`
		if _, err := tx.Exec(q, deudas, semanas_a_revisar, id); err != nil {
			return fmt.Errorf("Rectificar asuntos: no se pudo actualizar usuario >> %v", err)
		}

		q = `
			insert into cuotas_multas_activas(id_user, semana, valor) 
			values (?, ?, ?)
			on conflict(id_user, semana)
			do update set valor = valor + excluded.valor;
		`
		stmt, err := tx.Prepare(q)
		if err != nil {
			return fmt.Errorf("Rectificar asuntos: no se pudo precompilar la query >> %v", err)
		}
		defer stmt.Close()

		for llave, valor := range add_multas {
			_, err := stmt.Exec(id, llave, valor)
			if err != nil {
				return fmt.Errorf("Rectificar asuntos: no se pudo actualizar multas activas >> %v", err)
			}
		}

		// limpiamos todos los valores nulos
		q = `DELETE FROM cuotas_multas_activas WHERE id_user = ? AND valor = 0`
		if _, err := tx.Exec(q, id); err != nil {
			return fmt.Errorf("Rectificar asuntos: no se pudo limpiar multas activas >> %v", err)
		}

		q = `DELETE FROM cuotas_multas_pagas WHERE id_user = ? AND valor = 0`
		if _, err := tx.Exec(q, id); err != nil {
			return fmt.Errorf("Rectificar asuntos: no se pudo limpiar multas pagas >> %v", err)
		}
	} else {
		logMsg := fmt.Sprintf("La revision del usuario %d esta al dia", id)
		logs.RegistrarLog(logMsg)
	}

	return tx.Commit()
}

func (s *centralStruct) ObtenerDatosChequeCuotas(codigo int) (*ChequePdf, error) {
	cheque := &ChequePdf{Codigo: codigo}

	var fecha time.Time

	q := `SELECT id_user,cuotas_pagas, multas_pagas, metodo_de_pago, fecha FROM cheques_de_cuotas WHERE codigo = ?`

	if err := s.Conexion.QueryRow(q, codigo).Scan(&cheque.Id, &cheque.CuotasPagas, &cheque.MultasPagas, &cheque.MetodoPago, &fecha); err != nil {
		return nil, fmt.Errorf("store >> No se puso obtener los datos del usuario >> %v", err)
	}

	cheque.FechaYHora = fecha.Format("2006-01-02 15:04")

	cheque.ValorCuota = s.Ajustes.Cuotas.ValorCuota
	cheque.ValorMulta = s.Ajustes.Cuotas.ValorMulta

	q = `SELECT nombre, puestos FROM usuarios WHERE id = ?`
	if err := s.Conexion.QueryRow(q, cheque.Id).Scan(&cheque.Nombre, &cheque.Puestos); err != nil {
		return nil, fmt.Errorf("store >> error al obtener datos usuario : %v", err)
	}

	cheque.TotalCuotas = cheque.CuotasPagas * cheque.ValorCuota * cheque.Puestos
	cheque.TotalMultas = cheque.MultasPagas * cheque.ValorMulta * cheque.Puestos
	cheque.TotalPagado = cheque.TotalCuotas + cheque.TotalMultas

	return cheque, nil
}
