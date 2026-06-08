# Delta Spec - Extender InfoPage con Nuevas Entidades

**Status**: `approved`
**Fecha**: 2026-06-04
**Owner**: cristiansrc
**Proyecto**: `hv-go-ms-resume`
**Incremento**: `extend-info-page-new-entities`
**Tipo**: Mejora funcional
**Base**: `docs/specs/increments/migration-to-go.md` (implemented)

---

## 1. Objetivo

Extender el endpoint `GET /v1/ms-resume/public/info-page` para que incluya las 5 nuevas entidades (Course, Certification, Language, Reference, CustomSection) en la respuesta agregada. Los repositorios, DTOs, handlers y rutas ya existen del incremento `migration-to-go`; solo se requiere modificar `InfoPageResponse` y `PublicService.GetInfoPage()`.

---

## 2. Contexto y Justificación

El incremento `migration-to-go` (ya `implemented`) incorporó 5 nuevas entidades (Course, Certification, Language, Reference, CustomSection) con sus CRUD completos y endpoints públicos individuales (`/public/courses`, `/public/certifications`, `/public/languages-data`, `/public/references`, `/public/custom-sections`). Sin embargo, el endpoint agregador `/public/info-page` —que el frontend del portal (`hv-rt-fr-portal`) consume para renderizar la página principal— nunca se actualizó para incluirlas.

Esto fuerza al frontend a hacer 5 llamadas adicionales por carga de página, incrementando latencia y complejidad del cliente.

**Comportamiento actual**: `InfoPageResponse` contiene solo 6 secciones: home, basicData, skills, experiences, educations, altchaChallenge.
**Comportamiento deseado**: `InfoPageResponse` contiene 11 secciones (las 6 existentes + courses, certifications, languages, references, customSections).

---

## 3. Arquitectura y Diseño

### 3.1 Cambios en el Modelo de Datos

**No hay cambios en el esquema de base de datos.** Las tablas ya existen y los repositorios (`CourseRepository`, `CertificationRepository`, `LanguageRepository`, `ReferenceRepository`, `CustomSectionRepository`) ya están implementados y funcionando.

### 3.2 Cambios en la Capa de Aplicación

#### 3.2.1 DTO `InfoPageResponse` (único cambio estructural)

**Archivo**: `internal/application/dto/response/info_page.go`

Se agregarán 5 nuevos campos al struct `InfoPageResponse`:

| Campo Go | Tipo | JSON Tag | Descripción |
|---|---|---|---|
| `Courses` | `[]CourseResponse` | `json:"courses"` | Lista de cursos activos ordenados por `order ASC` |
| `Certifications` | `[]CertificationResponse` | `json:"certifications"` | Lista de certificaciones activas ordenadas por `order ASC` |
| `Languages` | `[]LanguageResponse` | `json:"languages"` | Lista de idiomas activos ordenados por `order ASC` |
| `References` | `[]ReferenceResponse` | `json:"references"` | Lista de referencias activas ordenadas por `order ASC` |
| `CustomSections` | `[]CustomSectionResponse` | `json:"customSections"` | Lista de secciones personalizadas activas **y visibles** ordenadas por `order ASC` |

**Orden en el struct Go**: Los 5 nuevos campos se agregarán al final del struct `InfoPageResponse`, después de `AltchaChallenge`, en el siguiente orden semántico: `Courses`, `Certifications`, `Languages`, `References`, `CustomSections`. Los 6 campos existentes se mantienen en su posición actual.

**Regla de serialización**: Los arrays DEBEN incluir `omitempty` en el tag JSON para mantener consistencia con el estilo existente de `InfoPageResponse`. Sin embargo, el servicio siempre los setea como slice no-nil (array vacío si no hay datos), por lo que el campo siempre aparecerá en la respuesta JSON. Esto es compatible con los frontends existentes (estándar JSON: campos desconocidos son ignorados).

#### 3.2.2 Servicio `PublicService.GetInfoPage()` (único cambio de lógica)

**Archivo**: `internal/application/service/public_service.go`

Se agregarán 5 consultas secuenciales dentro de `GetInfoPage()`:

1. `s.courseRepo.List(ctx)` → `[]entity.Course`
2. `s.certificationRepo.List(ctx)` → `[]entity.Certification`
3. `s.languageRepo.List(ctx)` → `[]entity.Language`
4. `s.referenceRepo.List(ctx)` → `[]entity.Reference`
5. `s.customSectionRepo.List(ctx)` → `[]entity.CustomSection`, filtrar `cs.Visible == true`

