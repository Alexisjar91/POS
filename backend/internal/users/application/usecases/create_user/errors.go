package createuser

import "errors"

// Errores de validación del caso de uso crear usuario.
var (
	ErrEmailRequerido    = errors.New("el email es requerido")
	ErrEmailInvalido     = errors.New("el email no tiene un formato válido")
	ErrNombreRequerido   = errors.New("el nombre completo es requerido")
	ErrPasswordRequerido = errors.New("la contraseña es requerida")
	ErrCreadorRequerido  = errors.New("el usuario creador es requerido")
	ErrEjecutorRequerido = errors.New("el usuario ejecutor es requerido")
)
