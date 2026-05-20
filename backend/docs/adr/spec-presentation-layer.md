---
title: Especificación de la Capa de Presentación — Gin + Huma v2 + OpenAPI/Swagger
version: 1.0
date_created: 2026-05-10
owner: Equipo Identidad
tags: presentation, api, huma, gin, swagger, openapi
---

# Especificación de la Capa de Presentación

## 1. Propósito y Alcance

Definir la arquitectura, frameworks, restricciones y contratos de la capa de presentación del servicio `identidad`. Esta especificación cubre desde el enrutamiento HTTP hasta la respuesta JSON, incluyendo la generación automática de documentación OpenAPI/Swagger.

**Alcance:**
- Definición de frameworks: Gin como router HTTP, Huma v2 (vía `humagin`) como generador de OpenAPI.
- Auto-documentación de la API mediante Swagger UI servida en `/docs`.
- Endpoints iniciales: `POST /api/v1/auth/register` y `POST /api/v1/auth/login`.
- Flujo de capas: Handler → Facade → Mapper → Domain (según `architecture-context.md`).
- Manejo de errores HTTP estandarizado.
- Middleware de autenticación JWT (extraer `usuarioID` y `sesionID` del token).

**Fuera de alcance:**
- Lógica de negocio de login, registro, refresh, logout (definida en `login_spec.md`).
- Implementación concreta de servicios de aplicación, repositorios, o infraestructura.
- Seguridad perimetral (rate limiting, IP blocking) — se integra como middleware separado.
- Despliegue, CI/CD, o configuración de contenedores.

## 2. Definiciones

| Término | Definición |
|---------|-----------|
| **Gin** | Framework HTTP para Go (github.com/gin-gonic/gin). Maneja el enrutamiento, parsing de request, y escritura de response. |
| **Huma v2** | Framework de API REST para Go que genera OpenAPI 3.1 a partir de tipos Go (huma.rocks). |
| **humagin** | Adaptador de Huma v2 para Gin. Permite definir recursos Huma sobre un router Gin. |
| **OpenAPI 3.1** | Especificación estándar para describir APIs REST. Huma la genera automáticamente. |
| **Swagger UI** | Interfaz web interactiva para explorar y probar endpoints de la API. Servida en `/docs`. |
| **Handler** | Función que recibe un request HTTP y orquesta la respuesta. No contiene lógica de negocio. |
| **Facade** | Capa que agrupa casos de uso relacionados y aplica mappers. Ver `architecture-context.md`. |
| **ApiResponse[T]** | Estructura genérica para todas las respuestas JSON. Ver `architecture-context.md`. |
| **HATEOAS** | Hipermedia como motor del estado de la aplicación. Links en las respuestas construidos por el Handler. |
| **DTO** | Data Transfer Object. Estructura plana que viaja entre capas. No tiene comportamiento. |

## 3. Requisitos, Restricciones y Guías

### Frameworks (obligatorios)

- **REQ-PRES-001**: El router HTTP debe ser Gin (github.com/gin-gonic/gin).
- **REQ-PRES-002**: La documentación OpenAPI debe generarse automáticamente vía Huma v2 (github.com/danielgtaylor/huma/v2).
- **REQ-PRES-003**: El adaptador Huma-Gin debe ser `humagin` (github.com/danielgtaylor/huma/v2/adapters/humagin).
- **REQ-PRES-004**: Swagger UI debe servirse en la ruta `GET /docs`.
- **REQ-PRES-005**: La especificación OpenAPI debe servirse en `GET /openapi.json` o `GET /openapi.yaml`.

### Arquitectura de capas

- **CON-PRES-001**: Se sigue estrictamente el flujo `Handler → Facade → Mapper → Domain` definido en `architecture-context.md`.
- **CON-PRES-002**: El Handler **nunca** importa paquetes de `domain` ni `mapper`.
- **CON-PRES-003**: El Facade **nunca** importa Gin ni nada HTTP.
- **CON-PRES-004**: El `ApiResponse[T]` es genérico. Nunca se crea una versión tipada por recurso.
- **CON-PRES-005**: Los links HATEOAS los construye el Handler, no el Facade.

### Endpoints

- **REQ-PRES-006**: `POST /api/v1/auth/register` — Registro de nuevo usuario.
- **REQ-PRES-007**: `POST /api/v1/auth/login` — Inicio de sesión, devuelve tokens JWT.
- **REQ-PRES-008**: `GET /health` — Health check (existente, se mantiene).
- **REQ-PRES-009**: `GET /docs` — Swagger UI.
- **REQ-PRES-010**: `GET /openapi.json` — Especificación OpenAPI.

### Formato de respuestas

