package domain

// Role representa un rol del sistema.
// Los campos no exportados garantizan inmutabilidad desde fuera del paquete.
// No contiene colecciones de permisos; la relación many-to-many se maneja
// a través de la entidad RolePermission.
type Role struct {
	id          string
	name        string
	description string
	isSystem    bool
}

// NuevoRol crea un nuevo Role con isSystem=false.
// Si name está vacío retorna ErrNombreRolRequerido.
func NuevoRol(name, description string) (*Role, error) {
	if name == "" {
		return nil, ErrNombreRolRequerido
	}
	return &Role{
		name:        name,
		description: description,
		isSystem:    false,
	}, nil
}

// NuevoRolSistema crea un nuevo Role con isSystem=true.
// Si name está vacío retorna ErrNombreRolRequerido.
func NuevoRolSistema(name, description string) (*Role, error) {
	if name == "" {
		return nil, ErrNombreRolRequerido
	}
	return &Role{
		name:        name,
		description: description,
		isSystem:    true,
	}, nil
}

// NuevoRolDesdeBD reconstruye un Role desde datos persistentes.
// No realiza validación: asume datos consistentes provenientes de la base de datos.
func NuevoRolDesdeBD(id, name, description string, isSystem bool) *Role {
	return &Role{
		id:          id,
		name:        name,
		description: description,
		isSystem:    isSystem,
	}
}

// IsSystem retorna true si el rol es de sistema.
func (r *Role) IsSystem() bool {
	return r.isSystem
}

// IsOwner retorna true si el rol es OWNER.
func (r *Role) IsOwner() bool {
	return r.name == "OWNER"
}

// IsAdmin retorna true si el rol es ADMIN.
func (r *Role) IsAdmin() bool {
	return r.name == "ADMIN"
}

// ID retorna el identificador único del rol.
func (r *Role) ID() string {
	return r.id
}

// Name retorna el nombre del rol.
func (r *Role) Name() string {
	return r.name
}

// Description retorna la descripción del rol.
func (r *Role) Description() string {
	return r.description
}
