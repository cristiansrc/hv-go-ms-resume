# Requirements Brief - Extender InfoPage con Nuevas Entidades

**Status**: `requirements-discovery`
**Fecha**: 2026-06-04
**Owner**: cristiansrc
**Proyecto**: `hv-go-ms-resume`
**Incremento**: `extend-info-page-new-entities`

---

## 1. Objetivo

Incluir las 5 entidades creadas en la migración (Course, Certification, Language, Reference, CustomSection) en la respuesta agregada del endpoint `GET /v1/ms-resume/public/info-page`, de modo que los frontends puedan obtener todos los datos públicos del CV/portfolio en una sola llamada.

Actualmente estas entidades solo se exponen mediante endpoints públicos individuales (`/public/courses`, `/public/certifications`, `/public/languages-data`, `/public/references`, `/public/custom-sections`), pero no forman parte de la respuesta consolidada de `info-page`.

---

## 2. Contexto

El proyecto `hv-go-ms-resume` se migró de Spring Boot a Go, y durante esa migración se incorporaron 5 nuevas entidades (Course, Certification, Language, Reference, CustomSection) con sus correspondientes CRUD, repositorios y endpoints públicos individuales. Sin embargo, el endpoint `/public/info-page` —que funge como respuesta consolidada para el portal público— no se actualizó para incluirlas.

El portal público (`hv-rt-fr-portal`) consume `info-page` para renderizar la página principal del portfolio. Al no estar estas entidades en la respuesta, el frontend debe hacer llamadas adicionales a los endpoints individuales, lo que incrementa la latencia y complejidad del cliente.

**Respuesta actual de `info-page`** (solo 6 secciones):

| Campo | Tipo | Origen |
|---|---|---|
| `home` | `HomeResponse` | Entidad Home (id=1) |
| `basicData` | `BasicDataResponse` | Entidad BasicData (id=1) |
| `skills` | `[]SkillResponse` | Skills con SkillSons anidados |
| `experiences` | `[]ExperienceResponse` | Lista completa |
| `educations` | `[]EducationResponse` | Lista completa |
| `altchaChallenge` | `AltchaChallengeResponse` | Generado por request |

---

## 3. Actores y Permisos

| Actor | Rol | Permisos |
|---|---|---|
| **Visitante público** | Usuario anónimo | Leer todos los datos del info-page (incluyendo las nuevas entidades) sin autenticación |
| **Frontend Portal (`hv-rt-fr-portal`)** | Cliente del API | Consumir el endpoint público `GET /v1/ms-resume/public/info-page` |

No hay cambios en los permisos existentes. El endpoint `/public/info-page` permanece público (sin JWT).

---

## 4. Alcance (Scope)

### 4.1 Incluye

1. **Agregar `courses` al `InfoPageResponse`**: Lista de `CourseResponse[]` con todos los cursos activos ordenados por `order ASC`.
2. **Agregar `certifications` al `InfoPageResponse`**: Lista de `CertificationResponse[]` con todas las certificaciones activas ordenadas por `order ASC`.
3. **Agregar `languages` al `InfoPageResponse`**: Lista de `LanguageResponse[]` con todos los idiomas activos ordenados por `order ASC`.
4. **Agregar `references` al `InfoPageResponse`**: Lista de `ReferenceResponse[]` con todas las referencias activas ordenadas por `order ASC`.
5. **Agregar `customSections` al `InfoPageResponse`**: Lista de `CustomSectionResponse[]` con **solo las secciones visibles** (`visible == true`) activas ordenadas por `order ASC`.
6. **Actualizar el OpenAPI** (`api/openapi.yaml`): Modificar el schema `InfoPageResponse` para reflejar los nuevos campos.

### 4.2 No Incluye (Out of Scope)