Cada lista se mapeará a su DTO de respuesta correspondiente usando las funciones de mapeo que **ya existen en el package `service`** (en archivos separados: `course_service.go`, `certification_service.go`, `language_service.go`, `reference_service.go`, `custom_section_service.go`):
- `mapCourseToResponse()` — definida en `course_service.go`
- `mapCertificationToResponse()` — definida en `certification_service.go`
- `mapLanguageToResponse()` — definida en `language_service.go`
- `mapReferenceToResponse()` — definida en `reference_service.go`
- `mapCustomSectionToResponse()` — definida en `custom_section_service.go`

Todas son accesibles desde `public_service.go` sin imports adicionales por estar en el mismo package.

**Estrategia de errores (todo o nada)**: Si cualquiera de las 5 consultas falla, el método retorna error inmediatamente (error 500). No se permiten respuestas parciales. El orden de las consultas no importa, pero se ejecutarán secuencialmente.

**Constructor**: NO requiere cambios. Los 5 repositorios ya están inyectados en `PublicService` (ver sección 4.1).

### 3.3 Capa de Dominio

**Sin cambios.** Las entidades de dominio existen y son correctas.

### 3.4 Capa de Infraestructura

**Sin cambios.** Los handlers, rutas, repositorios y configuración ya están implementados.

---

## 4. Contrato API

### 4.1 Endpoint Modificado

`GET /v1/ms-resume/public/info-page`

**Request**: Sin cambios (GET sin body ni parámetros).

**Response 200**: Se extiende el schema `InfoPageResponse` con 5 nuevos campos. Los campos existentes se mantienen idénticos (nombres, tipos, orden semántico).

### 4.2 Esquema `InfoPageResponse` (OpenAPI)

Ver cambios en `api/openapi.yaml`. Los nuevos campos son:

```yaml
courses:
  type: array
  items:
    $ref: "#/components/schemas/CourseResponse"
certifications:
  type: array
  items:
    $ref: "#/components/schemas/CertificationResponse"
languages:
  type: array
  items:
    $ref: "#/components/schemas/LanguageResponse"
references:
  type: array
  items:
    $ref: "#/components/schemas/ReferenceResponse"
customSections:
  type: array
  items:
    $ref: "#/components/schemas/CustomSectionResponse"
```

### 4.3 Códigos de Estado

| Código | Descripción | Cambio |
|---|---|---|
| 200 | Respuesta exitosa con datos completos | Se extiende el schema |
| 500 | Error interno (incluye fallo en consulta de nuevas entidades) | Sin cambios |

### 4.4 Rate Limiting

Sin cambios. Aplica el rate limiting público existente (100 req/min).

### 4.5 Seguridad

Sin cambios. El endpoint permanece público (sin JWT).

### 4.6 Idempotencia

El endpoint es GET, inherentemente idempotente. Sin cambios.

---

## 5. Integraciones

No se requieren nuevas integraciones externas. Todas las entidades se obtienen desde SQLite local.

---

## 6. Operación y Observabilidad

### 6.1 Logs

El endpoint ya genera logs vía el middleware de logging existente. No se requieren logs adicionales.

### 6.2 Métricas

Sin nuevas métricas. Las métricas de latencia y throughput del endpoint existente capturan automáticamente los cambios.

---

## 7. Estrategia de Testing

### 7.1 Tests Unitarios

- **`PublicService.GetInfoPage`**: Agregar test case que verifique que la respuesta incluye las 5 nuevas secciones con datos mock.
- **Caso borde**: Tablas vacías → verificar arrays vacíos (`[]`), no `null`.
- **Caso borde**: Error en un repositorio → verificar que retorna error 500 y no respuesta parcial.
- **Caso borde**: `CustomSection` con `Visible = false` → verificar que no se incluye en la respuesta.
- **Caso borde**: `CustomSection` con `Visible = true` → verificar que se incluye en la respuesta.

### 7.2 Tests de Integración

- **Mock de repositorios**: Usar SQLite en memoria (`:memory:`) con datos precargados o mocks para simular datos en las 5 tablas y verificar la respuesta agregada. Se descarta `testcontainers-go` porque SQLite es embebida y no requiere Docker.
- **Retrocompatibilidad**: Verificar que la respuesta JSON incluye todos los campos antiguos con los mismos nombres y tipos.

### 7.3 Cobertura

El umbral mínimo es 85% sobre el archivo completo `public_service.go` (incluyendo líneas existentes y nuevas). Se recomienda crear tests parametrizados que cubran tanto los flujos preexistentes como los nuevos.

---

## 8. Archivos Modificados

| Archivo | Cambio | Tipo |
|---|---|---|
| `internal/application/dto/response/info_page.go` | Agregar 5 nuevos campos al struct `InfoPageResponse` | Código |
| `internal/application/service/public_service.go` | Agregar lógica de consulta y mapeo en `GetInfoPage()` | Código |
| `api/openapi.yaml` | Extender schema `InfoPageResponse` con 5 nuevos campos | Contrato |

