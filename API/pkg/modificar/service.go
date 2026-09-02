package modificar

import (
	"errors"
	"strings"
)

func (s *centralStruct) serviceCrearUsuario(nuevo_usuario *FormatoCrearUsuario) (int, error) {
	if nuevo_usuario.Nombre == "" {
		return -1, errors.New("El nombre no puede estar vacío")
	}
	if nuevo_usuario.Puestos <= 0 {
		return -1, errors.New("No se pueden tener puestos iguales o menores a 0")
	}
	if nuevo_usuario.Telefono == "" {
		nuevo_usuario.Telefono = "n"
	}

	nombre := nuevo_usuario.Nombre
	nombre = strings.TrimSpace(nombre)
	nombre = strings.Join(strings.Fields(nombre), " ")
	nuevo_usuario.Nombre = strings.ToLower(nombre)

	return s.storeCrearUsuario(nuevo_usuario)
}

func (s *centralStruct) serviceObtenerUsuarios() ([]*Usuario, error) {
	return s.storeObtenerUsuarios()
}
