package prestamos

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (s *centralStruct) storeObtenerPropietario(codigo int) (int, error) {
	q := `SELECT id_user FROM prestamos_hechos WHERE codigo = ?`
	var id_user int
	if err := s.Conexion.QueryRow(q, codigo).Scan(&id_user); err != nil {
		return -1, fmt.Errorf("No se pudo obtener el propietario del prestamo %d >> %v", codigo, err)
	}

	return id_user, nil
}

func (s *centralStruct) storeObtenerDeudaEIntereses(codigo int) (int, int, error) {
	q := `SELECT deuda, intereses_vencidos FROM prestamos_hechos WHERE codigo = ?`
	var deuda, intereses int
	if err := s.Conexion.QueryRow(q, codigo).Scan(&deuda, &intereses); err != nil {
		return -1, -1, fmt.Errorf("No se pudo obtener la deuda del prestamo %d >> %v", codigo, err)
	}

	return deuda, intereses, nil
}

func (s *centralStruct) storeObtenerEstadoDelPresamo(codigo int) (bool, error) {
	q := `select ph.estado from prestamos_hechos ph where ph.codigo = ?`
	var estado bool
	if err := s.Conexion.QueryRow(q, codigo).Scan(&estado); err != nil {
		return false, fmt.Errorf("No se pudo obtener el estado del prestamo %d >> %v", codigo, err)
	}
	return estado, nil
}

func (s *centralStruct) storeConsultarCapitalDisponibleUsr(id int) (int, error) {
	if id == 1976 { // esto es para el super usuario
		return 1_000_000_000, nil
	}

	q := `
	SELECT u.capital_pago * ?/100.0 - (
			SELECT COALESCE(SUM(ph.deuda + ph.intereses_vencidos), 0)
			FROM prestamos_hechos ph
			WHERE ph.id_user = ?
	) - (
			SELECT COALESCE(SUM(p.monto_respaldado), 0)
			FROM prestamos_fiadores p
			WHERE p.id_user_fiador = ?
	)
	FROM usuarios u 
	WHERE u.id = ?;
	`
	var capital float64

	if err := s.Conexion.QueryRow(q, s.Ajustes.Prestamos.CapitalUsable, id, id, id).Scan(&capital); err != nil {
		return 0, fmt.Errorf("No se pudo consultar el capital disponible de %d: %v", id, err)
	}

	return int(capital), nil
}

