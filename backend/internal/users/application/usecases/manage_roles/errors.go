package manageroles

import "errors"

var (
	ErrNombreRequerido       = errors.New("el nombre del rol es requerido")
	ErrRolRequerido          = errors.New("el ID del rol es requerido")
	ErrEjecutorRequerido     = errors.New("el ID del ejecutor es requerido")
	ErrRolSistemaInmutable   = errors.New("el rol de sistema no puede ser modificado") // RN-ROL-002
	ErrDescripcionRequerida  = errors.New("la descripción del rol es requerida")
)
