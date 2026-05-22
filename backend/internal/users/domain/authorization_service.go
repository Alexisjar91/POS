package domain

import "context"

// AuthorizationService define el contrato para la verificación de permisos.
// La capa de aplicación inyecta este servicio en los casos de uso para
// validar si el usuario autenticado puede ejecutar la operación solicitada.
//
// Comportamiento esperado de la implementación:
//  1. Carga los roles del usuario.
//  2. Si algún rol es OWNER → autorizado (early return).
//  3. Si no: carga los permisos de todos los roles del usuario.
//  4. Si el permiso solicitado está en el conjunto → autorizado.
//  5. Si no → denegado.
type AuthorizationService interface {
	// VerificarPermiso verifica si un usuario tiene un permiso específico.
	// Recibe el contexto, el ID del usuario y el código del permiso
	// (constante del dominio, ej: CreateUser, DisableUser).
	//
	// Retorna true si el usuario está autorizado, false si no.
	//
	// Puede retornar ErrUsuarioNoEncontrado si el userID no existe.
	// Puede retornar ErrRepositorio ante un error de infraestructura.
	VerificarPermiso(ctx context.Context, userID string, permissionCode string) (bool, error)

	// EsUsuarioOWNER verifica si un usuario tiene el rol OWNER.
	// Retorna true si el usuario tiene el rol OWNER, false en caso contrario.
	//
	// Puede retornar ErrUsuarioNoEncontrado si el userID no existe.
	// Puede retornar ErrRepositorio ante un error de infraestructura.
	EsUsuarioOWNER(ctx context.Context, userID string) (bool, error)
}
