package modificar

import (
	"API-FONDO/pkg/prestamos"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (s *centralStruct) storeCrearUsuario(nuevoUsuario *FormatoCrearUsuario) (int, error) {
	// iniciamos una transaccion para guardar el usuario
	tx, err := s.Conexion.Begin()
	if err != nil {
		return -1, errors.New("No se pudo iniciar una transacción para guardar el usuario")
	}
	defer tx.Rollback()

	valorCuota := s.Ajustes.Cuotas.ValorCuota
	capitalProyectado := 50 * nuevoUsuario.Puestos * valorCuota

	usuarios := `INSERT INTO usuarios (nombre, puestos, telefono, capital_proyectado) VALUES (?, ?, ?, ?)`

	resultado, err := tx.Exec(usuarios, nuevoUsuario.Nombre, nuevoUsuario.Puestos, nuevoUsuario.Telefono, capitalProyectado)
	if err != nil {
		return -1, errors.New("No se pudo guardar al usuario y obtener su ID")
	}

	id, err := resultado.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("No se puede crear usuario >> %q", err)
	}

	err = s.Ajustes.IncrementarUsuarios()
	if err != nil {
		return -1, fmt.Errorf("No se puede incrementar el contador de usuarios >> %q", err)
	}

	// Todo ha salido bien
	err = tx.Commit()
	if err != nil {
		return -1, errors.New("No se pudo guardar el usuario, la transaccion no puede ser finalizada")
	}

	return int(id), nil
}

func (s *centralStruct) storeObtenerUsuarios() ([]*Usuario, error) {
	q := `SELECT u.id, u.nombre, u.telefono, u.puestos FROM usuarios u ORDER BY u.id DESC`
	rows, err := s.Conexion.Query(q)
	if err != nil {
		return nil, fmt.Errorf("No se pueden consultar los usuarios: %v", err)
	}
	defer rows.Close()

	var (
		id       int
		nombre   string
		telefono string
		puestos  int
	)

	users := []*Usuario{}

	for rows.Next() {
		if err = rows.Scan(&id, &nombre, &telefono, &puestos); err != nil {
			return nil, errors.New("No se pueden obtener los datos del usuario")
		}

		users = append(users, &Usuario{
			Id:       id,
			Nombre:   nombre,
			Telefono: telefono,
			Puestos:  puestos,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Hubo un error con la ejecucion del bucle rows >> %v", err)
	}

	return users, nil
}

func (s *centralStruct) storeCambiarFechasDeUnPrestamo(codigo int, nuevaFecha string) error {
	fecha, err := time.Parse("2006-01-02", nuevaFecha)
	if err != nil {
		return fmt.Errorf("No se pudo parsear la fecha >> %v", err)
	}

	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("No se pudo inicar la transaccion para escribir los datos >> %v", err)
	}
	defer tx.Rollback()

	q := `DELETE FROM prestamos_fechas_de_pago WHERE codigo_prestamo = ?`

	result, err := tx.Exec(q, codigo)
	if err != nil {
		return fmt.Errorf("No se pudo eliminar todas las fechas ya habidas >> %v", err)
	}

	cols, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("No se pudo saber cuantas columnas fueron eliminadas >> %v", err)
	}

	if cols == 0 {
		return errors.New("Me temo que no existe un prestamo con este codigo")
	}

	q = `INSERT INTO prestamos_fechas_de_pago (codigo_prestamo, fecha) VALUES (?,?)`
	stmt, err := tx.Prepare(q)
	if err != nil {
		return fmt.Errorf("No se pudo precompilar la consulta >> %q", err)
	}
	defer stmt.Close()

	for _, fecha := range prestamos.CalendarioMeses(fecha, s.Ajustes.Prestamos.FechaDeCierre) {
		if _, err := stmt.Exec(codigo, fecha); err != nil {
			return fmt.Errorf("No se pudo reescribir las fechas de pago del prestamo >> %v", err)
		}
	}

	// hacer la nota del prestamo
	q = `INSERT INTO notas_de_prestamos (codigo_prestamo, anotacion) VALUES (?,?)`
	msg := fmt.Sprintf("Se corrige la fecha de creacion a: %s", nuevaFecha)
	if _, err := tx.Exec(q, codigo, msg); err != nil {
		return fmt.Errorf("No se pudo corregir la anotacion del prestamo >> %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("No se pudo cerrar la transaccion >> %v", err)
	}

	return nil
}
