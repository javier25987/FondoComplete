package verusuaios

import (
	"API-FONDO/pkg/prestamos"
	"fmt"
	"time"
)

func (s *centralStruct) storeBuscarUsuarios(llave string) ([]*Usuario, error) {
	to_search := fmt.Sprintf("%%%s%%", llave)

	q := `
		SELECT 
			u.id, 
			u.nombre, 
			u.telefono, 
			u.puestos, 
			u.estado 
		FROM usuarios u 
		WHERE u.nombre 
		LIKE ? ORDER BY u.id
	`

	rows, err := s.Conexion.Query(q, to_search)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener los usuarios >> %v", err)
	}
	defer rows.Close()

	users := []*Usuario{}

	var estado bool

	for rows.Next() {
		user := &Usuario{}

		if err := rows.Scan(&user.Id, &user.Nombre, &user.Telefono, &user.Puestos, &estado); err != nil {
			return nil, fmt.Errorf("No se pudo obtener los datos de un usuario >> %v", err)
		}

		if estado {
			user.Estado = "✅ Activo"
		} else {
			user.Estado = "🚨 Inactivo"
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *centralStruct) storeFiltrarAcuerdo() ([]*Acuerdo, error) {
	q := `
		SELECT 
		u.id, 
		u.nombre, 
		u.capital_pago, 
		u.capital_retirado, 
		u.capital_pago/2 - u.capital_retirado
		FROM usuarios u 
		WHERE u.capital_retirado < u.capital_pago/2 
		ORDER BY u.id
	`

	rows, err := s.Conexion.Query(q)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obter los datos del usuario >> %v", err)
	}
	defer rows.Close()

	acuerdos := []*Acuerdo{}

	for rows.Next() {
		acuerdo := &Acuerdo{}

		if err := rows.Scan(&acuerdo.Id, &acuerdo.Nombre, &acuerdo.CapitalPago, &acuerdo.CapitalRetirado, &acuerdo.NecesarioRetirar); err != nil {
			return nil, fmt.Errorf("No se pudo obtener los datos de un usuario >> %v", err)
		}

		acuerdos = append(acuerdos, acuerdo)
	}

	return acuerdos, nil
}

func (s *centralStruct) storeBuscarBoletas(codigo_rifa int, numero string) ([]*Boleta, error) {
	q := `
		SELECT 
			br.dada_a, 
			u.nombre, 
			br.boleta 
		FROM boletas_rifas br  
		JOIN usuarios u 
			ON u.id = br.dada_a 
		WHERE br.boleta 
			LIKE ? AND br.codigo_rifa = ?`

	format_numero := fmt.Sprintf("%%%s%%", numero)

	rows, err := s.Conexion.Query(q, format_numero, codigo_rifa)
	if err != nil {
		return nil, fmt.Errorf("No se pudo obtener los datos de los usuarios >> %v", err)
	}

	boletas := []*Boleta{}

	for rows.Next() {
		boleta := &Boleta{}

		if err := rows.Scan(&boleta.Id, &boleta.Nombre, &boleta.Numeros); err != nil {
			return nil, fmt.Errorf("No se pudo obtener los datos de un usario >> %v", err)
		}

		boletas = append(boletas, boleta)
	}

	return boletas, nil
}

func (s *centralStruct) storeObtenerPrestamo(codigo int) (*GetPrestamo, error) {
	prestamo := &prestamos.Prestamo{}

	var (
		id_user int
		nombre  string

		id_user_fiador int
		deuda_fiador   int

		fecha_de_pago time.Time

		estado bool

		fecha_anotacion time.Time
		anotacion       string
		nota            string
	)

	q := `
		SELECT 
			codigo, 
			id_user, 
			interes, 
			intereses_vencidos, 
			intereses_generados, 
			deuda, 
			monto, 
			estado, 
			motivo
		FROM prestamos_hechos 
		WHERE codigo = ?
	`

	err := s.Conexion.QueryRow(q, codigo).Scan(
		&prestamo.Codigo,
		&id_user,
		&prestamo.Interes,
		&prestamo.InteresVencido,
		&prestamo.InteresGenerado,
		&prestamo.Deuda,
		&prestamo.Monto,
		&estado,
		&prestamo.Motivo,
	)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener los datos del prestamo >> %v", err)
	}

	prestamo.DeudaActual = prestamo.Deuda + prestamo.InteresVencido

	prestamo.Estado = "PAGO ✅"
	if estado {
		prestamo.Estado = "ACTIVO ⏳"
	}

	// cargar fiadores, fechas de pago y anotaciones
	q = `SELECT id_user_fiador, monto_respaldado FROM prestamos_fiadores WHERE codigo_prestamo = ?;`
	p := `SELECT fecha FROM prestamos_fechas_de_pago WHERE codigo_prestamo = ? ORDER BY fecha ASC`
	k := `SELECT anotacion, fecha FROM notas_de_prestamos WHERE codigo_prestamo = ? ORDER BY fecha ASC`

	// fiadores
	rows, err := s.Conexion.Query(q, prestamo.Codigo)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener los fiadores del prestamo %d >> %v", prestamo.Codigo, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id_user_fiador, &deuda_fiador); err != nil {
			return nil, fmt.Errorf("No se pudieron obtener los fiadores del prestamo %d >> %v", prestamo.Codigo, err)
		}
		prestamo.Fiadores = append(prestamo.Fiadores, &prestamos.Fiador{Id_user: id_user_fiador, Deuda: deuda_fiador})
	}
	rows.Close()

	// fechas de pago
	rows, err = s.Conexion.Query(p, prestamo.Codigo)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener las fechas de pago del prestamo %d >> %v", prestamo.Codigo, err)
	}

	for rows.Next() {
		if err := rows.Scan(&fecha_de_pago); err != nil {
			return nil, fmt.Errorf("No se pudo obtener la fecha del prestamo %d >> %v", prestamo.Codigo, err)
		}
		prestamo.FechasPago = append(prestamo.FechasPago, fecha_de_pago.Format("2006/01/02"))
	}
	rows.Close()

	// anotaciones
	rows, err = s.Conexion.Query(k, prestamo.Codigo)
	if err != nil {
		return nil, fmt.Errorf("No se pudieron obtener las anotaciones del prestamo %d >> %v", prestamo.Codigo, err)
	}

	for rows.Next() {
		if err := rows.Scan(&anotacion, &fecha_anotacion); err != nil {
			return nil, fmt.Errorf("No se pudo obtener las anotaciones del prestamo %d >> %v", prestamo.Codigo, err)
		}
		nota = fmt.Sprintf("[%s] %s", fecha_anotacion.Format("2006/01/02"), anotacion)
		prestamo.Anotaciones = append(prestamo.Anotaciones, nota)
	}
	rows.Close()

	q = `SELECT nombre FROM usuarios WHERE id = ?`
	if err := s.Conexion.QueryRow(q, id_user).Scan(&nombre); err != nil {
		return nil, fmt.Errorf("No se pudo consultar el nombre del usuario >> %v", err)
	}

	return &GetPrestamo{Id: id_user, Nombre: nombre, Prestamo: prestamo}, nil
}
