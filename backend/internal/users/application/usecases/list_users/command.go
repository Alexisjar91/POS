package listusers

import (
	"github.com/Alexisjar91/POS/pkg/especificacion"
	"github.com/Alexisjar91/POS/pkg/paginacion"
)

type ComandoListarUsuarios struct {
	EjecutorID    string
	Especificacion especificacion.Especificacion
	Paginacion    paginacion.Paginacion
}

func (cmd *ComandoListarUsuarios) Validar() error {
	if cmd.EjecutorID == "" {
		return ErrEjecutorRequerido
	}
	return nil
}
