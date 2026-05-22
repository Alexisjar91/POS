package domain

import "context"

// RoleRepository define el contrato para la persistencia de roles.
// Este repositorio maneja exclusivamente los datos del rol (id, name, description, isSystem).
// La relación many-to-many con permisos se gestiona a través de un repositorio separado.
type RoleRepository interface {
	// Crear guarda un nuevo rol en el repositorio.
	// Puede retornar ErrRolDuplicado si ya existe un rol con el mismo nombre.
	// Puede retornar ErrRepositorio ante un error interno.
	Crear(ctx context.Context, role *Role) (*Role, error)

	// ObtenerPorID busca un rol por su identificador único.
	// Retorna ErrRolNoEncontrado si no existe un rol con el id dado.
	// Puede retornar ErrRepositorio ante un error interno.
	ObtenerPorID(ctx context.Context, id string) (*Role, error)

	// ObtenerPorNombre busca un rol por su nombre.
	// Retorna ErrRolNoEncontrado si no existe un rol con el nombre dado.
	// Puede retornar ErrRepositorio ante un error interno.
	ObtenerPorNombre(ctx context.Context, name string) (*Role, error)

	// Actualizar actualiza los datos de un rol existente.
	// Retorna ErrRolNoEncontrado si el rol no existe.
	// Retorna ErrRolSistemaInmutable si se intenta modificar un rol de sistema.
	// Puede retornar ErrRolDuplicado si el nuevo nombre ya está en uso.
	// Puede retornar ErrRepositorio ante un error interno.
	Actualizar(ctx context.Context, role *Role) (*Role, error)

	// Eliminar borra un rol del repositorio.
	// Retorna ErrRolNoEncontrado si el rol no existe.
	// Retorna ErrRolSistemaInmutable si se intenta eliminar un rol de sistema.
	// Retorna ErrRolConUsuarios si el rol tiene usuarios asignados.
	// Puede retornar ErrRepositorio ante un error interno.
	Eliminar(ctx context.Context, id string) error

	// Listar retorna todos los roles del repositorio.
	// Puede retornar ErrRepositorio ante un error interno.
	Listar(ctx context.Context) ([]*Role, error)
}
