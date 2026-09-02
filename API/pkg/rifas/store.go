/*
	En esta seccion en teoria deberia manejar todo lo relacionado con rifas unicamente pero
	la parte de la gestion de boletas como eliminarlas, crearlas y demas son cosas muy
	sencillas que podria manejar aca sin problema, asi que esta seccion es mixta entre rifas
	la creacion de una rifa, el cierra de una rifa y la gestion de las boletas
*/

package rifas

import (
	"errors"
	"fmt"
)

// Creacion de una rifa

func (s *centralStruct) storeCrearUnaRifa(frmt *FormatoDatosDeRifa) (int64, error) {
	ganancia_por_boleta := (frmt.NumeroDeBoletas * frmt.CostoDeBoletas) - frmt.CostosDeAdministracion
	for _, value := range frmt.Premios {
		ganancia_por_boleta -= value.Premio
	}
	ganancia_por_boleta /= frmt.NumeroDeBoletas

	if ganancia_por_boleta < 1 {
		return -1, errors.New("rifas >> crear-una-rifa >> la ganacia por boleta no puede ser nula o negativa")
	}

	tx, err := s.Conexion.Begin()
	if err != nil {
		return -1, errors.New("rifas >> crear-una-rifa >> no se pudo inicar la transaccion")
	}
	defer tx.Rollback()

	q := `
		insert into datos_de_rifas(
			numero_de_boletas, costo_de_boletas, costos_de_administracion, ganancias_por_boleta
		) values (?, ?, ?, ?)
	`

	result, err := tx.Exec(q, frmt.NumeroDeBoletas, frmt.CostoDeBoletas, frmt.CostosDeAdministracion, ganancia_por_boleta)
	if err != nil {
		return -1, fmt.Errorf("rifas >> crear-una-rifa >> no se pudo insertar el usuario >> %q", err)
	}

	newCodigoRifa, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("rifas >> no se pudo obtener el numero de la rifa")
	}

	q = `insert into premios_de_rifas(codigo_rifa, posicion, premio) values (?, ?, ?)`

	for _, premio := range frmt.Premios {
		if _, err := tx.Exec(q, newCodigoRifa, premio.Numero, premio.Premio); err != nil {
			return -1, fmt.Errorf("rifas >> crear-una-rifa >> no se pudo guardar el premio %d, codigo de rifa: %d >> %q", premio, newCodigoRifa, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("rifas >> crear-una-rifa >> no se pudo hacer un commit de la operacion >> %q", err)
	}

	return newCodigoRifa, nil
}

// Iniciamos con las funciones de insertar

func (s *centralStruct) storeInsertarBoletas(boletas []Boleta) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("rifas >> insertar boletas >> no se pudo iniciar la transaccion: %q", err)
	}
	defer tx.Rollback()

	q := `insert into boletas_rifas(codigo_rifa, codigo_boleta, boleta) values (?, ?, ?)`
	for _, boleta := range boletas {
		if _, err := tx.Exec(q, s.Ajustes.Rifas.RifaAct, boleta.Codigo, boleta.Boleta); err != nil {
			return fmt.Errorf("rifas >> insertar boletas >> no se pudo insertar la boleta (%d) >> %q", boleta.Codigo, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("rifas >> insertar boletas >> no se pudo hacer el commit: %q", err)
	}

	return nil
}

func (s *centralStruct) storeBoletasLibres() ([]Boleta, error) {
	q := `select br.codigo_boleta, br.boleta from boletas_rifas br where br.dada_a is null and br.codigo_rifa = ? order by br.codigo_boleta`
	rows, err := s.Conexion.Query(q, s.Ajustes.Rifas.RifaAct)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener boletas las boletas libres >> %q", err)
	}
	defer rows.Close()

	var (
		codigo  int
		boleta  string
		boletas []Boleta
	)
	for rows.Next() {
		if err := rows.Scan(&codigo, &boleta); err != nil {
			return nil, fmt.Errorf("No se pudo obtner la boleta libre >> %q", err)
		}

		boletas = append(boletas, Boleta{Codigo: codigo, Boleta: boleta})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("No se pudieron leer todas las boletas disponibles >> %v", err)
	}

	return boletas, nil
}

