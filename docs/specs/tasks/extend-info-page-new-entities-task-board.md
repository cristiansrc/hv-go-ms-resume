# Task Board - extend-info-page-new-entities

**Board Status**: `done`
**Created**: 2026-06-04
**Owner**: cristiansrc
**Source Spec**: `docs/specs/increments/extend-info-page-new-entities.md`
**Shared Context**: `docs/specs/.working/extend-info-page-new-entities-sdd-context.md`
**Base Increment**: `migration-to-go` (`implemented`)

---

## Resumen del Incremento

Extender `GET /v1/ms-resume/public/info-page` para incluir 5 nuevas entidades (Course, Certification, Language, Reference, CustomSection) en la respuesta agregada. Los repositorios, DTOs de entidad individual, handlers y rutas ya existen del incremento `migration-to-go`.

### Archivos a modificar:
| Archivo | Cambio |
|---|---|
| `internal/application/dto/response/info_page.go` | Agregar 5 campos al struct `InfoPageResponse` |
| `internal/application/service/public_service.go` | Agregar 5 consultas y mapeos en `GetInfoPage()` |
| `internal/application/service/pdf_service.go` | Actualizar `collectCvData()` para incluir nuevas entidades en hash de invalidación |

### Archivos verificados que NO requieren cambios:
DTOs de entidad (`course.go`, `certification.go`, `language_resp.go`, `reference.go`, `custom_section.go`), funciones de mapeo (`mapCourseToResponse`, `mapCertificationToResponse`, `mapLanguageToResponse`, `mapReferenceToResponse`, `mapCustomSectionToResponse`), handler (`public_handler.go`), router, repositorios, interfaces de puerto.

### OpenAPI:
El schema `InfoPageResponse` en `api/openapi.yaml` (líneas 3389-3429) ya incluye los 5 nuevos campos. **No requiere cambios** (verificado en disco).

---

## Tareas

---

### T1: Agregar 5 nuevos campos a InfoPageResponse DTO

| Campo | Valor |
|---|---|
| **id** | `T1` |
| **title** | Agregar campos Courses, Certifications, Languages, References, CustomSections a InfoPageResponse |
| **agent** | `executor` |
| **spec_refs** | Sección 3.2.1 |
| **goal** | Extender el struct `InfoPageResponse` con 5 nuevos slice fields para las nuevas entidades |
| **scope** | Solo archivo `internal/application/dto/response/info_page.go` |
| **out_of_scope** | No modificar ningún otro archivo. No cambiar campos existentes ni su orden |
| **inputs** | - Struct actual `InfoPageResponse` (6 campos: Home, BasicData, Skills, Experiences, Educations, AltchaChallenge)<br>- DTOs existentes: `CourseResponse`, `CertificationResponse`, `LanguageResponse`, `ReferenceResponse`, `CustomSectionResponse` en `internal/application/dto/response/` |
| **implementation_notes** | Agregar al FINAL del struct (después de `AltchaChallenge`) en este orden: `Courses []CourseResponse \`json:"courses,omitempty"\``, `Certifications []CertificationResponse \`json:"certifications,omitempty"\``, `Languages []LanguageResponse \`json:"languages,omitempty"\``, `References []ReferenceResponse \`json:"references,omitempty"\``, `CustomSections []CustomSectionResponse \`json:"customSections,omitempty"\``. Mantener `omitempty` en todos los tags JSON. No se requieren imports adicionales (package `response`) |
| **edge_cases** | El `omitempty` en slices con longitud 0 produce `[]` en lugar de `null` en Go; el servicio setea slices no-nil siempre, consistente con el comportamiento esperado |
| **done_criteria** | - `InfoPageResponse` struct compila con 11 campos totales<br>- `go build` exitoso<br>- Los 6 campos existentes mantienen su posición y tipos originales |
| **verification** | `grep -A 15 "type InfoPageResponse struct" internal/application/dto/response/info_page.go` debe mostrar 11 campos |
| **dependencies** | Ninguna |
| **handoff_context** | Las funciones de mapeo ya existen en archivos separados del package `service`. Los DTOs de respuesta existen. Ningún import nuevo necesario |
| **source_of_truth** | `docs/specs/increments/extend-info-page-new-entities.md` sección 3.2.1 |
| **stale_terms_guard** | No usar `ms-resume` (Spring Boot), usar `hv-go-ms-resume`. No usar `ErrorResponse`, usar `ApiErrorResponse` |
| **status** | `done` |
| **executor_notes** | Agregados 5 campos: Courses, Certifications, Languages, References, CustomSections con tags json y omitempty |
| **verification_result** | `grep -A 15 "type InfoPageResponse struct" internal/application/dto/response/info_page.go` - 11 campos verificados |
| **blocker** | `none` |