- **No crear nuevas entidades**: Las entidades ya existen y tienen sus CRUD completos.
- **No modificar los endpoints públicos individuales**: `/public/courses`, `/public/certifications`, `/public/languages-data`, `/public/references`, `/public/custom-sections` continúan funcionando igual.
- **No modificar la estructura de las entidades**: Los schemas `CourseResponse`, `CertificationResponse`, `LanguageResponse`, `ReferenceResponse`, `CustomSectionResponse` se mantienen exactamente como están.
- **No afectar la generación de PDF**: El hash de caché de PDFs actualmente **no** incluye estas 5 entidades en su hash de invalidación (ver `collectCvData` en `pdf_service.go`). Fuera del scope de este incremento, pero se debe verificar que el `PdfCache` se actualice para incluirlas (ver sección 12 de la Delta Spec).
- **No modificar la respuesta de error** del endpoint.

---

## 5. Flujos de Usuario

### 5.1 Flujo Principal: Visitante carga el portal

1. El visitante público accede al portal (`cristiansrc.com`).
2. El frontend (`hv-rt-fr-portal`) hace `GET /v1/ms-resume/public/info-page`.
3. El servicio `PublicService.GetInfoPage` ahora además de los datos actuales:
   - Obtiene todos los **courses** activos desde `CourseRepository.List()`.
   - Obtiene todas las **certifications** activas desde `CertificationRepository.List()`.
   - Obtiene todos los **languages** activos desde `LanguageRepository.List()`.
   - Obtiene todas las **references** activas desde `ReferenceRepository.List()`.
   - Obtiene todas las **customSections** activas y visibles desde `CustomSectionRepository.List()`, filtrando `Visible == true`.
4. El frontend recibe la respuesta completa con las 11 secciones y renderiza la página sin llamadas adicionales.

### 5.2 Flujo Alterno: Dato no disponible

Si alguna de las tablas está vacía (ej: no hay cursos registrados), la lista correspondiente se retorna como `[]` (array vacío), no como `null` ni omitida.

---

## 6. Entidades Funcionales

### 6.1 Nuevas secciones en InfoPageResponse

| Sección (JSON) | Tipo DTO | Origen | Filtro |
|---|---|---|---|
| `courses` | `[]CourseResponse` | `CourseRepository.List()` | `deleted_at IS NULL`, orden `order ASC` |
| `certifications` | `[]CertificationResponse` | `CertificationRepository.List()` | `deleted_at IS NULL`, orden `order ASC` |
| `languages` | `[]LanguageResponse` | `LanguageRepository.List()` | `deleted_at IS NULL`, orden `order ASC` |
| `references` | `[]ReferenceResponse` | `ReferenceRepository.List()` | `deleted_at IS NULL`, orden `order ASC` |
| `customSections` | `[]CustomSectionResponse` | `CustomSectionRepository.List()` | `deleted_at IS NULL AND visible == true`, orden `order ASC` |

### 6.2 Campos existentes de cada DTO (sin cambios)

Los DTOs de respuesta ya existen y no se modifican:

- **`CourseResponse`**: id, name, nameEng, institution, institutionEng, completionDate, description, descriptionEng, summaryPdf, summaryPdfEng, certificateUrl
- **`CertificationResponse`**: id, name, nameEng, issuingOrganization, issuingOrganizationEng, issueDate, expirationDate, verificationUrl, credentialId, description, descriptionEng, summaryPdf, summaryPdfEng
- **`LanguageResponse`**: id, language, languageEng, readingLevel, writingLevel, speakingLevel
- **`ReferenceResponse`**: id, fullName, position, company, companyEng, email, phone, relationship, relationshipEng
- **`CustomSectionResponse`**: id, title, titleEng, content, contentEng, summaryPdf, summaryPdfEng, visible

---

## 7. Integraciones

No se requieren nuevas integraciones externas. Todas las entidades se obtienen desde la base de datos SQLite local.

---

## 8. Seguridad y Restricciones

