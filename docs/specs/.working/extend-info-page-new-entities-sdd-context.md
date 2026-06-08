# Shared Context - extend-info-page-new-entities

**Increment**: `extend-info-page-new-entities`
**Status**: `human-approved`
**Created**: 2026-06-04
**Last Updated**: 2026-06-04

---

## Current status

`executed` - Implementación completada por Executor el 2026-06-04. Todas las tareas del task board están `done`. Build, vet y tests OK.

---

## Canonical artifacts

| Artifact | Path | Status |
|---|---|---|
| Delta Spec | `docs/specs/increments/extend-info-page-new-entities.md` | `awaiting-human-plan-approval` |
| OpenAPI Contract | `api/openapi.yaml` | `updated` |
| Requirements Brief | `docs/specs/requirements/extend-info-page-new-entities-requirements-brief.md` | `ready-for-planner` |
| Master Spec (base) | `docs/specs/increments/migration-to-go.md` | `implemented` |
| Shared Context | `docs/specs/.working/extend-info-page-new-entities-sdd-context.md` | `awaiting-human-plan-approval` |

---

## Artifact evidence

### Repository Existente (código ya implementado)
- **`InfoPageResponse` struct**: 6 campos existentes (Home, BasicData, Skills, Experiences, Educations, AltchaChallenge)
- **`PublicService.GetInfoPage()`**: Consulta 4 repos (basicData, home, experiences, educations) + skills anidadas + altcha
- **`PublicService` constructor**: Ya recibe los 5 repos nuevos (courseRepo, certificationRepo, languageRepo, referenceRepo, customSectionRepo)
- **Funciones de mapeo**: `mapCourseToResponse`, `mapCertificationToResponse`, `mapLanguageToResponse`, `mapReferenceToResponse`, `mapCustomSectionToResponse` ya existen en archivos separados del package `service` (`course_service.go`, `certification_service.go`, `language_service.go`, `reference_service.go`, `custom_section_service.go`) y son accesibles desde `public_service.go` sin imports adicionales
- **Handlers públicos**: Ya existen y están enrutados (`GetPublicCourses`, `GetPublicCertifications`, `GetPublicLanguages`, `GetPublicReferences`, `GetPublicCustomSections`)
- **DTOs de respuesta**: CourseResponse, CertificationResponse, LanguageResponse, ReferenceResponse, CustomSectionResponse existen en `internal/application/dto/response/`

### OpenAPI Contract (`api/openapi.yaml`)
- **Schemas existentes y verificados**: CourseResponse (L3476+), CertificationResponse (L3554+), LanguageResponse (L3617+), ReferenceResponse (L3667+), CustomSectionResponse (L3721+)
- **InfoPageResponse actualizado**: Ahora incluye 11 propiedades (6 originales + courses, certifications, languages, references, customSections)
- **Endpoint description actualizada**: Menciona las 11 secciones

### Conformidad con Requirements Brief
- **Req 1** (courses): Pendiente de implementar
- **Req 2** (certifications): Pendiente de implementar
- **Req 3** (languages): Pendiente de implementar
- **Req 4** (references): Pendiente de implementar
- **Req 5** (customSections): Pendiente de implementar, con filtro `visible == true`
- **Req 6** (OpenAPI): ✅ Actualizado - InfoPageResponse con 11 campos

---

## Spec Validator Approval

| Campo | Valor |
|---|---|
| verdict | `ready` |
| reviewed_at | 2026-06-04T20:00:00-06:00 |
| validator_agent | `spec-validator` |
| artifact_set_reviewed | `docs/specs/increments/extend-info-page-new-entities.md`, `docs/specs/requirements/extend-info-page-new-entities-requirements-brief.md`, `docs/specs/.working/extend-info-page-new-entities-sdd-context.md`, `api/openapi.yaml` |
| summary | Re-validación final exitosa. Los 9 hallazgos (1 Major, 1 Medium, 2 Minor, 1 Low-Medium, 2 Low, 1 Info) están verificados como resueltos en los archivos. El delta spec, requirements brief y shared context reflejan correctamente: funciones de mapeo en archivos separados del package service, constructor sin cambios requeridos, y PdfCache fuera de scope con advertencia explícita. OpenAPI actualizado con 11 campos en InfoPageResponse. |
| invalidated_by_changes_since | `none` |

---

## Human Plan Approval

| Campo | Valor |
|---|---|
| approved_by | cristiansrc (usuario) |
| approved_at | 2026-06-04 |
| status | `approved` |
| nota | Usuario confirmó: "apruebo el incremento" |

---

## Decisions locked

| Decisión | Valor | Fuente |
|---|---|---|
| Estrategia errores | Todo o nada: error 500 si falla cualquier repositorio | Requirements Brief + Delta Spec |
| Arrays vacíos | Siempre `[]` (no `null`, no omitir campo) | Requirements Brief + Delta Spec |
| Filtro CustomSection | Solo `visible == true` | Requirements Brief + Delta Spec |
| Sin cambios en constructor | Los 5 repos ya están inyectados en PublicService | Verificación de código |
| Serialización JSON | `omitempty` en tags Go, pero servicio setea slices no-nil | Consistencia con estilo existente |
| Sin cambios en dominio/infraestructura | Solo cambios en application layer | Scope definido |

---

## Validator findings

