package disableuser

import "errors"

var (
	ErrUsuarioRequerido  = errors.New("el ID del usuario a desactivar es requerido")
	ErrEjecutorRequerido = errors.New("el ID del ejecutor es requerido")
	ErrAutoDesactivacion = errors.New("no puedes desactivarte a ti mismo")        // RN-USR-008
	ErrOWNERInmune       = errors.New("el usuario OWNER no puede ser desactivado") // RN-USR-009
)
