package assignrole

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// AsignarRolCasoDeUso orquesta la asignación de un rol a un usuario.
type AsignarRolCasoDeUso struct {
	userRepo     domain.UserRepository
	roleRepo     domain.RoleRepository
	userRoleRepo domain.UserRoleRepository
	authSvc      domain.AuthorizationService
}

// NewAsignarRolCasoDeUso crea una nueva instancia del caso de uso.
func NewAsignarRolCasoDeUso(
	userRepo domain.UserRepository,
	roleRepo domain.RoleRepository,
	userRoleRepo domain.UserRoleRepository,
	authSvc domain.AuthorizationService,
) *AsignarRolCasoDeUso {
	return &AsignarRolCasoDeUso{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		authSvc:      authSvc,
	}
}

// Ejecutar ejecuta el caso de uso asignar rol.
// Realiza: validación → autorización → verificar usuario → verificar rol (y RN-ROL-005) → persistir → respuesta.
func (uc *AsignarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoAsignarRol) (*RespuestaAsignarRol, error) {
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
		return nil, ErrRolOWNERNoAsignable
	}

	// 5. Crear la relación UserRole
	userRole, err := domain.NuevoUserRole(cmd.UserID, cmd.RoleID)
	if err != nil {
		return nil, err
	}

	// 6. Persistir la asignación
	if err := uc.userRoleRepo.Asignar(ctx, userRole.UserID(), userRole.RoleID()); err != nil {
		return nil, err
	}

	// 7. Armar respuesta
	return &RespuestaAsignarRol{
		UserID: cmd.UserID,
		RoleID: cmd.RoleID,
	}, nil
}
