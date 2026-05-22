package domain

// — Users module permissions —

// CreateUser permite crear nuevas cuentas de usuario.
const CreateUser = "create_user"

// DisableUser permite desactivar usuarios. No aplica al usuario OWNER.
const DisableUser = "disable_user"

// EnableUser permite reactivar usuarios desactivados.
const EnableUser = "enable_user"

// AssignRole permite asignar o remover roles de un usuario. No aplica al rol OWNER.
const AssignRole = "assign_role"

// ManageRoles permite crear, editar y eliminar roles, así como gestionar su asignación de permisos.
const ManageRoles = "manage_roles"

// ViewUsers permite listar usuarios y ver detalles de perfil.
const ViewUsers = "view_users"

// ResetUserPassword permite forzar el cambio de contraseña a otro usuario.
const ResetUserPassword = "reset_user_password"
