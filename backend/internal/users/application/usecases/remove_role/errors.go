package removerole

import "errors"

var (
	ErrUsuarioRequerido    = errors.New("el ID del usuario es requerido")
	ErrRolRequerido        = errors.New("el ID del rol es requerido")
	ErrEjecutorRequerido   = errors.New("el ID del ejecutor es requerido")
	ErrRolOWNERNoRemovible = errors.New("el rol OWNER no puede ser removido") // RN-ROL-005
)
