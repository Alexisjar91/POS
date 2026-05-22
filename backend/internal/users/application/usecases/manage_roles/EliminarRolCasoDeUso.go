package manageroles

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

type EliminarRolCasoDeUso struct {
	roleRepo domain.RoleRepository
	authSvc  domain.AuthorizationService
}

func NewEliminarRolCasoDeUso(roleRepo domain.RoleRepository, authSvc domain.AuthorizationService) *EliminarRolCasoDeUso {
	return &EliminarRolCasoDeUso{roleRepo: roleRepo, authSvc: authSvc}
}

func (uc *EliminarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoEliminarRol) (*RespuestaEliminarRol, error) {
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.ManageRoles)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	// Cargar rol para verificar RN-ROL-002
	rol, err := uc.roleRepo.ObtenerPorID(ctx, cmd.RoleID)
	if err != nil {
		return nil, err
	}

	// RN-ROL-002: roles de sistema no se eliminan
	if rol.IsSystem() {
		return nil, ErrRolSistemaInmutable
	}

	if err := uc.roleRepo.Eliminar(ctx, cmd.RoleID); err != nil {
		return nil, err
	}

	return &RespuestaEliminarRol{
		ID:      cmd.RoleID,
		Exitoso: true,
	}, nil
}
