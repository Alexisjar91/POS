package changepassword

// ComandoCambiarContrasena contiene los datos necesarios para cambiar la contraseña de un usuario.
type ComandoCambiarContrasena struct {
	UserID          string // usuario autenticado
	CurrentPassword string // contraseña actual en texto plano
	NewPassword     string // nueva contraseña en texto plano
}

// Validar ejecuta las validaciones de aplicación sobre el comando.
// Retorna el primer error encontrado.
func (cmd *ComandoCambiarContrasena) Validar() error {
	if cmd.UserID == "" {
		return ErrUsuarioRequerido
	}
	if cmd.CurrentPassword == "" {
		return ErrContrasenaActualRequerida
	}
	if cmd.NewPassword == "" {
		return ErrNuevaContrasenaRequerida
	}
	if cmd.CurrentPassword == cmd.NewPassword {
		return ErrContrasenaIgual
	}
	return nil
}
