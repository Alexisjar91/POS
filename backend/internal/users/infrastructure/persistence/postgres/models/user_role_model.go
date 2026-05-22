package models

import (
	"github.com/Alexisjar91/POS/internal/users/domain"
)

// UserRoleModel representa la tabla pivote user_roles.
type UserRoleModel struct {
	UserID string `gorm:"column:user_id;type:varchar(26);primaryKey"`
	RoleID string `gorm:"column:role_id;type:varchar(26);primaryKey"`
}

func (UserRoleModel) TableName() string { return "user_roles" }

func (m *UserRoleModel) ToDomain() *domain.UserRole {
	return domain.NuevoUserRoleDesdeBD(m.UserID, m.RoleID)
}
