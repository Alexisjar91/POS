# Plan de Implementación — Capa de Presentación (Módulo Usuario)

> Pendiente de implementar.

## Deuda técnica previa

Antes de comenzar la capa de presentación, hay que saldar una deuda de la capa de infraestructura:

### Especificación y Búsqueda (Specification + Pagination)

La interfaz `Listar` del `UserRepository` debe aceptar filtros dinámicos y paginación, según lo especificado en `spec-user-domain.md`:

```
| Listar | Context, filtros (opcional: activo, rol, paginación) | []User | ErrRepositorio |
```

Se requiere:

1. **Pagination**: struct compartido en `pkg/pagination/` con Page, PageSize, Total, Offset
2. **UserSpecification**: struct en `internal/users/domain/` con campos opcionales (Active *bool, RoleID *string, Query *string para búsqueda textual)
3. **Actualizar UserRepository.Listar**: aceptar spec + pagination, construir query con filtros dinámicos y LIMIT/OFFSET
4. **Actualizar list_users**: aceptar parámetros de filtro y paginación, pasarlos al repositorio

Esta deuda se resuelve como **Paso 0** antes de la presentación.

## Estado actual del proyecto (pre-presentación)

```
backend/
├── pkg/
│   └── permissions/          ← registry global (nuevo)
├── internal/
│   ├── users/
│   │   ├── domain/           ← entidades, interfaces, errores, constantes
│   │   ├── application/      ← 10 casos de uso
│   │   └── infrastructure/   ← persistencia Postgres completa
│   └── (presentación)        ← vacío
└── docs/
    ├── spec-presentation-layer.md  ← especificación Gin + Huma v2
    └── adr/                        ← ADRs
```

## Estructura esperada

```
internal/
└── users/
    ├── domain/
    ├── application/
    ├── infrastructure/
    └── presentation/                     ← nuevo
        ├── dto/                          ← DTOs de entrada/salida
        │   ├── user_dto.go
        │   ├── role_dto.go
        │   └── permission_dto.go
        ├── handler/                      ← Handlers HTTP
        │   ├── user_handler.go
        │   ├── role_handler.go
        │   ├── auth_handler.go
        │   └── routes.go
        ├── middleware/                   ← Middleware (auth, logging)
        │   └── auth_middleware.go
        ├── facade/                       ← Fachada que orquesta casos de uso
        │   ├── user_facade.go
        │   └── role_facade.go
        └── mapper/                       ← Mappers DTO ↔ dominio
            ├── user_mapper.go
            └── role_mapper.go
```

## Tecnologías

- **Router HTTP**: Gin (`github.com/gin-gonic/gin`)
- **Documentación API**: Huma v2 (OpenAPI 3.1 automático)
- **Adaptador**: humagin (Huma + Gin)
- **Swagger UI**: servido en GET /docs

## Principios

- **CON-PRES-001**: Handler → Facade → Mapper → Domain (nunca al revés)
- **CON-PRES-002**: Handler no importa domain ni mapper
- **CON-PRES-003**: Facade no importa Gin ni HTTP
- **CON-PRES-004**: DTOs son structs planos sin comportamiento
- **CON-PRES-005**: Errores HTTP siguen RFC 9457 (Problem Details)

## Pasos

| # | Componente | Estado |
|---|-----------|--------|
| 0 | Specification + Pagination en repositorios | ⬜ Pendiente |
| 1 | DTOs (user_dto, role_dto, permission_dto) | ⬜ Pendiente |
| 2 | Mappers (dominio ↔ DTO) | ⬜ Pendiente |
| 3 | Facades (orquestan casos de uso) | ⬜ Pendiente |
| 4 | Handlers usuarios (CRUD + listado + detalle) | ⬜ Pendiente |
| 5 | Handlers roles (CRUD) | ⬜ Pendiente |
| 6 | Handlers auth (login, refresh — si aplica) | ⬜ Pendiente |
| 7 | Middleware autenticación JWT | ⬜ Pendiente |
| 8 | Rutas + bootstrap (registrar handlers en Gin) | ⬜ Pendiente |
| 9 | Documentación OpenAPI (Huma) | ⬜ Pendiente |

### Paso 0: Specification + Pagination

Crear en orden:

1. **`pkg/pagination/pagination.go`**:
```go
type Pagination struct {
    Page     int   `json:"page"`
    PageSize int   `json:"page_size"`
}

func (p Pagination) Offset() int {
    if p.Page < 1 { p.Page = 1 }
    if p.PageSize < 1 { p.PageSize = 10 }
    return (p.Page - 1) * p.PageSize
}

type PaginatedResult[T any] struct {
    Data       []T        `json:"data"`
    Page       int        `json:"page"`
    PageSize   int        `json:"page_size"`
    Total      int64      `json:"total"`
    TotalPages int        `json:"total_pages"`
}
```

2. **Agregar `UserSpecification` en `internal/users/domain/user_repository.go`** (mismo archivo que la interfaz):

```go
type UserSpecification struct {
    Active *bool
    RoleID *string
    Query  *string // búsqueda por email o nombre
}
```

3. **Modificar `domain.UserRepository`**: cambiar `Listar(ctx)` por `Listar(ctx, spec UserSpecification, page pagination.Pagination) (*pagination.PaginatedResult[*User], error)`

4. **Actualizar implementación en postgres/user_repository.go**: construir query dinámica con SQL raw para filtros + paginación

5. **Actualizar caso de uso list_users** para aceptar y pasar los filtros

### Notas sobre implementación

- Los handlers se implementan con Gin puro (rutas estándar). Huma se agrega al final para la documentación OpenAPI.
- El facade recibe casos de uso por constructor y expone métodos que los handlers llaman.
- El mapper convierte entre DTOs y objetos de dominio.
- El middleware JWT extrae el userID del token y lo inyecta en el contexto de Gin.
- Las rutas protegidas usan el middleware de autenticación.
