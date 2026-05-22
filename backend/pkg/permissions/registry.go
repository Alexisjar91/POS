// Package permissions provee un registro global de permisos del sistema.
// Cada módulo registra sus propios permisos al iniciar (init() o bootstrap).
// El seed usa permissions.All() para poblar la base de datos.
//
// Uso en un módulo:
//
//	func init() {
//	    permissions.Register(
//	        permissions.Permission{Code: domain.CreateUser, Description: "Crear cuentas de usuario", Module: "users"},
//	    )
//	}
package permissions

// Permission representa un permiso registrable en el sistema.
type Permission struct {
	Code        string
	Description string
	Module      string
}

var registry []Permission

// Register agrega uno o más permisos al registro global.
// Es seguro llamarlo múltiples veces desde distintos módulos.
// No valida duplicados: si dos módulos registran el mismo Code,
// el seed insertará ambos (el segundo fallará por unique constraint).
func Register(perms ...Permission) {
	registry = append(registry, perms...)
}

// All retorna una copia de todos los permisos registrados.
func All() []Permission {
	result := make([]Permission, len(registry))
	copy(result, registry)
	return result
}
