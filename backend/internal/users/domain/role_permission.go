package domain

// RolePermission representa la relación many-to-many entre un rol y un permiso.
// Solo contiene referencias (IDs), no objetos completos.
type RolePermission struct {
	roleID       string
	permissionID string
}

// NuevoRolePermission crea una nueva relación RolePermission.
// Si roleID está vacío retorna ErrRolNoEncontrado.
// Si permissionID está vacío retorna ErrPermisoNoEncontrado.
func NuevoRolePermission(roleID, permissionID string) (*RolePermission, error) {
	if roleID == "" {
		return nil, ErrRolNoEncontrado
	}
	if permissionID == "" {
		return nil, ErrPermisoNoEncontrado
	}
	return &RolePermission{
		roleID:       roleID,
		permissionID: permissionID,
	}, nil
}

// NuevoRolePermissionDesdeBD reconstruye un RolePermission desde datos persistentes.
// No realiza validación: asume datos consistentes provenientes de la base de datos.
func NuevoRolePermissionDesdeBD(roleID, permissionID string) *RolePermission {
	return &RolePermission{
		roleID:       roleID,
		permissionID: permissionID,
	}
}

// RoleID retorna el identificador del rol asociado.
func (rp *RolePermission) RoleID() string {
	return rp.roleID
}

// PermissionID retorna el identificador del permiso asociado.
func (rp *RolePermission) PermissionID() string {
	return rp.permissionID
}
