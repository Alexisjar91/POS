package changepassword

import (
	"context"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// CambiarContrasenaCasoDeUso orquesta el cambio de contraseña de un usuario.
type CambiarContrasenaCasoDeUso struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
}

// NewCambiarContrasenaCasoDeUso crea una nueva instancia del caso de uso.
func NewCambiarContrasenaCasoDeUso(
	userRepo domain.UserRepository,
	passwordHasher domain.PasswordHasher,
) *CambiarContrasenaCasoDeUso {
	return &CambiarContrasenaCasoDeUso{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// Ejecutar ejecuta el caso de uso cambiar contraseña.
// Realiza: validación → cargar usuario → verificar contraseña actual → hashear nueva → persistir → respuesta.
func (uc *CambiarContrasenaCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCambiarContrasena) (*RespuestaCambiarContrasena, error) {
	// 1. Validar comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 2. Cargar usuario
	user, err := uc.userRepo.ObtenerPorID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// 3. Verificar contraseña actual
	if err := uc.passwordHasher.Compare(cmd.CurrentPassword, user.PasswordHash()); err != nil {
		return nil, ErrContrasenaActualIncorrecta
	}

	// 4. Hashear nueva contraseña
	newHash, err := uc.passwordHasher.Hash(cmd.NewPassword)
	if err != nil {
		return nil, err
	}

	// 5. Actualizar en la entidad
	user.SetPasswordHash(newHash)

	// 6. Persistir
	actualizado, err := uc.userRepo.Actualizar(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RespuestaCambiarContrasena{
		ID:     actualizado.ID(),
		Active: actualizado.IsActive(),
	}, nil
}