/*
+ obtener intereses e interes generado
+ aumentar capital generado
+ insertar el presamo
+ crear las fechas
+ cargar fiadores
+ crear nota del prestamo
*/
func (s *centralStruct) storeCargarPrestamo(id int, c *FormatoCrearPrestamo) (int64, error) {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	interes := s.Ajustes.Intereses.InteresmTope // definimos el interes del prestamo
	if c.Monto >= s.Ajustes.Intereses.TopeIntereses {
		interes = s.Ajustes.Intereses.InteresMTope
	}
	interes_inicial := int(c.Monto * interes / 100)

	sum_deudas := 0 // rectificamos el capital retirado por el usuario
	for _, fiador := range c.Fiadores {
		sum_deudas += fiador.Deuda
	}
	capital_retirado := c.Monto - sum_deudas

	q := `UPDATE usuarios SET capital_retirado = capital_retirado + ? WHERE id = ?`
	if _, err := tx.Exec(q, capital_retirado, id); err != nil {
		return -1, fmt.Errorf("No se puede cargar el capital retirado: %v", err)
	}

	// escribimos el prestamo en la base de datos
	q = `
		INSERT INTO prestamos_hechos(id_user, interes, intereses_generados, deuda, monto, motivo) 
		VALUES (?,?,?,?,?,?)
	`

	result, err := tx.Exec(q, id, interes, interes_inicial, c.Monto, c.Monto, c.Motivo)
	if err != nil {
		return -1, fmt.Errorf("No se pudo escribir el prestamo >> %v", err)
	}
	codigo, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("No se pudo obtener el codigo del nuevo prestamo >> %q", err)
	}

	q = `INSERT INTO prestamos_fechas_de_pago (codigo_prestamo, fecha) VALUES (?,?)`
	stmt, err := tx.Prepare(q)
	if err != nil {
		return -1, fmt.Errorf("No se pudo precompilar la consulta >> %q", err)
	}
	defer stmt.Close()

	for _, fecha := range CalendarioMeses(time.Now(), s.Ajustes.Prestamos.FechaDeCierre) {
		if _, err := stmt.Exec(codigo, fecha); err != nil {
			return -1, fmt.Errorf("No se pudo escribir las fechas de pago del prestamo >> %v", err)
		}
	}

	q = `INSERT INTO prestamos_fiadores (id_user_fiador, codigo_prestamo, monto_respaldado) VALUES (?,?,?)`

	p := `INSERT INTO super_usuario (codigo_prestamo, monto_respaldado) VALUES (?,?)`

	for _, fiador := range c.Fiadores {
		if fiador.Id_user == 1976 {
			if _, err := tx.Exec(p, codigo, fiador.Deuda); err != nil {
				return -1, fmt.Errorf("No se pudo guardar al super usuario >> %v", err)
			}
		} else {
			if _, err := tx.Exec(q, fiador.Id_user, codigo, fiador.Deuda); err != nil {
				return -1, fmt.Errorf("No se pudo guardar a los fiaadores >> %v", err)
			}
		}
	}

	var cadena strings.Builder
	for _, fiador := range c.Fiadores {
		cadena.WriteString(fmt.Sprintf(" %d:%d", fiador.Id_user, fiador.Deuda))
	}

	msg := "Prestamo creado "

	if fiadoresStr := cadena.String(); fiadoresStr != "" {
		msg += fmt.Sprintf("fiadores;%s", cadena.String())
	}

	q = `INSERT INTO notas_de_prestamos (codigo_prestamo, anotacion) VALUES (?,?)`

	if _, err := tx.Exec(q, codigo, msg); err != nil {
		return -1, fmt.Errorf("No se pudo crear la nota de prestamo >> %v", err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("No se pudo guardar la operacion >> %q", err)
	}

	return codigo, nil
}