- **REQ-PRES-011**: Todas las respuestas exitosas usan `ApiResponse[T]` del paquete `shared/presentation/api_response.go`.
- **REQ-PRES-012**: Las respuestas de error usan `Problem Details` (RFC 9457) via Huma (`huma.ErrorModel` o `huma.StatusError`).
- **REQ-PRES-013**: Los errores HTTP no revelan detalles internos de implementación.

### Documentación

- **REQ-PRES-014**: Todos los endpoints deben tener `op` (operation) definido con `OperationID`, `Summary`, y `Description` via Huma.
- **REQ-PRES-015**: Los modelos de request/response se documentan automáticamente via struct tags `json`, `doc`, `example`.
- **REQ-PRES-016**: Los errores posibles deben documentarse con `errors.Code` en Huma o responses documentados.

### Versiones

- **REQ-PRES-017**: La API se versiona via prefijo de ruta `/api/v1/`.
- **REQ-PRES-018**: La versión se refleja en la especificación OpenAPI (`info.version`).

### Middleware

- **REQ-PRES-019**: Middleware de autenticación JWT para rutas protegidas (extrae `usuarioID`, `sesionID` de claims).
- **REQ-PRES-020**: Middleware de CORS configurable via entorno.
- **REQ-PRES-021**: Middleware de logging de requests (Gin default o personalizado).

## 4. Interfaces y Contratos de Datos

### Estructura de endpoints

```
Método  Ruta                      Auth     Handler              Descripción
------  ------------------------  -------  -------------------  ---------------------------
GET     /health                   No       HealthHandler        Health check
POST    /api/v1/auth/register     No       RegisterHandler      Registro de usuario
POST    /api/v1/auth/login        No       LoginHandler         Inicio de sesión
GET     /docs                     No       (Swagger UI)         Documentación interactiva
GET     /openapi.json             No       (Huma)               Especificación OpenAPI
```

### Contratos de los endpoints

#### POST /api/v1/auth/register

**Request body:**
| Campo    | Tipo   | Requerido | Descripción |
|----------|--------|-----------|-------------|
| name     | string | sí        | Nombre completo del usuario |
| email    | string | sí        | Correo electrónico (válido) |
| password | string | sí        | Contraseña (mínimo 8 caracteres) |

**Response 201 (success):**
| Campo            | Tipo   | Descripción |
|------------------|--------|-------------|
| data.id          | string | ID del usuario creado |
| data.name        | string | Nombre del usuario |
| data.email       | string | Email del usuario |
| _links.self      | object | Link al perfil del usuario |

**Response 400 (validation error):**
RFC 9457 Problem Details.

#### POST /api/v1/auth/login

**Request body:**
| Campo    | Tipo   | Requerido | Descripción |
|----------|--------|-----------|-------------|
| email    | string | sí        | Correo electrónico |
| password | string | sí        | Contraseña |

**Response 200 (success):**
| Campo                   | Tipo   | Descripción |
|-------------------------|--------|-------------|
| data.access_token       | string | JWT access token |
| data.refresh_token      | string | JWT refresh token |
| data.expires_in         | int    | Segundos hasta expiración del access token |
| data.token_type         | string | Siempre "Bearer" |
| data.usuario_id         | string | ID del usuario autenticado |
| _links.self             | object | Link al perfil |
| _links.refresh          | object | Link para renovar token |

**Response 401 (unauthorized):**
RFC 9457 Problem Details.

### ApiResponse genérico

Definido en `shared/presentation/api_response.go`. Ver `architecture-context.md` para la estructura exacta.

```go
// Estructura conceptual:
type ApiResponse[T any] struct {
    Data  T                `json:"data"`
    Links map[string]Link  `json:"_links,omitempty"`
}

type Link struct {
    Href   string `json:"href"`
    Method string `json:"method"`
}
```

## 5. Criterios de Aceptación

- **AC-PRES-001**: Dado que la aplicación inicia, Cuando se accede a `GET /docs`, Entonces se renderiza Swagger UI con los endpoints documentados.
- **AC-PRES-002**: Dado que la aplicación inicia, Cuando se accede a `GET /openapi.json`, Entonces se retorna un JSON válido conforme a OpenAPI 3.1.
- **AC-PRES-003**: Dado un request válido a `POST /api/v1/auth/register`, Cuando se procesa, Entonces la respuesta tiene código 201 y body con `ApiResponse`.
- **AC-PRES-004**: Dado un request inválido a `POST /api/v1/auth/register` (email mal formado), Cuando se procesa, Entonces la respuesta tiene código 400 con `Problem Details`.
- **AC-PRES-005**: Dado un request válido a `POST /api/v1/auth/login`, Cuando las credenciales son correctas, Entonces la respuesta tiene código 200 con tokens JWT.
- **AC-PRES-006**: Dado un request válido a `POST /api/v1/auth/login`, Cuando las credenciales son incorrectas, Entonces la respuesta tiene código 401.
- **AC-PRES-007**: Dado que el Handler recibe un request, Cuando se ejecuta, Entonces no importa paquetes de `domain` ni `mapper` (validación de imports).
- **AC-PRES-008**: Dado un request sin token a una ruta protegida, Cuando se evalúa el middleware, Entonces la respuesta tiene código 401.