| Restricción | Descripción |
|---|---|
| **Endpoint público** | `/public/info-page` no requiere autenticación. Se mantiene igual. |
| **Solo datos activos** | Todas las listas excluyen registros con soft delete (`deleted_at IS NOT NULL`). |
| **CustomSections visibles** | Solo se incluyen secciones con `visible == true`. El administrador controla esto vía CRUD. |
| **Sin datos sensibles adicionales** | Los DTOs actuales ya excluyen datos sensibles (timestamps internos, deleted_at, etc.). No se agrega nueva información sensible. |

---

## 9. Edge Cases

| # | Escenario | Comportamiento Esperado |
|---|---|---|
| EC1 | **No hay courses registrados** | `courses: []` (array vacío) |
| EC2 | **No hay certifications registradas** | `certifications: []` |
| EC3 | **No hay languages registrados** | `languages: []` |
| EC4 | **No hay references registradas** | `references: []` |
| EC5 | **No hay customSections o todas invisibles** | `customSections: []` |
| EC6 | **Todas las tablas vacías** | La respuesta incluye los campos con arrays vacíos, más los datos existentes (Home, BasicData, etc.) y el AltchaChallenge |
| EC7 | **Error en repositorio de alguna entidad** | El servicio debe manejar el error de forma granular: si falla la obtención de una entidad, debe retornar error 500 sin respuesta parcial. Todas las secciones se cargan o ninguna (transaccional en el servicio). |
| EC8 | **Concurrencia** | No hay riesgo de concurrencia distinto al actual. Cada lista se obtiene en una sola consulta. |

---

## 10. Criterios de Aceptación

| # | Criterio | Verificación |
|---|---|---|
| CA1 | La respuesta de `GET /v1/ms-resume/public/info-page` incluye `courses: CourseResponse[]` | Response JSON contiene el campo `courses` con la misma estructura que `GET /public/courses` |
| CA2 | La respuesta incluye `certifications: CertificationResponse[]` | Idéntico a `GET /public/certifications` |
| CA3 | La respuesta incluye `languages: LanguageResponse[]` | Idéntico a `GET /public/languages-data` |
| CA4 | La respuesta incluye `references: ReferenceResponse[]` | Idéntico a `GET /public/references` |
| CA5 | La respuesta incluye `customSections: CustomSectionResponse[]` con solo las secciones visibles | Idéntico a `GET /public/custom-sections` |
| CA6 | Los campos existentes (`home`, `basicData`, `skills`, `experiences`, `educations`, `altchaChallenge`) no se ven afectados | La respuesta sigue siendo compatible con los frontends actuales |
| CA7 | Si no hay datos en una tabla, se retorna `[]` (no `null`, no se omite el campo) | Respuesta siempre incluye el campo con array |
| CA8 | El endpoint público sigue funcionando sin autenticación | `curl` sin token funciona correctamente |
| CA9 | El OpenAPI se actualiza reflejando los nuevos campos en `InfoPageResponse` | `api/openapi.yaml` contiene los nuevos schemas |

---

## 11. Preguntas Abiertas

### Críticas (bloquean handoff a Planner)

_No hay preguntas críticas abiertas._

### No Críticas (Planner puede decidir)

| # | Pregunta | Impacto |
|---|---|---|
| PA1 | ¿Se debe aplicar algún tratamiento especial en el mapeo de `CustomSection` (ej: transformación de contenido HTML) o se mapea directamente? | Bajo. Los DTOs actuales ya resuelven esto. |
| PA2 | ¿Los nuevos campos en `InfoPageResponse` deben ir con `omitempty` o siempre presentes? | Afecta el contrato JSON. Se recomienda arrays siempre presentes (sin `omitempty`) para consistencia. |

---

## 12. Supuestos