func (s *centralStruct) storeAbonarPrestamo(id int, a *AbonarPrestamo) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	fechaActual := time.Now()

	if a.Fecha != "" {
		fechaActual, err = time.Parse("2006-01-02", a.Fecha)
		if err != nil {
			return fmt.Errorf("La fecha ingresada `%s` para remplazar la fecha actual no puede ser aceptada por el programa >> %v", a.Fecha, err)
		}
	}

	// hacer la nota del prestamo
	q := `INSERT INTO notas_de_prestamos (codigo_prestamo, anotacion, fecha) VALUES (?,?,?)`
	msg := fmt.Sprintf("Se abono %d a %s", a.Monto, a.Tipo)
	if _, err := tx.Exec(q, a.Codigo, msg, fechaActual); err != nil {
		return fmt.Errorf("Pagar prestamo: %d no se pudo crear la anotacion del prestamo >> %v", a.Codigo, err)
	}

	// obtener los datos del prestamo (deuda, intereses)
	q = `SELECT deuda, intereses_vencidos FROM prestamos_hechos WHERE codigo = ?`
	var deuda, intereses int
	if err := tx.QueryRow(q, a.Codigo).Scan(&deuda, &intereses); err != nil {
		return fmt.Errorf("Pagar prestamo: %d no se pudo obtener deuda e intereses >> %v", a.Codigo, err)
	}

	var (
		deudaAntes       = deuda
		deudaDespues     int
		interesesAntes   = intereses
		interesesDespues int
	)

	switch a.Tipo {
	case "DEUDA":
		deudaDespues = deudaAntes - a.Monto
		interesesDespues = interesesAntes
	case "INTERESES":
		interesesDespues = interesesAntes - a.Monto
		deudaDespues = deudaAntes
	default:
		return errors.New("Como si quiera es posible que hayas llegado hasta este punto?")
	}

	q = `INSERT INTO cheques_de_prestamos (id_user, codigo_prestamo, monto_pago, tipo, deuda_antes, deuda_despues, intereses_antes, intereses_despues, fecha) VALUES (?,?,?,?,?,?,?,?,?)`
	if _, err := tx.Exec(q, id, a.Codigo, a.Monto, a.Tipo, deudaAntes, deudaDespues, interesesAntes, interesesDespues, fechaActual); err != nil {
		return fmt.Errorf("Pagar prestamo: %d no se pudo crear el cheque >> %v", a.Codigo, err)
	}

	// rectificar si el prestamo queda totalmente pago
	if deudaDespues+interesesDespues == 0 {
		q = `UPDATE prestamos_hechos SET estado = FALSE WHERE codigo = ?`
		if _, err := tx.Exec(q, a.Codigo); err != nil {
			return fmt.Errorf("Pagar prestamo: %d no se pudo actualizar estado del prestamo >> %v", a.Codigo, err)
		}
	}

	/*
		Ya antes hubo un switch para definir los valores de deuda e intereses antes y despues pero no
		quiero mesclar logica para ahorrarme un switch, la parte anterior es la logica necesaria para
		hacer el chque y rectifiicar sie l restamoe sta pago, esta parte se encarga de realizar todos
		los cambios pertienetes en la base de datos para hacer el pago
	*/
	switch a.Tipo {
	case "DEUDA":
		// pagar a fiadores
		q = `SELECT id_user_fiador, monto_respaldado FROM prestamos_fiadores WHERE codigo_prestamo = ?`
		rows, err := tx.Query(q, a.Codigo)
		if err != nil {
			return fmt.Errorf("Pagar prestamo: %d no se pudo obtener los fiadores >> %v", a.Codigo, err)
		}
		defer rows.Close()

		pago_fiadores := a.Monto
		var fiador, monto_respaldado int
		cargar_fiadores := map[int]int{}

		for rows.Next() {
			if err := rows.Scan(&fiador, &monto_respaldado); err != nil {
				return fmt.Errorf("Pagar prestamo: %d no se pudo obtener fiadores y montos >> %v", a.Codigo, err)
			}

			if monto_respaldado >= pago_fiadores {
				cargar_fiadores[fiador] = monto_respaldado - pago_fiadores
				break
			} else {
				cargar_fiadores[fiador] = 0
				pago_fiadores -= monto_respaldado
			}
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("Hubo un error en la lectura de los fiadores >> %v", err)
		}

		q = `UPDATE prestamos_fiadores SET monto_respaldado = ? WHERE id_user_fiador = ? AND codigo_prestamo = ?`
		for fiador, monto := range cargar_fiadores {
			if _, err := tx.Exec(q, monto, fiador, a.Codigo); err != nil {
				return fmt.Errorf("Pagar prestamo: %d no se pudo cargar monto a fiadores >> %v", a.Codigo, err)
			}
		}

		// descontar de deuda
		q = `UPDATE prestamos_hechos SET deuda = deuda - ? WHERE codigo = ?`
		if _, err := tx.Exec(q, a.Monto, a.Codigo); err != nil {
			return fmt.Errorf("Pagar prestamo %d: no se pudo descontar la deuda >> %v", a.Codigo, err)
		}

	case "INTERESES":
		// pagar los intereses
		q = `
			UPDATE prestamos_hechos 
			SET intereses_vencidos = intereses_vencidos - ?, intereses_generados = intereses_generados + ?
			WHERE codigo = ?
		`
		if _, err := tx.Exec(q, a.Monto, a.Monto, a.Codigo); err != nil {
			return fmt.Errorf("Pagar prestamo: %d no se pudo descontar los intereses >> %v", a.Codigo, err)
		}
	}

	return tx.Commit()
}

