package postgres

import (
	"context"
	"errors"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Crear(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	id := role.ID()
	if id == "" {
		id = ulid.Make().String()
	}

	model := &models.RoleModel{
		ID:          id,
		Name:        role.Name(),
		Description: role.Description(),
		IsSystem:    role.IsSystem(),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrRolDuplicado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	return model.ToDomain(), nil
}

func (r *roleRepository) ObtenerPorID(ctx context.Context, id string) (*domain.Role, error) {
	var model models.RoleModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRolNoEncontrado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}
	return model.ToDomain(), nil
}

func (r *roleRepository) ObtenerPorNombre(ctx context.Context, name string) (*domain.Role, error) {
	var model models.RoleModel
	if err := r.db.WithContext(ctx).First(&model, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRolNoEncontrado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}
	return model.ToDomain(), nil
}

func (r *roleRepository) Actualizar(ctx context.Context, role *domain.Role) (*domain.Role, error) {
	// 1. Load existing role — if not found, returns ErrRolNoEncontrado
	existing, err := r.ObtenerPorID(ctx, role.ID())
	if err != nil {
		return nil, err
	}

	// 2. System roles are immutable
	if existing.IsSystem() {
		return nil, domain.ErrRolSistemaInmutable
	}

	// Build model preserving the system flag from the database
	model := &models.RoleModel{
		ID:          role.ID(),
		Name:        role.Name(),
		Description: role.Description(),
		IsSystem:    existing.IsSystem(),
	}

	// 3. Save — unique_violation on name maps to ErrRolDuplicado
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrRolDuplicado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	return model.ToDomain(), nil
}

func (r *roleRepository) Eliminar(ctx context.Context, id string) error {
	// 1. Load existing role — if not found, returns ErrRolNoEncontrado
	existing, err := r.ObtenerPorID(ctx, id)
	if err != nil {
		return err
	}

	// 2. System roles cannot be deleted
	if existing.IsSystem() {
		return domain.ErrRolSistemaInmutable
	}

	// 3. Check if role has assigned users
	var count int64
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("role_id = ?", id).
		Count(&count).Error; err != nil {
		return errors.Join(domain.ErrRepositorio, err)
	}
	if count > 0 {
		return domain.ErrRolConUsuarios
	}

	// 4. Delete the role
	if err := r.db.WithContext(ctx).Delete(&models.RoleModel{ID: id}).Error; err != nil {
		return errors.Join(domain.ErrRepositorio, err)
	}

	return nil
}

func (r *roleRepository) Listar(ctx context.Context) ([]*domain.Role, error) {
	var modelsList []models.RoleModel
	if err := r.db.WithContext(ctx).Find(&modelsList).Error; err != nil {
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	roles := make([]*domain.Role, len(modelsList))
	for i := range modelsList {
		roles[i] = modelsList[i].ToDomain()
	}
	return roles, nil
}
