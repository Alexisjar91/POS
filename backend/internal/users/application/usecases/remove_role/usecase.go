package removerole

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// RemoverRolCasoDeUso orquesta la remoción de un rol de un usuario.
type RemoverRolCasoDeUso struct {
	userRepo     domain.UserRepository
	roleRepo     domain.RoleRepository
	userRoleRepo domain.UserRoleRepository
	authSvc      domain.AuthorizationService
}

// NewRemoverRolCasoDeUso crea una nueva instancia del caso de uso.
func NewRemoverRolCasoDeUso(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	userRoleRepo domain.UserRoleRepository,
	authSvc domain.AuthorizationService,
) *RemoverRolCasoDeUso {
	return &RemoverRolCasoDeUso{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		authSvc:      authSvc,
	}
}

// Ejecutar ejecuta el caso de uso remover rol.
// Realiza: validación → autorización → verificar usuario → verificar rol (y RN-ROL-005) → persistir → respuesta.
func (uc *RemoverRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoRemoverRol) (*RespuestaRemoverRol, error) {
	// 1. Validar comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 2. Autorizar: el ejecutor debe tener permiso AssignRole
	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.AssignRole)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	// 3. Verificar que el usuario existe
	if _, err := uc.userRepo.ObtenerPorID(ctx, cmd.UserID); err != nil {
		return nil, err
	}

	// 4. Verificar que el rol existe y no es OWNER (RN-ROL-005)
	rol, err := uc.roleRepo.ObtenerPorID(ctx, cmd.RoleID)
	if err != nil {
		return nil, err
	}
	if rol.IsOwner() {
		return nil, ErrRolOWNERNoRemovible
	}

	// 5. Remover la relación UserRole
	if err := uc.userRoleRepo.Remover(ctx, cmd.UserID, cmd.RoleID); err != nil {
		return nil, err
	}

	// 6. Armar respuesta
	return &RespuestaRemoverRol{
		UserID: cmd.UserID,
		RoleID: cmd.RoleID,
	}, nil
}