## 6. Estrategia de Automatización de Pruebas

- **Niveles de prueba**: Unitarias (Handler con Facade mockeado), Integración (end-to-end con servidor de prueba).
- **Framework de testing**: `testing` estándar de Go + `testify` (ya existente en el proyecto).
- **Mocks**: Los Handlers reciben interfaces de Facade, se mockean con `testify/mock`.
- **Pruebas de integración**: Usar `httptest.NewServer` con Gin router real y Huma configurado.
- **Validación OpenAPI**: Generar el `openapi.json` y validar contra schema de OpenAPI 3.1 en CI.
- **Cobertura mínima**: 80% en capa de presentación (handlers, middleware).

## 7. Justificación y Contexto

### ¿Por qué Gin + Huma v2?

- Gin es el framework HTTP más adoptado en Go, con alto rendimiento y ecosistema maduro.
- Huma v2 genera OpenAPI 3.1 automáticamente a partir de los tipos Go, eliminando la necesidad de anotaciones manuales de Swagger.
- La combinación `humagin` permite usar Gin como router real mientras Huma maneja la documentación, union de lo mejor de ambos mundos.
- Huma valida automáticamente los request bodies contra los struct tags, reduciendo código boilerplate en los handlers.

### ¿Por qué OpenAPI/Swagger?

- Estandariza la documentación de la API para consumo por otros equipos y generación de clientes.
- Swagger UI permite a desarrolladores frontend y otros servicios explorar la API sin necesidad de leer código.
- La generación automática via Huma garantiza que la documentación siempre esté sincronizada con la implementación.

### Relación con la arquitectura existente

- El handler existente (`internal/handler/handler.go`) usa `net/http` directo. Esta especificación migra a Gin + Huma manteniendo el health check existente.
- La estructura `ApiResponse` ya definida en `architecture-context.md` se mantiene y se integra con Huma.
- Los servicios de aplicación existentes (login, refresh, logout) se exponen via Facades, no se tocan directamente desde handlers.

## 8. Dependencias e Integraciones Externas

### Dependencias de Framework

| Paquete | Propósito | ¿Ya en go.mod? |
|---------|-----------|----------------|
| github.com/gin-gonic/gin | Router HTTP | No |
| github.com/danielgtaylor/huma/v2 | Framework OpenAPI | No |
| github.com/danielgtaylor/huma/v2/adapters/humagin | Adaptador Huma-Gin | No |

### Dependencias del Proyecto

- `internal/handler/handler.go` — Handler existente con health check. Migrar a Gin manteniendo endpoint `/health`.
- `internal/registry/registry.go` — Registry central. Los handlers obtendrán facades desde aquí.
- `internal/config/env.go` — Se agregarán variables de entorno para CORS, puerto, etc.
- `shared/presentation/api_response.go` — Struct genérico ApiResponse. Ver `architecture-context.md`.
- `../../adr/architecture-context.md` — Define el flujo de capas y restricciones.
- `../sesiones/login_spec.md` — Define los casos de uso de login y registro que los handlers expondrán.

### Dependencias de Infraestructura

- Puerto configurable via `PORT` (default `8080`).
- Variable `CORS_ORIGINS` para configuración de CORS.
- Variable `JWT_SECRET` para validación de tokens en middleware.

## 10. Criterios de Validación

- [ ] El servidor inicia con Gin como router y Huma configurado.
- [ ] `GET /docs` muestra Swagger UI con todos los endpoints documentados.
- [ ] `GET /openapi.json` retorna un JSON OpenAPI 3.1 válido.
- [ ] Los endpoints `POST /api/v1/auth/register` y `POST /api/v1/auth/login` existen y están documentados.
- [ ] Las respuestas exitosas usan `ApiResponse[T]` con estructura consistente.
- [ ] Las respuestas de error siguen RFC 9457 (Problem Details).
- [ ] Los handlers no importan paquetes de dominio.
- [ ] Los facades no importan Gin ni HTTP.
- [ ] El middleware JWT rechaza requests sin token en rutas protegidas.
- [ ] La documentación OpenAPI se genera sin anotaciones manuales de swagger (solo struct tags Go).

## 11. Especificaciones Relacionadas / Lectura Adicional

- `../../adr/architecture-context.md` — Arquitectura de capas de presentación.
- `../../adr/feature-template.md` — Template de features del proyecto.
- `../sesiones/login_spec.md` — Especificación del caso de uso login.
- [Documentación de Huma v2](https://huma.rocks/) — Framework de API REST.
- [Documentación de Gin](https://gin-gonic.com/) — Framework HTTP.
