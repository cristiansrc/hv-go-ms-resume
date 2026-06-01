# Shared Context - migration-to-go

**Increment**: `migration-to-go`
**Status**: `awaiting-human-plan-approval`
**Created**: 2026-05-31
**Last Updated**: 2026-05-31

---

## Current status

`awaiting-human-plan-approval` - Spec Validator ha emitido `verdict: ready`. Pendiente aprobación humana del plan.

---

## Canonical artifacts

| Artifact | Path | Status |
|---|---|---|
| Master Spec | `docs/specs/increments/migration-to-go.md` | `awaiting-human-plan-approval` |
| OpenAPI Contract | `api/openapi.yaml` | `awaiting-human-plan-approval` |
| Migration Contract | `docs/specs/increments/migration-to-go-migration.md` | `awaiting-human-plan-approval` |
| Requirements Brief | `docs/specs/requirements/migration-to-go-requirements-brief.md` | `ready-for-planner` |
| ADR-001 | `../../docs/architecture/decision-records/ADR-001-migration-to-go.md` | `accepted` |
| System Landscape | `../../docs/architecture/system-landscape.md` | `active` |
| Integration Map | `../../docs/architecture/integration-map.md` | `active` |
| Context Map | `../../docs/architecture/context-map.md` | `active` |
| RenderCV OpenAPI | `../../hv-py-ms-render-cv/docs/api/openapi.yaml` | `active` |
| Spring Boot OpenAPI (ref) | `/home/cristiansrc/Documentos/Proyectos/ms-resume/src/main/resources/openapi.yml` | `reference` |
| Graph Report | `../../graphify-out/GRAPH_REPORT.md` | `active` |

---

## Artifact evidence

### OpenAPI Contract (`api/openapi.yaml`)
- **Version**: OpenAPI 3.1.0
- **Endpoints públicos**: 16 (9 existentes + 2 nuevos templates/languages + 5 nuevos entidades públicas de solo lectura)
- **Endpoints protegidos**: 82 (CRUD para 18 entidades: BasicData y Home con 2 ops, ImageUrl/VideoUrl con 4 ops, Label con 5 ops (incluye PUT), resto con 5 ops c/u)
- **Schemas de respuesta 200**: Idénticos al OpenAPI Spring Boot para entidades existentes
- **Schemas de error**: ApiErrorResponse estandarizado (timestamp, status, error, code, message, path, trace_id, details)
- **Security**: BearerAuth JWT global, override `security: []` en endpoints públicos
- **Rate limiting**: Documentado en responses 429

### Master Spec (`docs/specs/increments/migration-to-go.md`)
- **Stack**: Go 1.22+, Chi, modernc.org/sqlite, golang-migrate, golang-jwt, go-playground/validator, slog, disintegration/imaging, golang.org/x/time/rate
- **Arquitectura**: Hexagonal (domain → application → infrastructure)
- **Entidades**: 13 existentes + 5 nuevas = 18 total
- **Tablas**: 18 entidades + 4 join tables + 2 soporte = 24 tablas
- **Integraciones**: RenderCV (HTTP), S3 (SDK v2), Telegram (HTTP)
- **Caché PDF**: SHA-256 hash + S3 storage + SQLite metadata
- **Rate limiting**: 100 req/min público, 5 req/min login, 200 req/min autenticado
- **Image optimization**: Redimensionar + eliminar EXIF (disintegration/imaging)

### Migration Contract (`docs/specs/increments/migration-to-go-migration.md`)
- **Migraciones**: 001_initial_schema + 002_new_entities
- **Script**: `go run ./cmd/migrate --source old.db --target new.db`
- **Preservación**: IDs, FKs, relaciones many-to-many
- **Order auto-asignado**: ROW_NUMBER() por entidad
- **Validación**: COUNT(*) comparison post-migración

---

## Spec Validator Approval

verdict: ready
reviewed_at: 2026-05-31T23:00:00Z
validator_agent: spec-validator
artifact_set_reviewed:
  - /home/cristiansrc/Documentos/Proyectos/hv-workspace/projects/hv-go-ms-resume/docs/specs/increments/migration-to-go.md
  - /home/cristiansrc/Documentos/Proyectos/hv-workspace/projects/hv-go-ms-resume/api/openapi.yaml
  - /home/cristiansrc/Documentos/Proyectos/hv-workspace/projects/hv-go-ms-resume/docs/specs/increments/migration-to-go-migration.md
  - /home/cristiansrc/Documentos/Proyectos/hv-workspace/projects/hv-go-ms-resume/docs/specs/.working/migration-to-go-sdd-context.md
  - /home/cristiansrc/Documentos/Proyectos/hv-workspace/projects/hv-go-ms-resume/docs/specs/requirements/migration-to-go-requirements-brief.md
summary: Re-validación post-remediación. 12/12 hallazgos previos resueltos. 1 hallazgo nuevo (N-001, low): conteo descriptivo de ops de Label en Artifact evidence (82 real vs 81 declarado). OpenAPI, Master Spec y Migration Contract alineados. Sin blockers.
invalidated_by_changes_since: none

