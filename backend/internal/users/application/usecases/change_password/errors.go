package changepassword

import "errors"

var (
	ErrUsuarioRequerido            = errors.New("el ID del usuario es requerido")
	ErrContrasenaActualRequerida   = errors.New("la contraseña actual es requerida")
	ErrNuevaContrasenaRequerida    = errors.New("la nueva contraseña es requerida")
	ErrContrasenaActualIncorrecta  = errors.New("la contraseña actual es incorrecta")
	ErrContrasenaIgual             = errors.New("la nueva contraseña debe ser diferente a la actual")
)
