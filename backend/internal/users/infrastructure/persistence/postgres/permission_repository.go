package postgres

import (
	"context"
	"errors"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// PermissionRepository define el contrato para la persistencia de permisos.
// Es una interfaz interna del paquete postgres (no está en dominio).
// La usa AuthorizationService para verificar permisos y Seed para poblar datos iniciales.
type PermissionRepository interface {
	Crear(ctx context.Context, permission *domain.Permission) (*domain.Permission, error)
	ObtenerPorID(ctx context.Context, id string) (*domain.Permission, error)
	ObtenerPorCode(ctx context.Context, code string) (*domain.Permission, error)
	ListarPorModulo(ctx context.Context, module string) ([]*domain.Permission, error)
	ListarTodos(ctx context.Context) ([]*domain.Permission, error)
	ExistePorCode(ctx context.Context, code string) (bool, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Crear(ctx context.Context, permission *domain.Permission) (*domain.Permission, error) {
	model := &models.PermissionModel{
		ID:          ulid.Make().String(),
		Code:        permission.Code(),
		Description: permission.Description(),
		Module:      permission.Module(),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrPermisoDuplicado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	return model.ToDomain(), nil
}

func (r *permissionRepository) ObtenerPorID(ctx context.Context, id string) (*domain.Permission, error) {
	var model models.PermissionModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPermisoNoEncontrado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}
	return model.ToDomain(), nil
}

func (r *permissionRepository) ObtenerPorCode(ctx context.Context, code string) (*domain.Permission, error) {
	var model models.PermissionModel
	if err := r.db.WithContext(ctx).First(&model, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPermisoNoEncontrado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}
	return model.ToDomain(), nil
}

func (r *permissionRepository) ListarPorModulo(ctx context.Context, module string) ([]*domain.Permission, error) {
	var modelsList []models.PermissionModel
	if err := r.db.WithContext(ctx).Where("module = ?", module).Find(&modelsList).Error; err != nil {
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	permissions := make([]*domain.Permission, len(modelsList))
	for i := range modelsList {
		permissions[i] = modelsList[i].ToDomain()
	}
	return permissions, nil
}

func (r *permissionRepository) ListarTodos(ctx context.Context) ([]*domain.Permission, error) {
	var modelsList []models.PermissionModel
	if err := r.db.WithContext(ctx).Find(&modelsList).Error; err != nil {
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	permissions := make([]*domain.Permission, len(modelsList))
	for i := range modelsList {
		permissions[i] = modelsList[i].ToDomain()
	}
	return permissions, nil
}

func (r *permissionRepository) ExistePorCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.PermissionModel{}).
		Where("code = ?", code).
		Count(&count).Error; err != nil {
		return false, errors.Join(domain.ErrRepositorio, err)
	}
	return count > 0, nil
}
