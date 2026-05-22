package createuser

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// CrearUsuarioCasoDeUso orquesta la creación de un nuevo usuario.
type CrearUsuarioCasoDeUso struct {
	userRepo domain.UserRepository
	authSvc  domain.AuthorizationService
}

// NewCrearUsuarioCasoDeUso crea una nueva instancia del caso de uso.
func NewCrearUsuarioCasoDeUso(userRepo domain.UserRepository, authSvc domain.AuthorizationService) *CrearUsuarioCasoDeUso {
	return &CrearUsuarioCasoDeUso{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

// Ejecutar ejecuta el caso de uso crear usuario.
// Realiza: validación → autorización → crear entidad → persistir → respuesta.
func (uc *CrearUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCrearUsuario) (*RespuestaCrearUsuario, error) {
	// 1. Validar comando (formato de email, campos requeridos)
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 2. Autorizar: el ejecutor debe tener permiso CreateUser
	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.CreateUser)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	// 3. Construir entidad de dominio
	user, err := domain.NuevoUsuario(cmd.Email, cmd.FullName, cmd.Password, cmd.CreatedBy)
	if err != nil {
		return nil, err
	}

	// 4. Persistir (el repositorio asigna ID y CreatedAt)
	userCreado, err := uc.userRepo.Crear(ctx, user)
	if err != nil {
		return nil, err
	}

	// 5. Armar respuesta
	return &RespuestaCrearUsuario{
		ID:        userCreado.ID(),
		Email:     userCreado.Email(),
		FullName:  userCreado.FullName(),
		Active:    userCreado.IsActive(),
		CreatedAt: userCreado.CreatedAt(),
	}, nil
}
