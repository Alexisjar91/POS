package manageroles

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

type ActualizarRolCasoDeUso struct {
	roleRepo domain.RoleRepository
	authSvc  domain.AuthorizationService
}

func NewActualizarRolCasoDeUso(roleRepo domain.RoleRepository, authSvc domain.AuthorizationService) *ActualizarRolCasoDeUso {
	return &ActualizarRolCasoDeUso{roleRepo: roleRepo, authSvc: authSvc}
}

func (uc *ActualizarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoActualizarRol) (*RespuestaActualizarRol, error) {
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

	// Cargar rol existente
	rolExistente, err := uc.roleRepo.ObtenerPorID(ctx, cmd.RoleID)
	if err != nil {
		return nil, err
	}

	// RN-ROL-002: roles de sistema no se renombran
	if rolExistente.IsSystem() {
		return nil, ErrRolSistemaInmutable
	}

	// Reconstruir con nuevos datos
	rolActualizado := domain.NuevoRolDesdeBD(rolExistente.ID(), cmd.Nombre, cmd.Descripcion, false)

	resultado, err := uc.roleRepo.Actualizar(ctx, rolActualizado)
	if err != nil {
		return nil, err
	}

	return &RespuestaActualizarRol{
		ID:          resultado.ID(),
		Nombre:      resultado.Name(),
		Descripcion: cmd.Descripcion,
		IsSystem:    resultado.IsSystem(),
	}, nil
}
