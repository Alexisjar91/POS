# Feature Template — Identidad

> Basado en el dominio `usuarios` como ejemplo canónico.
> Define lo obligatorio vs contextual. No es receta, es plantilla.

---

## Estructura

```
internal/{modulo}/
├── domain/          ← Obligatorio
│   ├── {entidad}.go           # Entidad + constructor(es) + comportamiento
│   ├── {entidad}_repositorio.go  # Interfaz del repositorio
│   └── errores.go              # Errores de dominio
├── application/
│   └── services/
│       └── {caso-uso}/
│           ├── comando.go      # Input DTO (contextual)
│           ├── respuesta.go    # Output DTO (contextual)
│           └── servicio.go     # Caso de uso (contextual)
└── infrastructure/
    └── persistence/
        └── postgres/
            ├── {entidad}_model.go     # Modelo GORM (obligatorio si hay DB)
            └── {entidad}_repositorio.go  # Impl del repositorio (obligatorio si hay DB)
```

---

## Domain — Lo obligatorio

### Entidad

```go
// {entidad}.go
package domain

type Entidad struct {
    // campos privados — solo getters públicos
}

// NuevoEntidad — constructor para crear nuevas instancias.
// NO valida ID (lo asigna app/infraestructura).
// Valida solo reglas de negocio.
func NuevoEntidad(campo1 string, campo2 string) (*Entidad, error) {
    // validar reglas de negocio
    if campo1 == "" {
        return nil, ErrCampoRequerido
    }
    return &Entidad{...}, nil
}

// NuevoEntidadDesdeBD — constructor para reconstruir desde persistencia.
// Sin validaciones, trae todo.
func NuevoEntidadDesdeBD(id string, campo1 string, ...) *Entidad {
    return &Entidad{...}
}
```

**Reglas:**
- `NuevoEntidad` → valida reglas de negocio, NO el ID
- `NuevoEntidadDesdeBD` → NO valida nada, asume datos consistentes
- Getters para todos los campos privados
- Comportamiento como métodos de la entidad (`EstaActiva()`, `CambiarEstado()`, etc.)

### Interfaz del repositorio

```go
// {entidad}_repositorio.go
type EntidadRepositorio interface {
    Crear(ctx context.Context, entidad *Entidad) (*Entidad, error)
    ObtenerPorID(ctx context.Context, id string) (*Entidad, error)
    // contextual:
    // Actualizar, Eliminar, Listar, ObtenerPorX, etc.
}
```

### Errores

```go
// errores.go
var (
    ErrCampoRequerido    = errors.New("...")
    ErrReglaDeNegocio    = errors.New("...")
    ErrTransicionInvalida = errors.New("...")
)
```

---

## Application — Contextual

Solo si hay un caso de uso que lo requiera.

```go
// comando.go — input DTO
type ComandoAlgo struct {
    Campo1 string
}

// respuesta.go — output DTO
type RespuestaAlgo struct {
    ID    string
    Estado string
}

// servicio.go — orquestador
type ServicioAlgo struct {
    repo EntidadRepositorio
}

func (s *ServicioAlgo) Ejecutar(ctx context.Context, cmd *ComandoAlgo) (*RespuestaAlgo, error) {
    // 1. validar comando
    // 2. entidad := NuevoEntidad(...)
    // 3. repo.Crear(ctx, entidad)
    // 4. armar respuesta
}
```

**Si toca múltiples tablas** (ej: Registro → `usuarios` + `credenciales`):
- Usar `UnitOfWork` en lugar de repo directo
- La interfaz `UnitOfWork` vive en dominio
- La implementación en infraestructura

---

## Infrastructure — Obligatorio si hay DB

### Modelo GORM

```go
type EntidadModel struct {
    ID    string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Campo string `gorm:"column:campo"`
}

func (EntidadModel) TableName() string { return "entidades" }

func (m *EntidadModel) ToDomain() *domain.Entidad {
    return domain.NuevoEntidadDesdeBD(m.ID, m.Campo, ...)
}

func FromDomain(e *domain.Entidad) (*EntidadModel, error) {
    return &EntidadModel{ID: e.ID(), Campo: e.Campo()}, nil
}
```

### Repositorio impl

```go
type entidadRepositorio struct {
    db *gorm.DB
}

func NewEntidadRepositorio(db *gorm.DB) domain.EntidadRepositorio {
    return &entidadRepositorio{db: db}
}

func (r *entidadRepositorio) Crear(ctx context.Context, e *domain.Entidad) (*domain.Entidad, error) {
    model, _ := FromDomain(e)
    r.db.WithContext(ctx).Create(model)
    return model.ToDomain(), nil
}
```

---

## Migración — Agregar a `config/database.go`

```go
func RunMigrations(db *gorm.DB) error {
    // ...
    db.AutoMigrate(&{modulo}_postgres.EntidadModel{})
    return nil
}
```

---

## Check-list para una feature nueva

| Capa | Archivo | ¿Obligatorio? |
|------|---------|---------------|
| Domain | `{entidad}.go` | ✅ Siempre |
| Domain | `{entidad}_repositorio.go` | ✅ Siempre |
| Domain | `errores.go` | ✅ Siempre |
| App | `comando.go` | Solo si hay caso de uso |
| App | `respuesta.go` | Solo si hay caso de uso |
| App | `servicio.go` | Solo si hay caso de uso |
| App | `unit_of_work.go` | Solo si toca múltiples tablas |
| Infra | `{entidad}_model.go` | ✅ Si hay DB |
| Infra | `{entidad}_repositorio.go` | ✅ Si hay DB |
| Config | `database.go` + `AutoMigrate` | ✅ Si hay tabla nueva |
| Registry | `registry.go` | ✅ Si hay nuevo repo/UoW |
