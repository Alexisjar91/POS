package createuser

import (
	"net/mail"
	"strings"
)

// ComandoCrearUsuario contiene los datos necesarios para crear un usuario.
type ComandoCrearUsuario struct {
	Email      string
	FullName   string
	Password   string // contraseña en texto plano (se hashea en infraestructura)
	CreatedBy  string // ID del usuario que crea
	EjecutorID string // ID del usuario autenticado (para autorización)
}

// Validar ejecuta las validaciones de aplicación sobre el comando.
// Retorna el primer error encontrado.
func (cmd *ComandoCrearUsuario) Validar() error {
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.FullName = strings.TrimSpace(cmd.FullName)

	if cmd.Email == "" {
		return ErrEmailRequerido
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return ErrEmailInvalido
	}
	if cmd.FullName == "" {
		return ErrNombreRequerido
	}
	if cmd.Password == "" {
		return ErrPasswordRequerido
	}
	if cmd.CreatedBy == "" {
		return ErrCreadorRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
