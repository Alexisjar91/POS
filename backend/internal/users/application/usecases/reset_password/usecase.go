package resetpassword

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

type ResetearContrasenaCasoDeUso struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
	authSvc        domain.AuthorizationService
}

func NewResetearContrasenaCasoDeUso(
	userRepo domain.UserRepository,
	passwordHasher domain.PasswordHasher,
	authSvc domain.AuthorizationService,
) *ResetearContrasenaCasoDeUso {
	return &ResetearContrasenaCasoDeUso{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		authSvc:        authSvc,
	}
}

func (uc *ResetearContrasenaCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoResetearContrasena) (*RespuestaResetearContrasena, error) {
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// Autorizar (el ejecutor necesita permiso, NO el objetivo)
	autorizado, err := uc.authSvc.VerificarPermiso(ctx, cmd.EjecutorID, domain.ResetUserPassword)
	if err != nil {
		return nil, err
	}
	if !autorizado {
		return nil, domain.ErrAccesoDenegado
	}

	// Cargar usuario objetivo
	user, err := uc.userRepo.ObtenerPorID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Hashear nueva contraseña
	newHash, err := uc.passwordHasher.Hash(cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	// Actualizar en la entidad
	user.SetPasswordHash(newHash)

	// Persistir
	actualizado, err := uc.userRepo.Actualizar(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RespuestaResetearContrasena{
		ID:     actualizado.ID(),
		Active: actualizado.IsActive(),
	}, nil
}
