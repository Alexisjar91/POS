package getuser

import "errors"

var (
	ErrUsuarioRequerido  = errors.New("el ID del usuario es requerido")
	ErrEjecutorRequerido = errors.New("el ID del ejecutor es requerido")
)
