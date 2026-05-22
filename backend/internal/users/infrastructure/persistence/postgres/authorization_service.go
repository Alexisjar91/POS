package postgres

import (
	"context"
	"errors"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"gorm.io/gorm"
)

// authorizationService implementa domain.AuthorizationService
// consultando permisos mediante SQL raw sobre JOINs múltiples.
type authorizationService struct {
	db *gorm.DB
}

// NewAuthorizationService crea un AuthorizationService que consulta la base de datos
// para verificar permisos y roles.
func NewAuthorizationService(db *gorm.DB) domain.AuthorizationService {
	return &authorizationService{db: db}
}

// VerificarPermiso verifica si un usuario tiene un permiso específico.
//
//  1. Si el usuario tiene rol OWNER → true (early return).
//  2. Si algún rol del usuario tiene el permiso solicitado → true.
//  3. En cualquier otro caso → false.
//
// Si el userID no existe en user_roles, COUNT devuelve 0 y retorna false, nil.
func (s *authorizationService) VerificarPermiso(ctx context.Context, userID string, permissionCode string) (bool, error) {
	// 1. Verificar si el usuario es OWNER — los OWNER tienen todos los permisos.
	isOwner, err := s.EsUsuarioOWNER(ctx, userID)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}

	// 2. Verificar si algún rol del usuario tiene el permiso solicitado.
	var count int64
	err = s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = ? AND p.code = ?
	`, userID, permissionCode).Scan(&count).Error
	if err != nil {
		return false, errors.Join(domain.ErrRepositorio, err)
	}

	return count > 0, nil
}

// EsUsuarioOWNER verifica si un usuario tiene el rol OWNER.
func (s *authorizationService) EsUsuarioOWNER(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.name = ?
	`, userID, "OWNER").Scan(&count).Error
	if err != nil {
		return false, errors.Join(domain.ErrRepositorio, err)
	}
	return count > 0, nil
}
