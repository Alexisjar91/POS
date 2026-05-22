package domain

import "errors"

// Errores de entidad User
var (
	ErrEmailRequerido        = errors.New("el email es requerido")
	ErrNombreRequerido       = errors.New("el nombre completo es requerido")
	ErrPasswordHashRequerido = errors.New("el hash de contraseña es requerido")
	ErrCreatedByRequerido    = errors.New("el usuario creador es requerido")
	ErrUsuarioYaActivo       = errors.New("el usuario ya está activo")
	ErrUsuarioYaInactivo     = errors.New("el usuario ya está inactivo")
)

// Errores de entidad Role
var (
	ErrNombreRolRequerido = errors.New("el nombre del rol es requerido")
	ErrPermisoYaAsignado   = errors.New("el permiso ya está asignado al rol")
	ErrPermisoNoAsignado   = errors.New("el permiso no está asignado al rol")
)

// Errores de entidad Permission
var (
	ErrCodeRequerido   = errors.New("el código del permiso es requerido")
	ErrModuleRequerido = errors.New("el módulo del permiso es requerido")
)

// Errores de repositorio
var (
	ErrEmailDuplicado      = errors.New("el email ya está registrado")
	ErrUsuarioNoEncontrado = errors.New("usuario no encontrado")
	ErrRolNoEncontrado     = errors.New("rol no encontrado")
	ErrPermisoNoEncontrado = errors.New("permiso no encontrado")
	ErrRolDuplicado        = errors.New("el nombre del rol ya existe")
	ErrPermisoDuplicado    = errors.New("el código del permiso ya existe")
	ErrRolSistemaInmutable = errors.New("el rol de sistema no puede ser modificado")
	ErrRolConUsuarios      = errors.New("el rol tiene usuarios asignados")
	ErrPermisoEnUso        = errors.New("el permiso está asociado a uno o más roles")
	ErrRepositorio         = errors.New("error interno del repositorio")
)

// Errores de autorización
var ErrAccesoDenegado = errors.New("acceso denegado")
