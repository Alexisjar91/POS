package manageroles

type ComandoEliminarRol struct {
	RoleID     string
	EjecutorID string
}

func (cmd *ComandoEliminarRol) Validar() error {
	if cmd.RoleID == "" {
		return ErrRolRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
