package enableuser

// ComandoActivarUsuario contiene los datos necesarios para activar un usuario.
type ComandoActivarUsuario struct {
	UserID     string // ID del usuario a activar
	EjecutorID string // ID del usuario autenticado
}

// Validar ejecuta las validaciones de aplicación sobre el comando.
// Retorna el primer error encontrado.
func (cmd *ComandoActivarUsuario) Validar() error {
	if cmd.UserID == "" {
		return ErrUsuarioRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
