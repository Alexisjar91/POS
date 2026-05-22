package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Crear(ctx context.Context, user *domain.User) (*domain.User, error) {
	id := user.ID()
	if id == "" {
		id = ulid.Make().String()
	}

	model := &models.UserModel{
		ID:           id,
		Email:        user.Email(),
		PasswordHash: user.PasswordHash(),
		FullName:     user.FullName(),
		Active:       user.IsActive(),
		CreatedBy:    user.CreatedBy(),
		CreatedAt:    parseTimeOrNow(user.CreatedAt()),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrEmailDuplicado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	return model.ToDomain(), nil
}

func (r *userRepository) ObtenerPorID(ctx context.Context, id string) (*domain.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUsuarioNoEncontrado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}
	return model.ToDomain(), nil
}

func (r *userRepository) ObtenerPorEmail(ctx context.Context, email string) (*domain.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUsuarioNoEncontrado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}
	return model.ToDomain(), nil
}

func (r *userRepository) Actualizar(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Verify the user exists before updating
	if _, err := r.ObtenerPorID(ctx, user.ID()); err != nil {
		return nil, err
	}

	model := &models.UserModel{
		ID:           user.ID(),
		Email:        user.Email(),
		PasswordHash: user.PasswordHash(),
		FullName:     user.FullName(),
		Active:       user.IsActive(),
		CreatedBy:    user.CreatedBy(),
		CreatedAt:    parseTimeOrNow(user.CreatedAt()),
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrEmailDuplicado
		}
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	return model.ToDomain(), nil
}

func (r *userRepository) Listar(ctx context.Context, spec especificacion.Especificacion, pag paginacion.Paginacion) ([]*domain.User, error) {
	tx := r.db.WithContext(ctx).Model(&models.UserModel{})

	// Aplicar filtros desde la especificación
	for _, criterio := range spec.Filtros {
		// Validar que el campo esté permitido
		if !domain.ColumnasPermitidas[criterio.Campo] {
			continue
		}

		// Mapear al nombre de columna en la BD
		columnaBD := domain.MapeoColumnas[criterio.Campo]

		// Aplicar filtro según el operador
		switch criterio.Operador {
		case "=":
			// Para email, hacer búsqueda case-insensitive
			if criterio.Campo == "email" {
				tx = tx.Where("LOWER("+columnaBD+") = LOWER(?)", criterio.Valor)
			} else {
				tx = tx.Where(columnaBD+" = ?", criterio.Valor)
			}
		case "!=":
			tx = tx.Where(columnaBD+" != ?", criterio.Valor)
		case "LIKE":
			// Para LIKE, siempre hacer case-insensitive
			tx = tx.Where("LOWER("+columnaBD+") LIKE LOWER(?)", criterio.Valor)
		}
	}

	// Aplicar ordenaciones desde paginación
	if len(pag.Ordenaciones) > 0 {
		for _, ord := range pag.Ordenaciones {
			// Validar que el campo esté permitido
			if !domain.ColumnasPermitidas[ord.Campo] {
				continue
			}

			columnaBD := domain.MapeoColumnas[ord.Campo]
			dirección := "ASC"
			if ord.Tipo == paginacion.DESC {
				dirección = "DESC"
			}
			tx = tx.Order(columnaBD + " " + dirección)
		}
	} else {
		// Ordenación por defecto
		tx = tx.Order("created_at DESC")
	}

	// Aplicar offset y limit desde paginación
	offset := pag.Offset()
	limit := pag.Limit()
	tx = tx.Offset(offset).Limit(limit)

	// Ejecutar consulta
	var modelsList []models.UserModel
	if err := tx.Find(&modelsList).Error; err != nil {
		return nil, errors.Join(domain.ErrRepositorio, err)
	}

	// Convertir modelos a dominio
	users := make([]*domain.User, len(modelsList))
	for i := range modelsList {
		users[i] = modelsList[i].ToDomain()
	}

	return users, nil
}

func (r *userRepository) ExistePorEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, errors.Join(domain.ErrRepositorio, err)
	}
	return count > 0, nil
}

// helpers

func parseTimeOrNow(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