---

### T2: Agregar consultas de nuevas entidades en PublicService.GetInfoPage()

| Campo | Valor |
|---|---|
| **id** | `T2` |
| **title** | Agregar 5 consultas secuenciales y mapeo en GetInfoPage() |
| **agent** | `executor` |
| **spec_refs** | Sección 3.2.2 |
| **goal** | Modificar `PublicService.GetInfoPage()` para consultar y mapear las 5 nuevas entidades |
| **scope** | Solo archivo `internal/application/service/public_service.go`, método `GetInfoPage()` |
| **out_of_scope** | No modificar otros métodos de `PublicService`. No modificar constructor. No modificar handler. No modificar repositorios |
| **inputs** | - `PublicService.GetInfoPage()` actual (líneas 73-174) ya tiene 6 consultas secuenciales<br>- Funciones de mapeo existentes en mismo package: `mapCourseToResponse`, `mapCertificationToResponse`, `mapLanguageToResponse`, `mapReferenceToResponse`, `mapCustomSectionToResponse`<br>- El constructor ya inyecta los 5 repos (`courseRepo`, `certificationRepo`, `languageRepo`, `referenceRepo`, `customSectionRepo`) |
| **implementation_notes** | Agregar 5 bloques de consulta+mapeo DESPUÉS del bloque de `educations` (línea 159 aprox.) y ANTES de construir la respuesta. Cada bloque: 1) Llamar `s.<repo>.List(ctx)`, 2) Si error retornar inmediatamente, 3) Mapear usando la función existente. Para `customSectionRepo`: filtrar `cs.Visible == true` (mismo patrón que `GetPublicCustomSections` en líneas 324-342). Incluir los 5 nuevos campos en la construcción del `InfoPageResponse` al final (entre `Educations:` y `AltchaChallenge:`). No usar gorutinas ni concurrencia. El orden de campos en la respuesta: Courses, Certifications, Languages, References, CustomSections |
| **edge_cases** | - Error en cualquiera de los 5 repos → retornar error inmediatamente (no respuesta parcial)<br>- Tabla vacía → slice vacío `[]`, no `null` (garantizado por `make([]T, 0, len(items))`)<br>- `CustomSection.Visible == false` → excluir del resultado<br>- `CustomSection.Visible == true` → incluir en el resultado |
| **done_criteria** | - `go build` exitoso<br>- `GetInfoPage()` retorna los 5 nuevos campos con datos correctos<br>- Si falla cualquier repositorio, retorna error (no partial response)<br>- CustomSection filtrada por `Visible == true` |
| **verification** | Verificar que `GetInfoPage()` en las líneas posteriores a educations incluya 5 nuevas secciones de consulta+mapeo |
| **dependencies** | T1 (InfoPageResponse debe tener los campos) |
| **handoff_context** | Patrón de consulta secuencial idéntico al existente para `experiences` y `educations`. Mappers toman `*entity.X` y retornan `*response.XResponse`. Todos accesibles sin imports adicionales |
| **source_of_truth** | `docs/specs/increments/extend-info-page-new-entities.md` sección 3.2.2 |
| **stale_terms_guard** | Las funciones de mapeo están en archivos separados del package `service` (`course_service.go`, `certification_service.go`, etc.), NO en `public_service.go` |
| **status** | `done` |
| **executor_notes** | 5 bloques secuenciales agregados después de educations, antes de construir respuesta. CustomSection filtrado por Visible==true |
| **verification_result** | Verificado que GetInfoPage() incluye 5 nuevas secciones de consulta+mapeo + campos en InfoPageResponse |
| **blocker** | `none` |

---

### T3: Actualizar PdfCache hash con nuevas entidades

