package resetpassword

type ComandoResetearContrasena struct {
	UserID      string // usuario objetivo
	NewPassword string // nueva contraseña en texto plano
	EjecutorID  string // quien ejecuta (debe tener ResetUserPassword)
}

func (cmd *ComandoResetearContrasena) Validar() error {
	if cmd.UserID == "" {
		return ErrUsuarioRequerido
	}
	if cmd.NewPassword == "" {
		return ErrNuevaContrasenaRequerida
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