func (s *centralStruct) storeArreglarAsuntos(id int) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("No se pudo crear la transacion para escribir los datos >> %v", err)
	}
	defer tx.Rollback()

	// obtenemos los datos del usuario
	q := `SELECT codigo, fechas_revisadas FROM prestamos_hechos WHERE id_user = ? AND estado = TRUE`
	rows1, err := tx.Query(q, id)
	if err != nil {
		return fmt.Errorf("No se pudieron obtener los prestamos activos del usuario %d >> %v", id, err)
	}
	defer rows1.Close()

	var (
		codigo           int
		fechas_revisadas int
	)

	prestamos := map[int]int{}

	for rows1.Next() {
		if err := rows1.Scan(&codigo, &fechas_revisadas); err != nil {
			return fmt.Errorf("No se pudieron escanear los prestamos activos del usuario %d >> %v", id, err)
		}
		prestamos[codigo] = fechas_revisadas
	}

	if err := rows1.Err(); err != nil {
		return fmt.Errorf("Hubo un error en el escaneo de las fechas revisadas >> %v", err)
	}

	q = `SELECT COUNT(*) FROM prestamos_fechas_de_pago WHERE codigo_prestamo = ? AND fecha < CURRENT_DATE`

	p := `
		UPDATE prestamos_hechos 
		SET 
			intereses_vencidos = intereses_vencidos + (deuda * interes * ?)/100, 
			fechas_revisadas = ? 
		WHERE codigo = ?
	`

	t := `INSERT INTO notas_de_prestamos (codigo_prestamo, anotacion) VALUES (?,?)`

	v := `select deuda from prestamos_hechos where codigo = ?`

	var (
		msg           string
		deudaPrestamo int
	)

	var pasadas, diff int

	for codigo, revisiones := range prestamos {
		if err := tx.QueryRow(q, codigo).Scan(&pasadas); err != nil {
			return fmt.Errorf("No se pudo contar las fechas pasadas del prestamo %d >> %v", codigo, err)
		}

		if pasadas > revisiones {
			diff = pasadas - revisiones
			if _, err := tx.Exec(p, diff, pasadas, codigo); err != nil {
				return fmt.Errorf("No se pudo cargar intereses vencidos al prestamo %d >> %v", codigo, err)
			}

			if err := tx.QueryRow(v, codigo).Scan(&deudaPrestamo); err != nil {
				return fmt.Errorf("No se pudo consultar la deuda del prestamo %d al cual hay que cargarle intereses >> %v", codigo, err)
			}

			msg = fmt.Sprintf("Se cargaron intereses por %d fechas vencidas sobre una deuda de %d", diff, deudaPrestamo)

			if _, err := tx.Exec(t, codigo, msg); err != nil {
				return fmt.Errorf("No se pudo crear el registro de la creacion de los intereses >> %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("No se pudo cerrar la transaccion >> %v", err)
	}

	return nil
}

func (s *centralStruct) storeGetDataUser(id int) (*GetUsuario, error) {
	user := GetUsuario{Id: id}

	q := `SELECT nombre FROM usuarios WHERE id = ?`
	if err := s.Conexion.QueryRow(q, id).Scan(&user.Nombre); err != nil {
		return nil, fmt.Errorf("No se pudo consultar el nombre del usuario >> %v", err)
	}

	var (
		codigo           int
		interes          int
		interes_vencido  int
		interes_generado int
		deuda            int
		monto            int
		deuda_actual     int

		id_user_fiador int
		deuda_fiador   int

		fecha_de_pago time.Time

		estado     bool
		estado_str string
		motivo     string

		fecha_anotacion time.Time
		anotacion       string
	)

	q = `SELECT codigo, interes, intereses_vencidos, intereses_generados, deuda, monto, estado, motivo
	FROM prestamos_hechos WHERE id_user = ? ORDER BY codigo DESC`

	rows1, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener los datos del usuario >> %v", err)
	}
	defer rows1.Close()

	deuda_prestamo := map[int]int{}
	total_deuda_prestamos := 0

	for rows1.Next() {
		if err := rows1.Scan(&codigo, &interes, &interes_vencido, &interes_generado, &deuda, &monto, &estado, &motivo); err != nil {
			return nil, fmt.Errorf("No se pudo leer los datos del usuario >> %v", err)
		}
		deuda_actual = deuda + interes_vencido

		if deuda_actual > 0 {
			deuda_prestamo[codigo] = deuda_actual
			total_deuda_prestamos += deuda_actual
		}

		user.PrestamosUser = append(user.PrestamosUser, codigo)

		estado_str = "PAGO ✅"
		if estado {
			estado_str = "ACTIVO ⏳"
		}

		user.Prestamos = append(user.Prestamos, &Prestamo{
			Codigo:          codigo,
			Interes:         interes,
			InteresVencido:  interes_vencido,
			InteresGenerado: interes_generado,
			Deuda:           deuda,
			Monto:           monto,
			DeudaActual:     deuda_actual,
			Estado:          estado_str,
			Motivo:          motivo,
		})
	}

	if err := rows1.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error en la obtencion de los datos del usuario >> %v", err)
	}

	// cargar fiadores, fechas de pago y anotaciones
	q = `SELECT id_user_fiador, monto_respaldado FROM prestamos_fiadores WHERE codigo_prestamo = ?;`
	p := `SELECT fecha FROM prestamos_fechas_de_pago WHERE codigo_prestamo = ? ORDER BY fecha ASC`
	k := `SELECT anotacion, fecha FROM notas_de_prestamos WHERE codigo_prestamo = ? ORDER BY fecha ASC`

	var nota string

	for _, prestamo := range user.Prestamos {
		// fiadores
		rows2, err := s.Conexion.Query(q, prestamo.Codigo)
		if err != nil {
			return nil, fmt.Errorf("No se pudieron obtener los fiadores del prestamo %d >> %v", prestamo.Codigo, err)
		}
		defer rows2.Close()

		for rows2.Next() {
			if err := rows2.Scan(&id_user_fiador, &deuda_fiador); err != nil {
				return nil, fmt.Errorf("No se pudieron obtener los fiadores del prestamo %d >> %v", prestamo.Codigo, err)
			}
			prestamo.Fiadores = append(prestamo.Fiadores, &Fiador{Id_user: id_user_fiador, Deuda: deuda_fiador})
		}

		if err := rows2.Err(); err != nil {
			return nil, fmt.Errorf("Hubo un error en el bucle de lectura rows2 >> %v", err)
		}

		// fechas de pago
		rows3, err := s.Conexion.Query(p, prestamo.Codigo)
		if err != nil {
			return nil, fmt.Errorf("No se pudieron obtener las fechas de pago del prestamo %d >> %v", prestamo.Codigo, err)
		}
		defer rows3.Close()

		for rows3.Next() {
			if err := rows3.Scan(&fecha_de_pago); err != nil {
				return nil, fmt.Errorf("No se pudo obtener la fecha del prestamo %d >> %v", prestamo.Codigo, err)
			}
			prestamo.FechasPago = append(prestamo.FechasPago, fecha_de_pago.Format("2006/01/02"))
		}

		if err := rows3.Err(); err != nil {
			return nil, fmt.Errorf("Hubo un error en el bucle de lectura rows3 >> %v", err)
		}

		// anotaciones
		rows4, err := s.Conexion.Query(k, prestamo.Codigo)
		if err != nil {
			return nil, fmt.Errorf("No se pudieron obtener las anotaciones del prestamo %d >> %v", prestamo.Codigo, err)
		}
		defer rows4.Close()

		for rows4.Next() {
			if err := rows4.Scan(&anotacion, &fecha_anotacion); err != nil {
				return nil, fmt.Errorf("No se pudo obtener las anotaciones del prestamo %d >> %v", prestamo.Codigo, err)
			}
			nota = fmt.Sprintf("[%s] %s", fecha_anotacion.Format("2006/01/02"), anotacion)
			prestamo.Anotaciones = append(prestamo.Anotaciones, nota)
		}

		if err := rows4.Err(); err != nil {
			return nil, fmt.Errorf("Hubo un error en el bucle de lectura rows4 >> %v", err)
		}
	}

	consulta_capital := &ConsultaCapital{TotalDeudasEnPrestamos: total_deuda_prestamos}
	for codigo, deuda := range deuda_prestamo {
		consulta_capital.DeudasEnPrestamos = append(consulta_capital.DeudasEnPrestamos, &DeudaPrestamo{
			Prestamo: codigo,
			Deuda:    deuda,
		})
	}

	q = `SELECT codigo_prestamo, monto_respaldado FROM prestamos_fiadores WHERE id_user_fiador = ? AND monto_respaldado != 0`
	rows5, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener las deudas como fiador del usuario %d >> %v", id, err)
	}
	defer rows5.Close()

	total_deudas_fiador := 0
	for rows5.Next() {
		if err := rows5.Scan(&codigo, &deuda); err != nil {
			return nil, fmt.Errorf("No se pudieron escanear las deudas como fiador del usuario %d >> %v", id, err)
		}
		consulta_capital.DeudasPorFiador = append(consulta_capital.DeudasPorFiador, &DeudaFiador{
			Cuanto:     deuda,
			EnPrestamo: codigo,
		})
		total_deudas_fiador += deuda
	}

	if err := rows5.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error en el bucle de lectura rows5 >> %v", err)
	}

	consulta_capital.TotalDeudasPorFiador = total_deudas_fiador

	q = `
		SELECT 
			capital_pago, 
			capital_pago*?/100 AS capital_disponible
		FROM usuarios WHERE id = ?`

	if err := s.Conexion.QueryRow(q, s.Ajustes.Prestamos.CapitalUsable, id).Scan(&consulta_capital.CapitalPago, &consulta_capital.CapitalDisponible); err != nil {
		return nil, fmt.Errorf("No se pudo obtener el capital del usuario %d >> %v", id, err)
	}

	consulta_capital.TotalDisponible = consulta_capital.CapitalDisponible - total_deudas_fiador - total_deuda_prestamos

	user.CapitalDisponible = consulta_capital

	q = `SELECT codigo, codigo_prestamo, monto_pago, deuda_antes, deuda_despues, fecha FROM cheques_de_prestamos WHERE id_user = ? ORDER BY codigo DESC`

	rows6, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener los cheques del usuario %d >> %v", id, err)
	}
	defer rows6.Close()

	var fecha_cheque time.Time

	for rows6.Next() {
		var cheque ChequeJson
		if err := rows6.Scan(&cheque.Codigo, &cheque.CodigoPrestamo, &cheque.Monto, &cheque.DeudaAntes, &cheque.DeudaDespues, &fecha_cheque); err != nil {
			return nil, fmt.Errorf("No se pudieron escanear los cheques del usuario %d >> %v", id, err)
		}
		cheque.Fecha = fecha_cheque.Format("2006/01/02 15:04")
		user.ChequesUser = append(user.ChequesUser, cheque.Codigo)
		user.Cheques = append(user.Cheques, &cheque)
	}

	if err := rows6.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error en el bucle de lectura rows6 >> %v", err)
	}

	return &user, nil
}

