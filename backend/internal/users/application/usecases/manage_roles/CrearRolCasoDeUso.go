package manageroles

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

type CrearRolCasoDeUso struct {
	roleRepo domain.RoleRepository
	authSvc  domain.AuthorizationService
}

func NewCrearRolCasoDeUso(roleRepo domain.RoleRepository, authSvc domain.AuthorizationService) *CrearRolCasoDeUso {
	return &CrearRolCasoDeUso{roleRepo: roleRepo, authSvc: authSvc}
}

func (uc *CrearRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCrearRol) (*RespuestaCrearRol, error) {
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

	rol, err := domain.NuevoRol(cmd.Nombre, cmd.Descripcion)
	if err != nil {
		return nil, err
	}

	rolCreado, err := uc.roleRepo.Crear(ctx, rol)
	if err != nil {
		return nil, err
	}

	return &RespuestaCrearRol{
		ID:          rolCreado.ID(),
		Nombre:      rolCreado.Name(),
		Descripcion: cmd.Descripcion,
		IsSystem:    rolCreado.IsSystem(),
	}, nil
}
