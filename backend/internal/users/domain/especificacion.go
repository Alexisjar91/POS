package domain

// ColumnasPermitidas define qué campos pueden filtrarse en usuarios.
var ColumnasPermitidas = map[string]bool{
	"email":       true,
	"fullName":    true,
	"active":      true,
	"createdAt":   true,
	"createdBy":   true,
}

// MapeoColumnas mapea nombres de dominio a columnas DB.
var MapeoColumnas = map[string]string{
	"email":       "email",
	"fullName":    "full_name",
	"active":      "active",
	"createdAt":   "created_at",
	"createdBy":   "created_by",
}
