package ajustes_http

import (
	"fmt"
)

func (s *centralStruct) storeObtenerAjusteNum(ajuste string) (int, error) {
	ajustesNum := map[string]int{
		"valor_cuota":    s.Ajustes.Cuotas.ValorCuota,
		"valor_multa":    s.Ajustes.Cuotas.ValorMulta,
		"interes_M_tope": s.Ajustes.Intereses.InteresMTope,
		"interes_m_tope": s.Ajustes.Intereses.InteresmTope,
		"tope_intereses": s.Ajustes.Intereses.TopeIntereses,
		"capital_usable": s.Ajustes.Prestamos.CapitalUsable,
		"rifa_act":       s.Ajustes.Rifas.RifaAct,
		"count_users":    s.Ajustes.Usuarios.Count,
	}

	ajustesValidosStr := "valor_cuota, valor_multa, interes_M_tope, interes_m_tope, tope_intereses, capital_usable, rifa_act, count_users"

	valor, ok := ajustesNum[ajuste]
	if !ok {
		return -1, fmt.Errorf("El ajuste numerico '%s' no existe, las opciones validas son: %s", ajuste, ajustesValidosStr)
	}

	return valor, nil
}

func (s *centralStruct) storeObtenerAjusteBool(ajuste string) (bool, error) {
	ajustesBool := map[string]bool{
		"cobrar_multas":     s.Ajustes.Cuotas.CobrarMultas,
		"bloquear_usuarios": s.Ajustes.Cuotas.BloquearUsuarios,
	}

	ajustesValidosBool := "cobrar_multas, bloquear_usuarios"

	valor, ok := ajustesBool[ajuste]
	if !ok {
		return false, fmt.Errorf("El ajuste boleano '%s' no existe, las opciones validas son: %s", ajuste, ajustesValidosBool)
	}

	return valor, nil
}

func (s *centralStruct) storeCrearCalendario(formato FormatNewCalendar) error {
	tx, err := s.Conexion.Begin()
	if err != nil {
		return fmt.Errorf("No se pudo crear la transaccion para la base de datos >> %q", err)
	}
	defer tx.Rollback()

	q := `delete from calendario`
	if _, err := tx.Exec(q); err != nil {
		return fmt.Errorf("No se pudo eliminal en calendario anterior >> %q", err)
	}

	fechasCalendario, err := fnCrearCalendario(formato.FechaInicial, formato.HoraDeCierre, formato.FechasDobles)
	if err != nil {
		return fmt.Errorf("No se pudo formatear las fechas dadas >> %q", err)
	}

	stmt, err := tx.Prepare(`insert into calendario(semana, fecha) values (?, ?);`)
	if err != nil {
		return fmt.Errorf("No se pudo precompilar la consulta de insercion >> %q", err)
	}
	defer stmt.Close()

	formatoDeFechas := "2006-01-02 15:04"
	var fechaTmp string

	for idx, fecha := range fechasCalendario {
		fechaTmp = fecha.Format(formatoDeFechas)
		if _, err = stmt.Exec(idx+1, fechaTmp); err != nil {
			return fmt.Errorf("no se pudo guardar la fecha '%s' >> %q", fechaTmp, err)
		}
	}

	return tx.Commit()
}