func (s *centralStruct) storeObtenerDatosChequePrestamos(codigo int) (*ChequePdf, error) {
	cheque := &ChequePdf{Codigo: codigo}

	var fecha time.Time

	q := `
		SELECT 
			id_user, 
			codigo_prestamo, 
			monto_pago, 
			deuda_antes, 
			deuda_despues, 
			fecha,
			tipo,
			intereses_antes,
			intereses_despues
		FROM cheques_de_prestamos 
		WHERE codigo = ?
	`

	if err := s.Conexion.QueryRow(q, codigo).Scan(
		&cheque.Id,
		&cheque.CodigoPrestamo,
		&cheque.MontoAbonado,
		&cheque.DeudaAntes,
		&cheque.DeudaDespues,
		&fecha,
		&cheque.AbonoA,
		&cheque.InteresesAntes,
		&cheque.InteresesDespues,
	); err != nil {
		return nil, fmt.Errorf("No se pudo obtener el cheque de prestamos %d >> %v", codigo, err)
	}

	cheque.FechaYHora = fecha.Format("2006/01/02 15:04")

	q = `SELECT nombre FROM usuarios WHERE id = ?`
	if err := s.Conexion.QueryRow(q, cheque.Id).Scan(&cheque.Nombre); err != nil {
		return nil, fmt.Errorf("No se pudo obtener el nombre del usuario %d >> %v", cheque.Id, err)
	}

	return cheque, nil
}
