// Package especificacion provee tipos genéricos para filtrado y búsqueda de resultados.
package especificacion

// Especificacion define los filtros disponibles para búsquedas.
type Especificacion struct {
	Filtros []CriterioFiltro
}

// CriterioFiltro define un criterio individual de filtrado.
type CriterioFiltro struct {
	Campo    string
	Operador string
	Valor    any
}