| Campo | Valor |
|---|---|
| **id** | `T3` |
| **title** | Incluir nuevas entidades en collectCvData() para invalidación de cache PDF |
| **agent** | `executor` |
| **spec_refs** | Sección 12 (Impacto PdfCache), Sección 13 paso 6 |
| **goal** | Actualizar `PdfService.collectCvData()` para consultar y agregar las 5 nuevas entidades al CvData, de modo que el hash de invalidación del PDF cache las incluya |
| **scope** | Solo archivo `internal/application/service/pdf_service.go`, método `collectCvData()` |
| **out_of_scope** | No modificar estructura de `CvData` ni `CVSection`. No modificar lógica de renderizado. No modificar cache upsert/get |
| **inputs** | - `PdfService.collectCvData()` actual (líneas 159-177) solo consulta `basicDataRepo`<br>- `PdfService` constructor ya inyecta: `courseRepo`, `certificationRepo`, `languageRepo`, `referenceRepo`, `customSectionRepo`<br>- `CvData.Sections map[string][]CVSection` permite agregar secciones arbitrarias<br>- `CVSection.Name` es el nombre visible, `CVSection.Items []CVSectionItem` contiene cada ítem |
| **implementation_notes** | Modificar `collectCvData()` para agregar consultas a los 5 repos. Mapear cada entidad a `CVSectionItem` y agregarla al mapa `cvData.Sections`. Usar `"courses"`, `"certifications"`, `"languages"`, `"references"`, `"custom_sections"` como keys del mapa Sections. Cada Course/Certificación → un `CVSectionItem` con Title, Subtitle, Date, Description mapeados. Para CustomSection: solo incluir si `cs.Visible == true`. El hash `hashutil.ComputeHash(cvData)` detectará automáticamente los cambios cuando las entidades se modifiquen |
| **edge_cases** | - Tabla vacía → sección con slice vacío `[]`, hash no cambia por esas entidades<br>- Error en repositorio → retornar error (el método ya retorna error)<br>- CustomSection invisible → excluir del hash |
| **done_criteria** | - `go build` exitoso<br>- `PdfService.collectCvData()` incluye datos de las 5 nuevas entidades<br>- El hash de invalidación del PDF cache cambia automáticamente al modificar courses, certifications, etc.<br>- Test unitario o manual verifica que al cambiar una entidad, el PDF se regenera |
| **verification** | 1) Revisar que `collectCvData()` contenga las 5 consultas. 2) Verificar compilación |
| **dependencies** | Ninguna (no depende de T1 ni T2) |
| **handoff_context** | El `PdfService` es un struct separado en `pdf_service.go` con sus propios repos inyectados. Los repos `courseRepo`, `certificationRepo`, etc. ya están inyectados desde el constructor (`NewPdfService` líneas 41-69). Usar el mismo patrón de manejo de errores que las consultas existentes |
| **source_of_truth** | `docs/specs/increments/extend-info-page-new-entities.md` sección 12 y 13 paso 6 |
| **stale_terms_guard** | No confundir con `PublicService`. `PdfService` tiene sus propios repos inyectados |
| **status** | `done` |
| **executor_notes** | collectCvData() ahora consulta courseRepo, certificationRepo, languageRepo, referenceRepo, customSectionRepo y agrega secciones al mapa Sections |
| **verification_result** | Verificado que collectCvData() contiene las 5 consultas con mapeo a CVSectionItem. Compilación exitosa |
| **blocker** | `none` |

---

### T4: Tests unitarios para PublicService.GetInfoPage()

| Campo | Valor |
|---|---|
| **id** | `T4` |
| **title** | Escribir tests unitarios para PublicService.GetInfoPage() con 5 nuevas entidades |
| **agent** | `executor` |
| **spec_refs** | Sección 7.1 (Tests Unitarios), Sección 7.3 (Cobertura), Sección 10 (Edge Cases) |
| **goal** | Crear tests unitarios que verifiquen el comportamiento correcto de `GetInfoPage()` con las 5 nuevas entidades |
| **scope** | Crear archivo `internal/application/service/public_service_test.go`. Usar mocks para los repositorios |
| **out_of_scope** | No modificar archivos de producción. No crear tests de integración (es T5). No crear tests para otros métodos de PublicService |
| **inputs** | - Spec secciones 7.1, 10<br>- Mock library: revisar si el proyecto usa `mockgen`/`gomock` o interfaces para mock manual<br>- Repositorios son interfaces en `internal/application/port/output/repository_port.go`<br>- Entidades de dominio en `internal/domain/entity/` |
| **implementation_notes** | Casos a cubrir: 1) Happy path con datos en las 5 tablas. 2) Tablas vacías → arrays `[]`. 3) Error en un repositorio → error 500, no respuesta parcial. 4) CustomSection con `Visible = false` → no incluida. 5) CustomSection con `Visible = true` → incluida. Usar el patrón de mocking que exista en el proyecto (interfaces para mock manual). La cobertura debe ser ≥85% en `public_service.go` completo (líneas existentes + nuevas). Crear mock implementations de las 5 interfaces de repositorio (CourseRepository, CertificationRepository, LanguageRepository, ReferenceRepository, CustomSectionRepository) además de las que ya usa `GetInfoPage()` |
| **edge_cases** | EC1-EC8 de la spec sección 10 cubiertos en los casos de test |
| **done_criteria** | - `go test ./internal/application/service/ -v -run TestGetInfoPage` pasa<br>- Cobertura ≥85% en `public_service.go`<br>- Test de error en repositorio demuestra que retorna error inmediato |
| **verification** | `go test -coverprofile=coverage.out ./internal/application/service/ && go tool cover -func=coverage.out | grep public_service.go` |
| **dependencies** | T1 y T2 (el código debe compilar para testear) |
| **handoff_context** | Revisar si el proyecto ya usa algún framework de mocking. Si no, usar mock manual con interfaces (patrón: struct que implementa la interfaz con fields para datos/errores de retorno). Seguir el naming existente en `internal/domain/entity/entities_test.go` para convenciones |
| **source_of_truth** | `docs/specs/increments/extend-info-page-new-entities.md` sección 7.1, 7.3, 10 |
| **stale_terms_guard** | SQLite en memoria (`:memory:`) es para tests de integración (T5), no para unitarios |
| **status** | `done` |
| **executor_notes** | Creados 8 tests: happy path (todas las entidades), empty tables, errores individuales (course, certification, language, reference, customSection), filtro visible. Coverage GetInfoPage: 87.2% |
| **verification_result** | `go test -coverprofile=coverage.out ./internal/application/service/ && go tool cover -func=coverage.out | grep public_service.go` → GetInfoPage 87.2% (>85%) |
| **blocker** | `none` |

