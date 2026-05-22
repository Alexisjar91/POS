package domain

// Permission representa un permiso del sistema.
// Los campos no exportados garantizan inmutabilidad desde fuera del paquete.
type Permission struct {
	id          string
	code        string
	description string
	module      string
}

// NuevoPermiso crea un nuevo Permission validando los campos requeridos.
// code y module son obligatorios. Si code está vacío retorna ErrCodeRequerido.
// Si module está vacío retorna ErrModuleRequerido.
func NuevoPermiso(code, description, module string) (*Permission, error) {
	if code == "" {
		return nil, ErrCodeRequerido
	}
	if module == "" {
		return nil, ErrModuleRequerido
	}
	return &Permission{
		code:        code,
		description: description,
		module:      module,
	}, nil
}

// NuevoPermisoDesdeBD reconstruye un Permission desde datos persistentes.
// No realiza validación: asume datos consistentes provenientes de la base de datos.
func NuevoPermisoDesdeBD(id, code, description, module string) *Permission {
	return &Permission{
		id:          id,
		code:        code,
		description: description,
		module:      module,
	}
}

// Code retorna el código del permiso (ej: "create_user").
func (p *Permission) Code() string {
	return p.code
}

// Module retorna el módulo al que pertenece el permiso (ej: "users").
func (p *Permission) Module() string {
	return p.module
}

// Description retorna la descripción del permiso.
func (p *Permission) Description() string {
	return p.description
}

// BelongsToModule verifica si el permiso pertenece al módulo indicado.
// Retorna true si p.module == moduleName, false en caso contrario.
func (p *Permission) BelongsToModule(moduleName string) bool {
	return p.module == moduleName
}
