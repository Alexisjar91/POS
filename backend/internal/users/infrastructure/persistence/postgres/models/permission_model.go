package models

import (
	"github.com/Alexisjar91/POS/internal/users/domain"
)

// PermissionModel representa la tabla permissions en la base de datos.
type PermissionModel struct {
	ID          string `gorm:"type:varchar(26);primaryKey"`
	Code        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	Module      string `gorm:"type:varchar(50);not null"`
}

func (PermissionModel) TableName() string { return "permissions" }

func (m *PermissionModel) ToDomain() *domain.Permission {
	return domain.NuevoPermisoDesdeBD(m.ID, m.Code, m.Description, m.Module)
}