| # | Supuesto | Validación |
|---|---|---|
| S1 | Los repositorios `CourseRepository`, `CertificationRepository`, `LanguageRepository`, `ReferenceRepository`, `CustomSectionRepository` ya están implementados y funcionando correctamente. | Verificar tests existentes. |
| S2 | Los métodos `List()` de cada repositorio ya excluyen registros con soft delete (`deleted_at IS NULL`). | Verificar queries SQL. |
| S3 | Los DTOs de respuesta (`CourseResponse`, `CertificationResponse`, `LanguageResponse`, `ReferenceResponse`, `CustomSectionResponse`) están correctamente definidos y son consumibles por los frontends. | Verificar OpenAPI actual y frontends. |
| S4 | La estructura de carpetas y la inyección de dependencias de `PublicService` permite agregar nuevas dependencias de repositorio sin refactor mayor. | Verificar constructor de `PublicService`. |
| S5 | Los frontends existentes que consumen `info-page` ignoran campos desconocidos en la respuesta JSON. | Comportamiento estándar de JSON en clientes REST. |

---

## 13. Handoff para Planner

### Resumen Funcional

Extender `InfoPageResponse` para incluir 5 nuevas secciones: `courses`, `certifications`, `languages`, `references`, `customSections`. Cada sección se obtiene del repositorio correspondiente (ya existente) y se mapea con el DTO de respuesta existente. La lógica de `PublicService.GetInfoPage()` debe inyectar los 5 nuevos repositorios y agregar las secciones a la respuesta.

### Scope

- **Incluye**: Modificar `InfoPageResponse` (struct Go), modificar `PublicService.GetInfoPage()` (lógica de agregación), actualizar OpenAPI. **Constructor `PublicService`**: Sin cambios (los 5 repos ya están inyectados).
- **Excluye**: CRUD, nuevos endpoints, cambios en DTOs existentes, cambios en PDF cache, cambios en seguridad.

### Archivos a modificar (identificados)

| Archivo | Cambio |
|---|---|
| `internal/application/dto/response/info_page.go` | Agregar campos `Courses`, `Certifications`, `Languages`, `References`, `CustomSections` al struct `InfoPageResponse` |
| `internal/application/port/input/public_usecase.go` | Sin cambios (la interfaz `GetInfoPage` retorna `*InfoPageResponse`) |
| `internal/application/service/public_service.go` | Agregar lógica de consulta y mapeo en `GetInfoPage()` (los 5 repositorios ya están inyectados en el constructor) |
| `api/openapi.yaml` | Actualizar schema `InfoPageResponse` y description del endpoint |

### Archivos que NO requieren cambios (ya existen)

- `internal/application/dto/response/course.go` ✅
- `internal/application/dto/response/certification.go` ✅
- `internal/application/dto/response/language_resp.go` ✅
- `internal/application/dto/response/reference.go` ✅
- `internal/application/dto/response/custom_section.go` ✅
- `internal/application/port/output/repository_port.go` (ya tiene las interfaces) ✅
- `internal/infrastructure/adapter/repository/sqlite_repository.go` (ya tiene las implementaciones) ✅
- `internal/infrastructure/adapter/http/handler/public_handler.go` (el handler no cambia) ✅

### Criterios de Aceptación Clave

- Retrocompatibilidad: los campos existentes de `InfoPageResponse` no cambian de nombre, tipo ni posición.
- Nuevos campos siempre presentes como arrays (incluso vacíos).
- Mismo orden que los endpoints públicos individuales (`order ASC` para todos, filtro `visible == true` para CustomSection).

### Decisiones Pendientes

| Decisión | Opciones |
|---|---|
| `omitempty` en arrays | Se recomienda **no usar `omitempty`** para que los arrays siempre aparezcan en la respuesta aunque estén vacíos. Planner debe validar que los frontends existentes no dependan de `omitempty` en campos de array (riesgo bajo). |
| Estrategia de errores | Si falla la consulta de una entidad, ¿se retorna error 500 o se omite esa sección? Se recomienda error 500 para consistencia (todo o nada). |
