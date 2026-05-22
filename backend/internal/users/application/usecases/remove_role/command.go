package removerole

// ComandoRemoverRol contiene los datos necesarios para remover un rol de un usuario.
type ComandoRemoverRol struct {
	UserID     string // usuario que pierde el rol
	RoleID     string // rol a remover
	EjecutorID string // usuario autenticado (para autorización)
}

// Validar ejecuta las validaciones de aplicación sobre el comando.
// Retorna el primer error encontrado.
func (cmd *ComandoRemoverRol) Validar() error {
	if cmd.UserID == "" {
		return ErrUsuarioRequerido
	}
	if cmd.RoleID == "" {
		return ErrRolRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