---

## Human Plan Approval: approved_by_user

_approved_at_: 2026-05-31
_approved_by_: cristiansrc

---

## Decisions locked

| Decisión | Valor | Fuente |
|---|---|---|
| Router | Chi | ADR-001 + Master Spec |
| SQLite driver | modernc.org/sqlite | Master Spec (pure Go, sin CGO) |
| Migraciones | golang-migrate | Master Spec |
| JWT library | golang-jwt/jwt/v5 | Master Spec |
| Validation | go-playground/validator/v10 | Master Spec |
| Logging | log/slog | Master Spec (Go 1.21+ stdlib) |
| Image processing | disintegration/imaging | PA11 = A |
| Rate limiting | golang.org/x/time/rate | PA5 = A |
| HTML sanitization | N/A (frontend) | PA12 = B |
| Soft delete | `deleted_at` TIMESTAMP | Master Spec |
| JSON arrays en SQLite | TEXT column | Master Spec |
| IDs | INTEGER auto-increment | Compatibilidad Spring Boot |
| Error format | ApiErrorResponse | Skills configurados |
| Respuestas 200 | Idénticas al OpenAPI actual | Brief + ADR-001 |

---

## Validator findings

| ID | Hallazgo | Severidad | Archivo | Descripción |
|---|---|---|---|---|
| N-001 | Conteo ops Label en Artifact evidence | `low` | `migration-to-go-sdd-context.md` L39 | La descripción agrupa Label con ImageUrl/VideoUrl como "4 ops" pero Label tiene 5 operaciones protegidas. El conteo real de endpoints protegidos es 82, no 81. Corregir texto: "Label con 5 ops, ImageUrl/VideoUrl con 4 ops" y actualizar total a 82. |

## Resolved findings

| ID | Hallazgo | Tipo | Acción | Resultado |
|---|---|---|---|---|
| H-001 | Response schemas sin `order` | `contract-drift` | `superseded` | El `order` es interno de BD, no se expone en API. Policy clarificada en Master Spec 3.2 |
| H-002 | Request schemas sin `order` | `contract-drift` | `superseded` | El `order` no se expone en API. Reordenamiento vía BD directa. Policy clarificada |
| H-003 | Rutas canónicas inexistentes | `mechanical` | `fixed` | Paths corregidos en shared context. RenderCV en `../../hv-py-ms-render-cv/docs/api/openapi.yaml`, Spring Boot ref en ruta absoluta |
| H-004 | Label sin PUT | `design-decision` | `fixed` | Agregado endpoint PUT `/v1/ms-resume/label/{id}` en OpenAPI + Master Spec |
| H-005 | visible INTEGER vs boolean | `contract-drift` | `fixed` | Documentado mapeo de conversión en Master Spec (sección custom_section) |
| H-006 | Conteos/estados inconsistentes | `mechanical` | `fixed` | Estados corregidos a `planning`, conteos: 16 públicos, 81 protegidos |
| H-007 | Concurrencia PDF ambigua | `design-decision` | `fixed` | Simplificada: operación idempotente por hash, último gana, sin locking |
| H-008 | PA2/PA4 no documentadas | `design-decision` | `superseded` | Referencias stale eliminadas del Requirements Brief. No existían como contenido real |
| H-010 | Migration contract sin .down.sql | `mechanical` | `fixed` | Agregada mención explícita de `002_new_entities.down.sql` + nota de rollback |
| H-011 | ExperienceRequest minLength opcionales | `mechanical` | `reviewed-no-change` | Comportamiento válido: minLength en opcionales evita strings vacíos |
| H-012 | Rutas canónicas con `../` | `mechanical` | `fixed` | Corregidas a `../../` desde el proyecto hacia workspace root |

---

## Open questions

| # | Pregunta | Estado |
|---|---|---|
| OQ1 | ¿El bucket de S3 ya existe y las credenciales están activas? | Pendiente verificación |
| OQ2 | ¿La URL de RenderCV está configurada en la red Docker? | Pendiente verificación |
| OQ3 | ¿El Telegram Bot Token y Chat ID están activos? | Pendiente verificación |

---

## Stale terms guard

Los siguientes términos NO deben usarse en este incremento:
- `ms-resume` (Spring Boot) → usar `hv-go-ms-resume` (Go)
- `ErrorResponse` → usar `ApiErrorResponse`
- `camelCase` en columnas SQLite → usar `snake_case`
- `@Entity`, `@Table`, `@Column` → no aplica en Go
- `JPA`, `Hibernate` → no aplica
- `Flyway` → usar `golang-migrate`
- `Spring Security` → usar JWT middleware + Altcha

---

## Next action

1. **Spec Validator**: Re-validación completada con `verdict: ready` (1 hallazgo low pendiente: N-001)
2. **Planner**: Corregir N-001 (conteo de ops Label y total protegidos en Artifact evidence) antes del handoff
3. **Awaiting Human Plan Approval**: Usuario debe aprobar el plan para que Task Decomposer proceda
4. Tras aprobación humana (`## Human Plan Approval: approved_by_user`), **Task Decomposer** crea el task board
