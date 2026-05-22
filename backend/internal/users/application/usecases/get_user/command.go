package getuser

type ComandoObtenerUsuario struct {
	UserID     string
	EjecutorID string
}

func (cmd *ComandoObtenerUsuario) Validar() error {
	if cmd.UserID == "" {
		return ErrUsuarioRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
