package postgres

import (
	"github.com/Alexisjar91/POS/internal/users/domain"
	"github.com/Alexisjar91/POS/pkg/permissions"
)

func init() {
	permissions.Register(
		permissions.Permission{Code: domain.CreateUser, Description: "Crear cuentas de usuario", Module: "users"},
		permissions.Permission{Code: domain.DisableUser, Description: "Desactivar usuarios", Module: "users"},
		permissions.Permission{Code: domain.EnableUser, Description: "Reactivar usuarios", Module: "users"},
		permissions.Permission{Code: domain.AssignRole, Description: "Asignar y remover roles de usuarios", Module: "users"},
		permissions.Permission{Code: domain.ManageRoles, Description: "Crear, editar y eliminar roles del sistema", Module: "users"},
		permissions.Permission{Code: domain.ViewUsers, Description: "Listar usuarios y ver detalles de perfil", Module: "users"},
		permissions.Permission{Code: domain.ResetUserPassword, Description: "Forzar el cambio de contraseña a otro usuario", Module: "users"},
	)
}
