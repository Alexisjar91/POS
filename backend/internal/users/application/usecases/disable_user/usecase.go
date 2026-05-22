package disableuser

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// DesactivarUsuarioCasoDeUso orquesta la desactivación de un usuario.
type DesactivarUsuarioCasoDeUso struct {
	userRepo domain.UserRepository
	authSvc  domain.AuthorizationService
}

// NewDesactivarUsuarioCasoDeUso crea una nueva instancia del caso de uso.
func NewDesactivarUsuarioCasoDeUso(userRepo domain.UserRepository, authSvc domain.AuthorizationService) *DesactivarUsuarioCasoDeUso {
	return &DesactivarUsuarioCasoDeUso{
		userRepo: userRepo,
		authSvc:  authSvc,
	}
}

// Ejecutar ejecuta el caso de uso desactivar usuario.
// Realiza: validación → autorización → carga del usuario → RN-USR-008 → RN-USR-009 → desactivar → persistir → respuesta.
func (uc *DesactivarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoDesactivarUsuario) (*RespuestaDesactivarUsuario, error) {
	// 1. Validar comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 2. Autorizar: el ejecutor debe tener permiso DisableUser
	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.DisableUser)
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

	// 4. RN-USR-008: el ejecutor no puede desactivarse a sí mismo
	if cmd.EjecutorID == cmd.UserID {
		return nil, ErrAutoDesactivacion
	}

	// 5. RN-USR-009: el usuario con rol OWNER es inmune a desactivación
	esOwner, err := uc.authSvc.EsUsuarioOWNER(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if esOwner {
		return nil, ErrOWNERInmune
	}

	// 6. Desactivar entidad (valida que no esté ya inactivo)
	if err := target.Disable(); err != nil {
		return nil, err
	}

	// 7. Persistir
	usuarioActualizado, err := uc.userRepo.Actualizar(ctx, target)
	if err != nil {
		return nil, err
	}

	// 8. Responder
	return &RespuestaDesactivarUsuario{
		ID:     usuarioActualizado.ID(),
		Active: usuarioActualizado.IsActive(),
	}, nil
}
