package models

import (
	"github.com/Alexisjar91/POS/internal/users/domain"
)

// RoleModel representa la tabla roles en la base de datos.
type RoleModel struct {
	ID          string `gorm:"type:varchar(26);primaryKey"`
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	IsSystem    bool   `gorm:"column:is_system;not null;default:false"`
}

func (RoleModel) TableName() string { return "roles" }

func (m *RoleModel) ToDomain() *domain.Role {
	return domain.NuevoRolDesdeBD(m.ID, m.Name, m.Description, m.IsSystem)
}
