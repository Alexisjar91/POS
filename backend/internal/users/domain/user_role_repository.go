package domain

import "context"

// UserRoleRepository define el contrato para la gestión de la relación
// many-to-many entre usuarios y roles.
type UserRoleRepository interface {
	// Asignar asigna un rol a un usuario.
	// Retorna ErrUsuarioNoEncontrado si el usuario no existe.
	// Retorna ErrRolNoEncontrado si el rol no existe.
	// Puede retornar ErrRepositorio ante un error interno.
	Asignar(ctx context.Context, userID string, roleID string) error

	// Remover remueve un rol de un usuario.
	// Retorna ErrUsuarioNoEncontrado si el usuario no existe.
	// Retorna ErrRolNoEncontrado si el rol no existe.
	// Puede retornar ErrRepositorio ante un error interno.
	Remover(ctx context.Context, userID string, roleID string) error
}
