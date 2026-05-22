package enableuser

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// ActivarUsuarioCasoDeUso orquesta la activación de un usuario.
type ActivarUsuarioCasoDeUso struct {
	userRepo domain.UserRepository
	authSvc  domain.AuthorizationService
}

// NewActivarUsuarioCasoDeUso crea una nueva instancia del caso de uso.
func NewActivarUsuarioCasoDeUso(userRepo domain.UserRepository, authSvc domain.AuthorizationService) *ActivarUsuarioCasoDeUso {
	return &ActivarUsuarioCasoDeUso{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

// Ejecutar ejecuta el caso de uso activar usuario.
// Realiza: validación → autorización → carga del usuario → activar → persistir → respuesta.
func (uc *ActivarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoActivarUsuario) (*RespuestaActivarUsuario, error) {
	// 1. Validar comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 2. Autorizar: el ejecutor debe tener permiso EnableUser
	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.EnableUser)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	// 3. Cargar usuario objetivo
	target, err := uc.userRepo.ObtenerPorID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// 4. Activar entidad (valida que no esté ya activo)
	if err := target.Enable(); err != nil {
		return nil, err
	}

	// 5. Persistir
	usuarioActualizado, err := uc.userRepo.Actualizar(ctx, target)
	if err != nil {
		return nil, err
	}

	// 6. Responder
	return &RespuestaActivarUsuario{
		ID:     usuarioActualizado.ID(),
		Active: usuarioActualizado.IsActive(),
	}, nil
}
