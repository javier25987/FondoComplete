package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func RegistrarLog(mensaje string) error {
	fechaNow := time.Now().Format("2006-01-02")

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error obteniendo ruta del ejecutable: %w", err)
	}

	logDir := filepath.Join(filepath.Dir(exePath), "logs")

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("error creando carpeta logs: %w", err)
	}

	nombreArchivo := fechaNow + ".txt"
	rutaArchivo := filepath.Join(logDir, nombreArchivo)

	// O_WRONLY = solo escritura, O_CREATE = crear si no existe, O_APPEND = agregar al final
	archivo, err := os.OpenFile(rutaArchivo, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de log: %w", err)
	}
	defer archivo.Close()

	logMensage := fmt.Sprintf("[%s] %s\n", time.Now().Format("15:04:05"), mensaje)
	if _, err := archivo.WriteString(logMensage); err != nil {
		return fmt.Errorf("error escribiendo en el log: %w", err)
	}

	return nil
}