func (s *centralStruct) storeEliminarBoleta(codigo int) error {
	q := `delete from boletas_rifas where codigo_rifa = ? and codigo_boleta = ?`
	result, err := s.Conexion.Exec(q, s.Ajustes.Rifas.RifaAct, codigo)
	if err != nil {
		return fmt.Errorf("No se pudo eliminar la boleta >> %q", err)
	}

	delRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("No se sabe cuantas filas fueron afectadas")
	}

	if delRows == 0 {
		return fmt.Errorf(
			"no existe una boleta en la rifa '%d' con el identificador '%d'",
			s.Ajustes.Rifas.RifaAct,
			codigo,
		)
	}

	if delRows > 1 {
		return fmt.Errorf(
			"como putas habia %d boletas repetidas en la rifa '%d'? la boleta '%d' es incorrecta, por favor llamar a fernando !!!",
			delRows,
			s.Ajustes.Rifas.RifaAct,
			codigo,
		)
	}

	return nil
}

func (s *centralStruct) storeObtenerDatosDelUsuario(id int) (*FormatoGet, error) {
	user := &FormatoGet{Puesto: id}

	// Nombre
	q := `select nombre from usuarios where id = ?`
	if err := s.Conexion.QueryRow(q, id).Scan(&user.Nombre); err != nil {
		return nil, fmt.Errorf("No se pudo consultar el nombre del usuario >> %v", err)
	}

	// Boletas Disponibles
	q = `select br.codigo_boleta from boletas_rifas br where br.dada_a is null and br.codigo_rifa = ? order by br.codigo_boleta`
	rows, err := s.Conexion.Query(q, s.Ajustes.Rifas.RifaAct)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener boletas las boletas libres >> %q", err)
	}
	defer rows.Close()

	var codigo int

	for rows.Next() {
		if err := rows.Scan(&codigo); err != nil {
			return nil, fmt.Errorf("No se pudo obtner la boleta libre >> %q", err)
		}

		user.BoletasDisponibles = append(user.BoletasDisponibles, codigo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("No se pudieron leer todas las boletas disponibles >> %v", err)
	}

	// Deuda por Boletas
	q = `select deuda_rifas from usuarios where id = ?`
	if err := s.Conexion.QueryRow(q, id).Scan(&user.DeudaPorBoletas); err != nil {
		return nil, fmt.Errorf("No se pudo obtener la deuda por boletas del usuario >> %v", err)
	}

	// Boletas Entregadas
	q = `SELECT br.boleta FROM boletas_rifas br WHERE br.dada_a = ? order by br.codigo_boleta`
	rows1, err := s.Conexion.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener las boletas dadas al usuario >> %v", err)
	}
	defer rows1.Close()

	var boleta string

	for rows1.Next() {
		if err := rows1.Scan(&boleta); err != nil {
			return nil, fmt.Errorf("No se pudo leer la boleta dentro del bucle >> %v", err)
		}

		user.BoletasEntregadas = append(user.BoletasEntregadas, boleta)
	}

	if err := rows1.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error al escaneal los balores dentro del bucle rows1 >> %v", err)
	}

	// retornal los valores
	return user, nil
}

