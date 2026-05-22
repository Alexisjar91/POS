package assignrole

// ComandoAsignarRol contiene los datos necesarios para asignar un rol a un usuario.
type ComandoAsignarRol struct {
	UserID     string // usuario que recibe el rol
	RoleID     string // rol a asignar
	EjecutorID string // usuario autenticado (para autorización)
}

// Validar ejecuta las validaciones de aplicación sobre el comando.
// Retorna el primer error encontrado.
func (cmd *ComandoAsignarRol) Validar() error {
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
