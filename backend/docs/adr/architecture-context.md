# Architecture Context — Presentation Layer

## Decision
La capa de presentación sigue un flujo estricto de responsabilidades.
Ninguna capa puede saltarse a la siguiente. El Handler no conoce el dominio.

## Flujo de capas

```
Domain → Mapper → Facade → Handler → ApiResponse
```

## Responsabilidades por actor

| Actor | Consume | Produce | Regla |
|---|---|---|---|
| `Mapper` | `Domain struct` | `DTO struct` | No tiene lógica de negocio |
| `Facade` | `Domain (via use case)` + `Mapper` | `DTO` | Agrupa casos de uso relacionados |
| `Handler` | `Facade` | `ApiResponse[T]` | No toca dominio ni mapper directamente |
| `ApiResponse` | `DTO cualquiera` | `JSON HTTP response` | Genérico, sirve cualquier recurso |

---

## Estructura de carpetas

```
internal/
└── user/
    ├── domain/
    │   └── user.go               # User (domain struct)
    ├── application/
    │   └── user_use_cases.go     # casos de uso
    └── presentation/
        ├── dto/
        │   └── user_dto.go       # UserDTO
        ├── mapper/
        │   └── user_mapper.go    # UserMapper struct
        ├── facade/
        │   └── user_facade.go    # UserFacade struct
        └── handler/
            └── user_handler.go   # UserHandler struct

shared/
└── presentation/
    └── api_response.go           # ApiResponse[T] genérico + HATEOAS
```

---

## Ejemplos mínimos

### Domain — `internal/user/domain/user.go`
```go
package domain

type User struct {
    ID       string
    Name     string
    Email    string
    Active   bool
}
```

---

### DTO — `internal/user/presentation/dto/user_dto.go`
```go
package dto

type UserDTO struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

---

### Mapper — `internal/user/presentation/mapper/user_mapper.go`
```go
package mapper

import (
    "myapp/internal/user/domain"
    "myapp/internal/user/presentation/dto"
)

type UserMapper struct{}

func (m UserMapper) ToDTO(u domain.User) dto.UserDTO {
    return dto.UserDTO{
        ID:    u.ID,
        Name:  u.Name,
        Email: u.Email,
    }
}
```

---

### Facade — `internal/user/presentation/facade/user_facade.go`
```go
package facade

import (
    "myapp/internal/user/application"
    "myapp/internal/user/presentation/dto"
    "myapp/internal/user/presentation/mapper"
)

type UserFacade struct {
    useCases application.UserUseCases
    mapper   mapper.UserMapper
}

func NewUserFacade(uc application.UserUseCases) UserFacade {
    return UserFacade{useCases: uc, mapper: mapper.UserMapper{}}
}

func (f UserFacade) GetByID(id string) (dto.UserDTO, error) {
    user, err := f.useCases.FindByID(id)
    if err != nil {
        return dto.UserDTO{}, err
    }
    return f.mapper.ToDTO(user), nil
}
```

---

### ApiResponse — `shared/presentation/api_response.go`
```go
package presentation

type Link struct {
    Href   string `json:"href"`
    Method string `json:"method"`
}

type ApiResponse[T any] struct {
    Data  T                `json:"data"`
    Links map[string]Link  `json:"_links,omitempty"`
}

func NewResponse[T any](data T, links map[string]Link) ApiResponse[T] {
    return ApiResponse[T]{Data: data, Links: links}
}
```

---

### Handler — `internal/user/presentation/handler/user_handler.go`
```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "myapp/internal/user/presentation/facade"
    shared "myapp/shared/presentation"
)

type UserHandler struct {
    facade facade.UserFacade
}

func NewUserHandler(f facade.UserFacade) UserHandler {
    return UserHandler{facade: f}
}

func (h UserHandler) GetByID(c *gin.Context) {
    id := c.Param("id")

    userDTO, err := h.facade.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    response := shared.NewResponse(userDTO, map[string]shared.Link{
        "self":   {Href: "/users/" + id, Method: "GET"},
        "update": {Href: "/users/" + id, Method: "PUT"},
        "delete": {Href: "/users/" + id, Method: "DELETE"},
    })

    c.JSON(http.StatusOK, response)
}
```

---

## Reglas que la IA debe respetar siempre

- El `Handler` **nunca** importa `domain` ni `mapper`.
- El `Facade` **nunca** importa `gin` ni nada HTTP.
- El `Mapper` **nunca** tiene lógica de negocio, solo conversión de structs.
- El `ApiResponse` es genérico con `[T any]`, nunca tipado a un recurso específico.
- Los links HATEOAS los construye el `Handler`, no el `Facade`.
- Un nuevo recurso (Casa, Carro, etc.) sigue exactamente esta misma estructura de carpetas.
