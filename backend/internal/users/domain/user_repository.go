package domain

import (
	"context"

	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
)

// UserRepository define el contrato para la persistencia de usuarios.
type UserRepository interface {
	// Crear guarda un nuevo usuario en el repositorio.
	// Retorna ErrEmailDuplicado si el email ya existe.
	Crear(ctx context.Context, user *User) (*User, error)

	// ObtenerPorID busca un usuario por su identificador único.
	// Retorna ErrUsuarioNoEncontrado si no existe.
	ObtenerPorID(ctx context.Context, id string) (*User, error)

	// ObtenerPorEmail busca un usuario por su email.
	// Retorna ErrUsuarioNoEncontrado si no existe.
	ObtenerPorEmail(ctx context.Context, email string) (*User, error)

	// Actualizar actualiza los datos de un usuario existente.
	// Retorna ErrUsuarioNoEncontrado si no existe.
	// Retorna ErrEmailDuplicado si el nuevo email ya está en uso.
	Actualizar(ctx context.Context, user *User) (*User, error)

	// Listar retorna usuarios filtrados por especificación y paginados.
	// especificación vacía retorna todos los usuarios paginados.
	// Retorna ErrRepositorio ante un error interno.
	Listar(ctx context.Context, especificacion especificacion.Especificacion, paginacion paginacion.Paginacion) ([]*User, error)

	// ExistePorEmail verifica si existe un usuario con el email dado.
	ExistePorEmail(ctx context.Context, email string) (bool, error)
}