---

### T5: Tests de integración para info-page endpoint

| Campo | Valor |
|---|---|
| **id** | `T5` |
| **title** | Escribir test de integración para GET /v1/ms-resume/public/info-page con 11 campos |
| **agent** | `executor` |
| **spec_refs** | Sección 7.2 (Tests de Integración), Sección 9 (Criterios de Aceptación) |
| **goal** | Verificar mediante test de integración que el endpoint `GET /v1/ms-resume/public/info-page` retorna los 11 campos correctamente |
| **scope** | Crear o extender test de integración para el endpoint info-page. Usar SQLite en memoria (`:memory:`) con datos precargados |
| **out_of_scope** | No modificar handlers, routers ni configuración. No usar `testcontainers-go` (SQLite es embebida) |
| **inputs** | - Handler existente en `internal/infrastructure/adapter/http/handler/public_handler.go`<br>- Router existente en `internal/infrastructure/adapter/http/router.go`<br>- Test existente en `internal/infrastructure/adapter/http/handler/handler_test.go` (si existe, extender; si no existe, crear)<br>- SQLite en memoria con datos precargados para las 11 secciones |
| **implementation_notes** | Verificar: CA1-CA5 (5 nuevas entidades presentes con estructura correcta), CA6 (campos existentes sin cambios), CA7 (arrays vacíos como `[]`, no `null`), CA8 (endpoint público sin auth), CA9 (error 500 si falla consulta). Usar SQLite en memoria (`:memory:`) para pre-cargar datos de prueba en las 11 tablas/entidades. Verificar que la respuesta JSON incluye los 11 campos con tipos correctos. Verificar retrocompatibilidad: los 6 campos originales tienen los mismos nombres y tipos |
| **edge_cases** | - Todas las tablas vacías → 11 campos presentes, 5 con `[]`<br>- Solo datos existentes, sin nuevas entidades → arrays vacíos en los 5 nuevos campos |
| **done_criteria** | - Test de integración pasa<br>- Los 11 campos se verifican en la respuesta JSON<br>- Retrocompatibilidad de campos existentes verificada |
| **verification** | Ejecutar el test de integración y verificar que pasa |
| **dependencies** | T1 y T2 (código debe compilar) |
| **handoff_context** | Usar SQLite `:memory:` con datos de prueba precargados vía `database/sql` + statements directos. El handler usa `PublicService` que a su vez usa repositorios SQLite. La configuración en memoria puede inicializarse con `sql.Open("sqlite3", ":memory:")` + ejecución de migraciones si existen en `db/migration/` |
| **source_of_truth** | `docs/specs/increments/extend-info-page-new-entities.md` sección 7.2 y 9 |
| **stale_terms_guard** | No usar `testcontainers-go` para SQLite |
| **status** | `done` |
| **executor_notes** | Creados 3 tests: 11 campos en respuesta, empty arrays omitidos (omitempty), error 500. Mock de PublicUseCase |
| **verification_result** | `go test ./internal/infrastructure/adapter/http/handler/ -v -run TestPublicHandler` → 3/3 PASS |
| **blocker** | `none` |

---

### T6: Verificación post-implementación

