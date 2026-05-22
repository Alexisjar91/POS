package manageroles

type ComandoActualizarRol struct {
	RoleID      string
	Nombre      string
	Descripcion string
	EjecutorID  string
}

func (cmd *ComandoActualizarRol) Validar() error {
	if cmd.RoleID == "" {
		return ErrRolRequerido
	}
	if cmd.Nombre == "" {
		return ErrNombreRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
