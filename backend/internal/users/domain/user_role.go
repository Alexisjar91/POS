package domain

// UserRole representa la relación many-to-many entre un usuario y un rol.
// Solo contiene referencias (IDs), no objetos completos.
type UserRole struct {
	userID string
	roleID string
}

// NuevoUserRole crea una nueva relación UserRole.
// Si userID está vacío retorna ErrUsuarioNoEncontrado.
// Si roleID está vacío retorna ErrRolNoEncontrado.
func NuevoUserRole(userID, roleID string) (*UserRole, error) {
	if userID == "" {
		return nil, ErrUsuarioNoEncontrado
	}
	if roleID == "" {
		return nil, ErrRolNoEncontrado
	}
	return &UserRole{
		userID: userID,
		roleID: roleID,
	}, nil
}

// NuevoUserRoleDesdeBD reconstruye un UserRole desde datos persistentes.
// No realiza validación: asume datos consistentes provenientes de la base de datos.
func NuevoUserRoleDesdeBD(userID, roleID string) *UserRole {
	return &UserRole{
		userID: userID,
		roleID: roleID,
	}
}

// UserID retorna el identificador del usuario asociado.
func (ur *UserRole) UserID() string {
	return ur.userID
}

// RoleID retorna el identificador del rol asociado.
func (ur *UserRole) RoleID() string {
	return ur.roleID
}