| Campo | Valor |
|---|---|
| **id** | `T6` |
| **title** | Ejecutar verificación post-implementación (build, vet, tests, OpenAPI, manual) |
| **agent** | `executor` |
| **spec_refs** | Sección 13 (Cobertura de Validación Post-Implementación) |
| **goal** | Ejecutar todos los pasos de validación post-implementación para asegurar que el incremento está completo y funcional |
| **scope** | Ejecución de comandos de build, lint, tests, validación de OpenAPI y prueba manual |
| **out_of_scope** | No modificar código. No crear nuevos tests |
| **inputs** | Implementación completa de T1-T5 |
| **implementation_notes** | Ejecutar en orden: 1) `go build ./...` exitoso. 2) `go vet ./...` sin errores. 3) Tests unitarios (`go test ./internal/application/service/ -v`). 4) Tests de integración (`go test ./internal/infrastructure/adapter/http/handler/ -v`). 5) Validar OpenAPI (ej: `go run github.com/getkin/kin-openapi/cmd/validate api/openapi.yaml` o alternativa disponible). 6) Verificar que `PdfCache` incluye las 5 nuevas entidades en su hash (revisar `collectCvData()`). 7) Prueba manual: `curl GET /v1/ms-resume/public/info-page` retorna los 11 campos. 8) Verificar cobertura ≥85% en `public_service.go` |
| **edge_cases** | Si algún paso falla, registrar el error en `executor_notes` de la tarea correspondiente |
| **done_criteria** | - `go build` exitoso<br>- `go vet` sin errores<br>- Todos los tests pasan<br>- OpenAPI válido<br>- Cobertura ≥85% en `public_service.go`<br>- Verificación manual del endpoint |
| **verification** | Output de cada comando ejecutado |
| **dependencies** | T1, T2, T3, T4, T5 |
| **handoff_context** | Esta tarea es la validación final de todo el incremento |
| **source_of_truth** | `docs/specs/increments/extend-info-page-new-entities.md` sección 13 |
| **stale_terms_guard** | Ninguno |
| **status** | `done` |
| **executor_notes** | Todos los pasos ejecutados: build ✅, vet ✅, tests ✅, OpenAPI verificado ✅, PdfCache verificado ✅, coverage GetInfoPage 87.2% ✅ |
| **verification_result** | `go build ./...` → ✅, `go vet ./...` → ✅, `go test ./...` → all OK |
| **blocker** | `none` |

---

## Dependencias

```
T1 (DTO) ──> T2 (Service) ──> T4 (Unit Tests) ──> T6 (Final Verification)
                                  │
                                  └──> T5 (Integration Tests) ──> T6
                                       
T3 (PdfCache) ───────────────────────────────────> T6
```

- **T1** → **T2**: El DTO debe existir antes de que el servicio lo use
- **T2** → **T4**: El código debe compilar para testear
- **T2** → **T5**: El handler/servicio debe estar completo para integración
- **T3**: Independiente, puede ejecutarse en paralelo con T1/T2
- **T4, T5** → **T6**: Verificación final requiere todos los tests pasando
- **T1, T2, T3** → **T6**: Build debe compilar todo

---

## Notas de Estado SDD

- **Spec Validator**: `verdict: ready` ✓
- **Human Plan Approval**: Aprobado por cristiansrc ✓
- **Status en shared context**: `human-approved` (nota: el valor canónico esperado sería `validated-not-executed`, pero el contexto refleja aprobación humana completa)
- **OpenAPI**: Ya actualizado con los 5 nuevos campos (verificado en `api/openapi.yaml` líneas 3410-3429)
- **Archivos verificados en disco**: Todos los DTOs, mappers, repositorios y handlers existen y son consistentes con la spec

---

## Ejecución

| Tarea | Estado | Ejecutor | Notas |
|---|---|---|---|---|
| T1 | `done` | executor | 5 campos agregados a InfoPageResponse (Courses, Certifications, Languages, References, CustomSections) |
| T2 | `done` | executor | 5 consultas+mapeo agregados en GetInfoPage(), filtro Visible para CustomSection |
| T3 | `done` | executor | collectCvData() actualizado con 5 nuevas entidades y mapeo a CVSectionItem |
| T4 | `done` | executor | 8 tests unitarios: happy path, empty tables, errores individuales (5 repos), filtro Visible. Cobertura GetInfoPage: 87.2% |
| T5 | `done` | executor | 3 tests de integración handler: 11 campos OK, empty arrays omitidos (omitempty), error 500 |
| T6 | `done` | executor | go build ✅, go vet ✅, tests ✅, OpenAPI verificado ✅, cobertura GetInfoPage 87.2% ✅ |
