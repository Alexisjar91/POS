package postgres

import (
	"context"
	"errors"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) domain.UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) Asignar(ctx context.Context, userID, roleID string) error {
	// 1. Verificar que usuario existe
	var userCount int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", userID).Count(&userCount).Error; err != nil {
		return errors.Join(domain.ErrRepositorio, err)
	}
	if userCount == 0 {
		return domain.ErrUsuarioNoEncontrado
	}

	// 2. Verificar que rol existe
	var roleCount int64
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Where("id = ?", roleID).Count(&roleCount).Error; err != nil {
		return errors.Join(domain.ErrRepositorio, err)
	}
	if roleCount == 0 {
		return domain.ErrRolNoEncontrado
	}

	// 3. Insertar en user_roles (idempotente: si ya existe, no error)
	model := &models.UserRoleModel{UserID: userID, RoleID: roleID}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isUniqueViolation(err) {
			return nil // ya asignado → idempotente
		}
		return errors.Join(domain.ErrRepositorio, err)
	}
	return nil
}

func (r *userRoleRepository) Remover(ctx context.Context, userID, roleID string) error {
	// 1. Verificar que usuario existe
	var userCount int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", userID).Count(&userCount).Error; err != nil {
		return errors.Join(domain.ErrRepositorio, err)
	}
	if userCount == 0 {
		return domain.ErrUsuarioNoEncontrado
	}

	// 2. Verificar que rol existe
	var roleCount int64
	if err := r.db.WithContext(ctx).Model(&models.RoleModel{}).Where("id = ?", roleID).Count(&roleCount).Error; err != nil {
		return errors.Join(domain.ErrRepositorio, err)
	}
	if roleCount == 0 {
		return domain.ErrRolNoEncontrado
	}

	// 3. Eliminar de user_roles (idempotente: si no existe, no error)
	result := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRoleModel{})
	if result.Error != nil {
		return errors.Join(domain.ErrRepositorio, result.Error)
	}
	return nil
}
