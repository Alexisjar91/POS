package models

import (
	"github.com/Alexisjar91/POS/internal/users/domain"
)

// RolePermissionModel representa la tabla pivote role_permissions.
type RolePermissionModel struct {
	RoleID       string `gorm:"column:role_id;type:varchar(26);primaryKey"`
	PermissionID string `gorm:"column:permission_id;type:varchar(26);primaryKey"`
}

func (RolePermissionModel) TableName() string { return "role_permissions" }

func (m *RolePermissionModel) ToDomain() *domain.RolePermission {
	return domain.NuevoRolePermissionDesdeBD(m.RoleID, m.PermissionID)
}