| ID | Hallazgo | Severidad | Archivo | Estado |
|---|---|---|---|---|
| F1 | Ubicación de funciones de mapeo incorrectamente referenciada | Major | `docs/specs/increments/extend-info-page-new-entities.md` | ✅ Resuelto |
| F2 | Contradicción Requirements Brief vs Delta Spec sobre constructor | Minor | `docs/specs/requirements/extend-info-page-new-entities-requirements-brief.md` | ✅ Resuelto |
| F3 | Recomendación subóptima de testing con testcontainers-go | Minor | `docs/specs/increments/extend-info-page-new-entities.md` | ✅ Resuelto |
| F4 | Falta verificación de impacto en PDF cache hash | Low | `docs/specs/increments/extend-info-page-new-entities.md` | ✅ Resuelto |
| F5 | Orden de nuevos campos en struct Go no especificado | Low | `docs/specs/increments/extend-info-page-new-entities.md` | ✅ Resuelto |
| F6 | Ausencia de tests unitarios base y cobertura ambigua | Info | `docs/specs/increments/extend-info-page-new-entities.md` | ✅ Resuelto |
| NF1 | Shared context desactualizado sobre ubicación de funciones de mapeo | Minor | `docs/specs/.working/extend-info-page-new-entities-sdd-context.md` | ✅ Resuelto |
| NF2 | Requirements Brief Scope incluye modificación innecesaria de constructor | Medium | `docs/specs/requirements/extend-info-page-new-entities-requirements-brief.md` | ✅ Resuelto |
| NF3 | Contradicción Requirements Brief vs código real sobre PdfCache | Low-Medium | `docs/specs/requirements/extend-info-page-new-entities-requirements-brief.md` | ✅ Resuelto |

---

## Resolved findings

| ID | Severidad | Descripción | Corrección Aplicada |
|---|---|---|---|
| F1 | Major | Funciones de mapeo referenciadas como si estuvieran en `public_service.go` cuando están en archivos separados del package `service` | Sección 3.2.2 actualizada con rutas exactas de cada función de mapeo |
| F2 | Minor | Requirements Brief indicaba "inyectar 5 repos" pero constructor ya los tiene | Línea 206 del Requirements Brief corregida |
| F3 | Minor | Recomendación de `testcontainers-go` para BD SQLite embebida | Cambiado a SQLite in-memory (`:memory:`) como opción principal |
| F4 | Low | No se verificaba impacto en PDF cache hash | Agregada fila en sección 12 (Impacto) y paso 6 en sección 13 (Validación Post-Implementación) |
| F5 | Low | Orden de campos nuevos en struct Go no especificado | Agregada nota explícita en sección 3.2.1: Courses, Certifications, Languages, References, CustomSections al final del struct |
| F6 | Info | Umbral de cobertura ambiguo (85% sobre líneas nuevas vs completo) | Aclarado: "85% sobre el archivo completo `public_service.go`" |
| NF1 | Minor | Shared context línea 34 decía funciones "ya existen en public_service.go" | Corregido: indica archivos separados del package `service` |
| NF2 | Medium | Requirements Brief Scope línea 197 decía "modificar constructor (inyectar 5 repos)" | Corregido: "Sin cambios (los 5 repos ya están inyectados)" |
| NF3 | Low-Medium | Requirements Brief línea 65 decía "hash PDF ya considera estas entidades" | Corregido: refleja que `collectCvData` no las incluye aún |

---

## Open questions

| # | Pregunta | Estado |
|---|---|---|
| OQ1 | ¿Se requiere probar manualmente la respuesta con frontend `hv-rt-fr-portal`? | Pendiente - depende del flujo de QA |
| OQ2 | ¿El orden de los campos en `InfoPageResponse` debe ser el mismo que en las respuestas individuales? | Sí, `order ASC` para todas. Filtro `visible == true` solo para CustomSection |

---

## Stale terms guard

Los siguientes términos NO deben usarse en este incremento:
- `ms-resume` (Spring Boot) → usar `hv-go-ms-resume` (Go)
- `ErrorResponse` → usar `ApiErrorResponse`
- `camelCase` en columnas SQLite → no aplica (sin cambios de BD)
- `@Entity`, `@Table`, `@Column` → no aplica en Go
- `JPA`, `Hibernate` → no aplica

---

## Next action

**Next action**: Task Decomposer

1. ✅ **Planner**: Delta Spec creada (`docs/specs/increments/extend-info-page-new-entities.md`)
2. ✅ **Planner**: OpenAPI actualizado (`api/openapi.yaml` - InfoPageResponse con 11 campos)
3. ✅ **Planner**: Shared Context creado (`docs/specs/.working/extend-info-page-new-entities-sdd-context.md`)
4. ✅ **Spec Validator**: Validación #1 completada - 6 hallazgos
5. ✅ **Spec Remediator**: 6 hallazgos corregidos (F1-F6)
6. ✅ **Spec Validator**: Re-validación #2 completada - 3 nuevos hallazgos
7. ✅ **Spec Remediator**: 3 hallazgos corregidos (NF1-NF3)
8. ✅ **Spec Validator (re-validación final)**: Veredicto `ready`
9. ✅ **Humano**: Plan aprobado por cristiansrc el 2026-06-04
10. **⏳ Task Decomposer**: Crear task board para implementación
