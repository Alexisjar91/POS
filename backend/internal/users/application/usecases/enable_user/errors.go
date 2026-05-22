package enableuser

import "errors"

var (
	ErrUsuarioRequerido  = errors.New("el ID del usuario a activar es requerido")
	ErrEjecutorRequerido = errors.New("el ID del ejecutor es requerido")
)
