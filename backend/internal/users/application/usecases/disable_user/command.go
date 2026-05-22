package disableuser

// ComandoDesactivarUsuario contiene los datos necesarios para desactivar un usuario.
type ComandoDesactivarUsuario struct {
	UserID     string // ID del usuario a desactivar
	EjecutorID string // ID del usuario autenticado
}

// Validar ejecuta las validaciones de aplicación sobre el comando.
// Retorna el primer error encontrado.
func (cmd *ComandoDesactivarUsuario) Validar() error {
	if cmd.UserID == "" {
		return ErrUsuarioRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
