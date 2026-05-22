package postgres

import (
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres/models"
	"gorm.io/gorm"
)

// RunMigrations ejecuta AutoMigrate para todos los modelos del módulo usuarios.
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.UserModel{},
		&models.RoleModel{},
		&models.PermissionModel{},
		&models.UserRoleModel{},
		&models.RolePermissionModel{},
	)
}
