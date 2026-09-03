package analizar

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (s *centralStruct) EstadoDeBoletas() (string, error) {
	formatoParaGuardar := FormatoDeConsultaCompletaDeBoletas{
		RifaAct: s.Ajustes.Rifas.RifaAct,
	}

	// obtenemos todos los datos
	q := `
		SELECT DISTINCT edb.entrega_a 
		FROM entregas_de_boletas edb 
		WHERE edb.en_rifa = ?
		ORDER BY edb.entrega_a;
	`
	rows1, err := s.Conexion.Query(q, s.Ajustes.Rifas.RifaAct)
	if err != nil {
		return "", fmt.Errorf("No se pudo buscar todos los usuarios a los que se les entrego boletas >> %v", err)
	}
	defer rows1.Close()

	// obtenemos los usuarios
	usuarios := []int{}
	var usuario int
	for rows1.Next() {
		if err := rows1.Scan(&usuario); err != nil {
			return "", fmt.Errorf("No se pudo obtener el usuario")
		}
		usuarios = append(usuarios, usuario)
	}

	if err := rows1.Err(); err != nil {
		return "", fmt.Errorf("Hubo un error en la ejecucion del bucle >> %v", err)
	}

	// obtenemos las entregas
	q = `
		SELECT 
			edb.codigo,
			edb.entrega_a,
			edb.fecha 
		FROM entregas_de_boletas edb 
		WHERE edb.en_rifa = ?
		ORDER BY edb.codigo;
	`
	rows2, err := s.Conexion.Query(q, s.Ajustes.Rifas.RifaAct)
	if err != nil {
		return "", fmt.Errorf("No se pudo obtener las entregas hechas >> %v", err)
	}
	defer rows2.Close()

	entregas := map[int][]FormatoDeEntregas{}
	var (
		codigo    int
		entrega_a int
		fecha     time.Time
	)

	for rows2.Next() {
		if err := rows2.Scan(&codigo, &entrega_a, &fecha); err != nil {
			return "", fmt.Errorf("No se pudo escanear las entregas >> %v", err)
		}

		entregas[entrega_a] = append(entregas[entrega_a], FormatoDeEntregas{
			Codigo: codigo,
			Fecha:  fecha.Format("2006-01-02"),
		})
	}

	if err := rows2.Err(); err != nil {
		return "", fmt.Errorf("Hubo un error en la ejecucion del bucle >> %v", err)
	}

	// obteneos todas las boletas
	q = `
		SELECT 
			br.boleta,
			br.entrega
		FROM boletas_rifas br 
		WHERE
			br.codigo_rifa = ? AND 
			br.entrega IS NOT NULL
		ORDER BY br.codigo_boleta;
	`
	rows3, err := s.Conexion.Query(q, s.Ajustes.Rifas.RifaAct)
	if err != nil {
		return "", fmt.Errorf("No se pudo obtener las boletas >> %v", err)
	}
	defer rows3.Close()

	boletas := map[int][]string{}
	var (
		boleta  string
		entrega int
	)
	for rows3.Next() {
		if err := rows3.Scan(&boleta, &entrega); err != nil {
			return "", fmt.Errorf("No se pudo escanear las boletas >> %v", err)
		}
		boletas[entrega] = append(boletas[entrega], boleta)
	}

	if err := rows3.Err(); err != nil {
		return "", fmt.Errorf("Hubo un error en la ejecucion del bucle 3 >> %v", err)
	}

	// creamos todos los datos
	q = `SELECT u.nombre FROM usuarios u WHERE u.id = ?;`
	var nombre string

	for _, usuario := range usuarios {
		if err := s.Conexion.QueryRow(q, usuario).Scan(&nombre); err != nil {
			return "", fmt.Errorf("No se pudo obtener el nombre del usuario %d >> %v", usuario, err)
		}
		formatoUsuario := FormatoDeUsuarioYBoletas{
			Puesto: usuario,
			Nombre: nombre,
		}

		for _, entrega := range entregas[usuario] {
			formatoUsuario.Entregas = append(formatoUsuario.Entregas, FormatoDeFechaYBoletas{
				Fecha:   entrega.Fecha,
				Boletas: boletas[entrega.Codigo],
			})
		}

		formatoParaGuardar.Usuarios = append(formatoParaGuardar.Usuarios, formatoUsuario)
	}

	// obtenemos las boletas disponibles
	q = `
		SELECT br.boleta 
		FROM boletas_rifas br 
		WHERE 
			br.codigo_rifa = ? AND 
			br.entrega IS NULL;
	`
	rows4, err := s.Conexion.Query(q, s.Ajustes.Rifas.RifaAct)
	if err != nil {
		return "", fmt.Errorf("No se pudo obtener las boletas disponibles >> %v", err)
	}
	defer rows4.Close()

	boletasLibres := []string{}
	var boletaLibre string

	for rows4.Next() {
		if err := rows4.Scan(&boletaLibre); err != nil {
			return "", fmt.Errorf("No se pudo escanear las boletas libres >> %v", err)
		}
		boletasLibres = append(boletasLibres, boletaLibre)
	}

	if err := rows4.Err(); err != nil {
		return "", fmt.Errorf("Hubo un error en la ejecucion del bucle 4 >> %v", err)
	}

	formatoParaGuardar.BoletasLibres = boletasLibres

	basePath, _ := os.Getwd()

	// creamos el archivo json
	nombreJson := time.Now().Format("0102150405") + ".json"
	rutaJson := filepath.Join(basePath, "tmp", nombreJson)

	file, _ := os.Create(rutaJson)
	defer file.Close()

	json.NewEncoder(file).Encode(formatoParaGuardar)

	// compilamos el pdf
	nombrePdf := time.Now().Format("0102150405") + ".pdf"
	rutaPdf := filepath.Join(basePath, "tmp", nombrePdf)

	formato := fmt.Sprintf(FORMATO_BOLETAS_GENERALES, "tmp/"+nombreJson)

	cmd := exec.Command("typst", "compile", "-", rutaPdf)
	cmd.Stdin = strings.NewReader(formato)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("no se pudo compilar el archivo (%s): %v\nSalida: %s", rutaPdf, err, string(output))
	}

	return nombrePdf, nil
}
