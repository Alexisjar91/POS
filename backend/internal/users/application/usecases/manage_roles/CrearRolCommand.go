package manageroles

type ComandoCrearRol struct {
	Nombre      string
	Descripcion string
	EjecutorID  string
}

func (cmd *ComandoCrearRol) Validar() error {
	if cmd.Nombre == "" {
		return ErrNombreRequerido
	}
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
