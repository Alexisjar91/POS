package models

import (
	"time"

	"github.com/Alexisjar91/POS/internal/users/domain"
)

// UserModel representa la tabla users en la base de datos.
type UserModel struct {
	ID           string    `gorm:"type:varchar(26);primaryKey"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null"`
	FullName     string    `gorm:"column:full_name;type:varchar(255);not null"`
	Active       bool      `gorm:"not null;default:true"`
	CreatedBy    *string   `gorm:"column:created_by;type:varchar(26)"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (UserModel) TableName() string { return "users" }

func (m *UserModel) ToDomain() *domain.User {
	return domain.NuevoUsuarioDesdeBD(
		m.ID,
		m.Email,
		m.FullName,
		m.PasswordHash,
		m.Active,
		m.CreatedBy,
		m.CreatedAt.Format(time.RFC3339),
	)
}
