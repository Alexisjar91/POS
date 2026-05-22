package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"github.com/Alexisjar91/POS/pkg/permissions"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// RunSeed ejecuta la siembra inicial de datos del módulo usuarios:
//  1. Inserta todos los permisos registrados en el sistema (idempotente)
//  2. Crea los roles de sistema: OWNER, ADMIN (si no existen)
//  3. Asigna todos los permisos al rol ADMIN
//
// No crea el usuario OWNER — ese se solicita en el primer inicio del sistema.
func RunSeed(db *gorm.DB) error {
	ctx := context.Background()

	// 1. Sembrar permisos desde el registro global
	for _, p := range permissions.All() {
		var count int64
		if err := db.WithContext(ctx).Model(&models.PermissionModel{}).
			Where("code = ?", p.Code).Count(&count).Error; err != nil {
			return fmt.Errorf("verificar permiso %s: %w", p.Code, err)
		}
		if count > 0 {
			continue
		}

		perm, err := domain.NuevoPermiso(p.Code, p.Description, p.Module)
		if err != nil {
			return fmt.Errorf("crear permiso %s: %w", p.Code, err)
		}

		model := &models.PermissionModel{
			ID:          ulid.Make().String(),
			Code:        perm.Code(),
			Description: p.Description,
			Module:      perm.Module(),
		}
		if err := db.WithContext(ctx).Create(model).Error; err != nil {
			return fmt.Errorf("insertar permiso %s: %w", p.Code, err)
		}
	}

	// 2. Crear roles de sistema
	systemRoles := []struct {
		Name        string
		Description string
	}{
		{"OWNER", "Propietario del sistema — tiene todos los permisos implícitamente"},
		{"ADMIN", "Administrador del sistema — tiene todos los permisos explícitamente"},
	}

	var adminRoleID string

	for _, sr := range systemRoles {
		var roleModel models.RoleModel
		err := db.WithContext(ctx).Where("name = ?", sr.Name).First(&roleModel).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role := &models.RoleModel{
				ID:          ulid.Make().String(),
				Name:        sr.Name,
				Description: sr.Description,
				IsSystem:    true,
			}
			if err := db.WithContext(ctx).Create(role).Error; err != nil {
				return fmt.Errorf("crear rol %s: %w", sr.Name, err)
			}
			if sr.Name == "ADMIN" {
				adminRoleID = role.ID
			}
		} else if err != nil {
			return fmt.Errorf("buscar rol %s: %w", sr.Name, err)
		} else {
			if sr.Name == "ADMIN" {
				adminRoleID = roleModel.ID
			}
		}
	}

	// 3. Asignar todos los permisos al rol ADMIN
	if adminRoleID != "" {
		var permissionModels []models.PermissionModel
		if err := db.WithContext(ctx).Find(&permissionModels).Error; err != nil {
			return fmt.Errorf("listar permisos: %w", err)
		}

		for _, pm := range permissionModels {
			rp := &models.RolePermissionModel{
				RoleID:       adminRoleID,
				PermissionID: pm.ID,
			}
			if err := db.WithContext(ctx).Create(rp).Error; err != nil {
				if !isUniqueViolation(err) {
					return fmt.Errorf("asignar permiso %s al rol ADMIN: %w", pm.Code, err)
				}
				// ya asignado → idempotente
			}
		}
	}

	return nil
}