### 8.1 Archivos que NO requieren cambios (verificados en disco)

- `internal/application/dto/response/course.go` ✅ (existe)
- `internal/application/dto/response/certification.go` ✅ (existe)
- `internal/application/dto/response/language_resp.go` ✅ (existe)
- `internal/application/dto/response/reference.go` ✅ (existe)
- `internal/application/dto/response/custom_section.go` ✅ (existe)
- `internal/application/port/input/public_usecase.go` ✅ (interfaz no cambia)
- `internal/application/port/output/repository_port.go` ✅ (interfaces existen)
- `internal/infrastructure/adapter/repository/sqlite_repository.go` ✅ (implementaciones existen)
- `internal/infrastructure/adapter/http/handler/public_handler.go` ✅ (handler no cambia)
- `internal/infrastructure/adapter/http/router.go` ✅ (rutas ya están registradas)
- `internal/infrastructure/adapter/http/middleware/*` ✅ (sin cambios)

---

## 9. Criterios de Aceptación

| # | Criterio | Verificación |
|---|---|---|
| CA1 | `GET /v1/ms-resume/public/info-page` retorna `courses: []CourseResponse` | La respuesta JSON contiene el campo `courses` con la misma estructura que `GET /public/courses` |
| CA2 | Retorna `certifications: []CertificationResponse` | Idéntico a `GET /public/certifications` |
| CA3 | Retorna `languages: []LanguageResponse` | Idéntico a `GET /public/languages-data` |
| CA4 | Retorna `references: []ReferenceResponse` | Idéntico a `GET /public/references` |
| CA5 | Retorna `customSections: []CustomSectionResponse` solo con `visible == true` | Idéntico a `GET /public/custom-sections` |
| CA6 | Los campos existentes (`home`, `basicData`, `skills`, `experiences`, `educations`, `altchaChallenge`) no cambian de nombre, tipo ni posición | Compatibilidad con frontends actuales |
| CA7 | Si no hay datos en una tabla, se retorna `[]` (no `null`, no se omite) | Respuesta siempre incluye el campo con array |
| CA8 | El endpoint público sigue funcionando sin autenticación | `curl` sin token funciona |
| CA9 | Si falla la consulta de UNA entidad, se retorna error 500 (todo o nada) | Respuesta parcial no es posible |
| CA10 | El OpenAPI `api/openapi.yaml` refleja los nuevos campos en `InfoPageResponse` | Schema actualizado |

---

## 10. Edge Cases

| # | Escenario | Comportamiento Esperado |
|---|---|---|
| EC1 | No hay courses registrados | `courses: []` |
| EC2 | No hay certifications registradas | `certifications: []` |
| EC3 | No hay languages registrados | `languages: []` |
| EC4 | No hay references registradas | `references: []` |
| EC5 | No hay customSections o todas invisibles | `customSections: []` |
| EC6 | Todas las tablas vacías | Arrays vacíos + datos existentes + AltchaChallenge |
| EC7 | Error en repositorio de alguna entidad | Error 500 sin respuesta parcial |
| EC8 | Concurrencia | Sin riesgo adicional. Consultas secuenciales, cada una en una transacción |

---

## 11. Deuda Técnica

No se introduce deuda técnica. Los repositorios, DTOs y mappers ya existen del incremento `migration-to-go`.

---

## 12. Impacto en Artefactos Existentes

| Artefacto | Impacto | Acción |
|---|---|---|
| `master_spec` (migration-to-go.md) | Bajo. Sección 3.3 (InfoPage) debe reflejar el cambio al consolidar | Documentar en próxima consolidación |
| `api/openapi.yaml` | Se extiende schema `InfoPageResponse` | Modificar directamente |
| `docs/specs/increments/migration-to-go-decomposition.md` | Ninguno. Fase 6 ya incluía InfoPage | Sin cambios |
| Task Board `migration-to-go` | Ninguno (está `implemented`) | Sin cambios |
| `PdfCache` (hash de invalidación) | **Potencial**: Si el hash no incluye las 5 nuevas entidades, el PDF quedaría desactualizado al cambiar cursos, certificaciones, etc. | Verificar que `PdfCache` incluya courses, certifications, languages, references, customSections en su hash de invalidación. Si no es así, actualizar el hash como parte del incremento. |

---

## 13. Cobertura de Validación Post-Implementación

1. `go build` exitoso
2. `go vet` sin errores
3. Tests unitarios de `PublicService.GetInfoPage` pasan
4. OpenAPI válido (sin errores de sintaxis ni referencias rotas)
5. Prueba manual: `curl GET /v1/ms-resume/public/info-page` retorna los 11 campos
6. Verificar que `PdfCache` incluya las 5 nuevas entidades en su hash de invalidación para evitar PDF desactualizado