func (s *centralStruct) storeObtenerRifasCreadas() ([]*FormatoDatosDeRifa, error) {
	rifas := []*FormatoDatosDeRifa{}

	q := `select codigo, numero_de_boletas, costo_de_boletas, costos_de_administracion, ganancias_por_boleta from datos_de_rifas`
	rows, err := s.Conexion.Query(q)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener los datos de las rifas >> %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		rifa := &FormatoDatosDeRifa{}

		if err := rows.Scan(&rifa.Codigo, &rifa.NumeroDeBoletas, &rifa.CostoDeBoletas, &rifa.CostosDeAdministracion, &rifa.GananciaPorBoleta); err != nil {
			return nil, fmt.Errorf("No se pudo copiar los datos a cada rifa >> %v", err)
		}

		rifas = append(rifas, rifa)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error en la ejecucion del bucle rows >> %v", err)
	}

	q = `select codigo_rifa, posicion, premio from premios_de_rifas`
	rows1, err := s.Conexion.Query(q)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener los premios de las rifas >> %v", err)
	}
	defer rows1.Close()

	premios := map[int][]Premio{}
	var codigo int

	for rows1.Next() {
		premio := Premio{}
		if err := rows1.Scan(&codigo, &premio.Numero, &premio.Premio); err != nil {
			return nil, fmt.Errorf("No se pudo leer los premios en el bucle >> %v", err)
		}

		premios[codigo] = append(premios[codigo], premio)
	}

	if err := rows1.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error en el bucle rows1 >> %v", err)
	}

	for _, rifa := range rifas {
		rifa.Premios = premios[rifa.Codigo]
	}

	return rifas, nil
}

func (s *centralStruct) storeEntregarBoletas(id int, boletas []int) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("No se pudo crear la transaccion para la operacion >> %v", err)
	}
	defer tx.Rollback()

	// sumamos el valor de todas las boletas a la deuda del usuario
	q := `
		update usuarios
		SET 
			deuda_rifas = deuda_rifas + (
				SELECT 
					ddr.costo_de_boletas 
				FROM datos_de_rifas ddr 
				WHERE ddr.codigo = ?
			)*?
		WHERE id = ?
	`
	if _, err := tx.Exec(q, s.Ajustes.Rifas.RifaAct, len(boletas), id); err != nil {
		return fmt.Errorf("No se pudo actualizar la deuda por las boletas >> %v", err)
	}

	// entregamos las boletas
	q = `
		UPDATE boletas_rifas
		SET 
			dada_a = ?
		WHERE 
			codigo_rifa = ? AND codigo_boleta = ?
	`
	for _, boleta := range boletas {
		if _, err := tx.Exec(q, id, s.Ajustes.Rifas.RifaAct, boleta); err != nil {
			return fmt.Errorf("No se puede dar la boleta %d al usuario >> %v", boleta, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("No se pudo guardar todos los cambio hechos >> %v", err)
	}

	return nil
}

func (s *centralStruct) storeObtenerDeudasPorBoletas(id int) (int, error) {
	q := `select deuda_rifas from usuarios where id = ?`
	var deuda int
	if err := s.Conexion.QueryRow(q, id).Scan(&deuda); err != nil {
		return -1, fmt.Errorf("No se pudo obtener la deuda del usuario >> %v", err)
	}

	return deuda, nil
}

func (s *centralStruct) storePagarDeudaPorBoletas(id, monto int) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("No se pudo iniciar la transaccion para realizar el pago >> %v", err)
	}
	defer tx.Rollback()

	// hecemos el cheque

	q := `select deuda_rifas from usuarios where id = ?`

	var (
		deudaAntes   int
		deudaDespues int
	)

	if err := tx.QueryRow(q, id).Scan(&deudaAntes); err != nil {
		return fmt.Errorf("No se pudo consultar la deuda del usuario >> %v", err)
	}

	deudaDespues = deudaAntes - monto

	q = `
		INSERT INTO cheques_de_rifas (
			id_user, monto_pago, deuda_antes, deuda_despues, en_rifa
		) VALUES (?,?,?,?,?)
	`

	if _, err := tx.Exec(q, id, monto, deudaAntes, deudaDespues, s.Ajustes.Rifas.RifaAct); err != nil {
		return fmt.Errorf("No se pudo crear el cheque de pago >> %v", err)
	}

	q = `
		update usuarios
		set
			deuda_rifas = deuda_rifas - ?
		where
			id = ? 
	`
	if _, err := tx.Exec(q, monto, id); err != nil {
		return fmt.Errorf("No se pudo descontar la deuda del usuario >> %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("No se pudo guardar toda la operacion >> %v", err)
	}

	return nil
}
