package main

import (
	"API-FONDO/ajustes"
	ajustes_http "API-FONDO/pkg/ajustes"
	"API-FONDO/pkg/cuotas"
	"API-FONDO/pkg/deudas"
	"API-FONDO/pkg/modificar"
	"API-FONDO/pkg/prestamos"
	"API-FONDO/pkg/rifas"
	verusuaios "API-FONDO/pkg/ver_usuarios"
	"database/sql"
	"log"
	"net/http"
	"os"

	"fmt"

	_ "modernc.org/sqlite"
)

func main() {

	/*
		Limpiamos la carpeta `tmp` eliminandola y volviendola a crear
	*/

	if err := os.RemoveAll("tmp"); err != nil {
		log.Fatalln("No se pudo eliminar la carpeta `tmp`")
	}
	if err := os.Mkdir("tmp", 0755); err != nil {
		log.Fatalln("No se pudo crear la carpeta `tmp`")
	}

	/*
		Esta es una nueva aplicacion del fondo la cual separa los ajustes de los datos y
		usa como base de datos SQLite en ves de postgres
	*/

	// Iniciamos abriendo los ajustes
	progConfig, err := ajustes.AbrirAjustes()
	if err != nil {
		ajustes.CrearAjustes()
		fmt.Println("Unos nuevos ajustes han sido creados. Error producido:")
		fmt.Println(err)
	}

	/*
		Abrrimos la conexion a la base de datos
	*/

	db, err := sql.Open("sqlite", progConfig.BaseDeDatos.Db)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.Exec("PRAGMA foreign_keys = ON;")

	/*
		Obtenemos los handlers de los paquetes con los que vamos a trajar
	*/

	hModificar := modificar.NewCentralStruct(progConfig, db)
	hAjustes := ajustes_http.NewCentralStruct(progConfig, db)
	hCuotas := cuotas.NewCentralStruct(progConfig, db)
	// hAnotaciones := anotaciones.HandlerAnotaciones(MainCtx, pool)
	// hRegistros := registros.HandlerRegistros(MainCtx, pool)
	hPrestamos := prestamos.NewCentralStruct(progConfig, db)
	hDeudas := deudas.NewCentralStruct(progConfig, db)
	hVer := verusuaios.NewCentralStruct(progConfig, db)
	hRifas := rifas.NewCentralStruct(progConfig, db)

	/*
		Creamos todos los endpoints necesarios para el sistema
	*/

	// Cuotas
	http.HandleFunc("/cuotas/{id}", hCuotas.RutaCuotas_)
	http.HandleFunc("/cuotas/cheque/{codigo}", hCuotas.RutaCuotasCheque_)
	// http.HandleFunc("/cuotas/bloc/{semana}", void)

	// Prestamos
	http.HandleFunc("/prestamos/{id}", hPrestamos.RutaPrestamos_)
	http.HandleFunc("/prestamos/cheque/{codigo}", hPrestamos.RutaPrestamosCheque_)
	http.HandleFunc("/prestamos/carta", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./documents/carta.pdf")
	})

	// // Analisis
	// http.HandleFunc("/analizar", void)
	// http.HandleFunc("/analizar/", void)

	// // Transferencias
	// http.HandleFunc("/transferencias", void)
	// http.HandleFunc("/transferencias/", void)

	// Rifas
	http.HandleFunc("/rifas/{id}", hRifas.RutaRifas)

	// Deudas
	http.HandleFunc("/deudas/{id}", hDeudas.RutaDeudas_)

	// // Anotaciones
	// http.HandleFunc("/anotaciones/{id}", HAnotaciones.Anotaciones_)

	// Ver usuarios
	http.HandleFunc("/ver/usuarios", hVer.RutaVerUsuarios)
	http.HandleFunc("/ver/acuerdo", hVer.RutaVerAcuerdo)
	http.HandleFunc("/ver/boletas/{rifa}", hVer.RutaVerBoletas)
	http.HandleFunc("/ver/prestamo/{codigo}", hVer.RutaVerPrestamo)

	// // Registros
	// http.HandleFunc("/registros/{mes}", HRegistros.Registros_)

	// Insertar
	http.HandleFunc("/insertar", hRifas.RutaInsertar)

	// Modificar Usuarios
	http.HandleFunc("/modificar", hModificar.RutaModificar)
	http.HandleFunc("/modificar/fechas", hModificar.RutaModificarFechas)
	// http.HandleFunc("/modificar/semana", void)

	// Ajustes
	http.HandleFunc("/ajustes", hAjustes.RutaAjustes)
	http.HandleFunc("/ajuste/calendario", hAjustes.RutaAjusteCalendario)
	// http.HandleFunc("/ajuste", HAjustes.Ajuste)
	http.HandleFunc("/ajustes/boletas", hRifas.RutaAjustesBoletas)

	// http.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
	// 	// Reemplaza "./archivo.html" con la ruta a tu archivo real
	// 	http.ServeFile(w, r, "./docs/cuotas/cheque.pdf")
	// })

	fmt.Println("Servidor iniciado en el puerto: 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
