package getuser

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

type ObtenerUsuarioCasoDeUso struct {
	userRepo domain.UserRepository
	authSvc  domain.AuthorizationService
}

func NewObtenerUsuarioCasoDeUso(userRepo domain.UserRepository, authSvc domain.AuthorizationService) *ObtenerUsuarioCasoDeUso {
	return &ObtenerUsuarioCasoDeUso{userRepo: userRepo, authSvc: authSvc}
}

func (uc *ObtenerUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoObtenerUsuario) (*RespuestaObtenerUsuario, error) {
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.ViewUsers)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	user, err := uc.userRepo.ObtenerPorID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &RespuestaObtenerUsuario{
		ID:        user.ID(),
		Email:     user.Email(),
		FullName:  user.FullName(),
		Active:    user.IsActive(),
		CreatedAt: user.CreatedAt(),
	}, nil
}
