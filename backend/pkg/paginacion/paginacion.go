// Package paginacion provee tipos genéricos para paginación y ordenación de resultados.
package paginacion

// Paginacion define los parámetros de paginación de una consulta.
type Paginacion struct {
	Pagina       int
	TamanoPagina int
	Ordenaciones []Ordenacion
}

// Ordenacion define un criterio de ordenación.
type Ordenacion struct {
	Campo string
	Tipo  TipoOrdenacion
}

// TipoOrdenacion define el tipo de ordenación.
type TipoOrdenacion string

const (
	ASC  TipoOrdenacion = "ASC"
	DESC TipoOrdenacion = "DESC"
)

// Offset calcula el desplazamiento para la consulta SQL.
func (p Paginacion) Offset() int {
	if p.Pagina < 1 {
		p.Pagina = 1
	}
	return (p.Pagina - 1) * p.TamanoPagina
}

// Limit retorna el tamaño de página normalizado.
func (p Paginacion) Limit() int {
	if p.TamanoPagina < 1 {
		return 10
	}
	return p.TamanoPagina
}
