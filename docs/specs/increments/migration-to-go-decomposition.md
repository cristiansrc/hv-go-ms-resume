# Decomposition Contract - migration-to-go

**Status**: `superseded`
**Increment**: `migration-to-go`
**Fecha**: 2026-05-31

---

## 1. Propósito

Este contrato define los límites, dependencias y orden de ejecución para que Task Decomposer pueda crear el task board. No contiene tareas implementables; solo define la estructura de descomposición.

---

## 2. Canonical Endpoint Paths

Todos los endpoints bajo el prefijo `/v1/ms-resume`:

### Públicos (sin auth)
```
POST   /v1/ms-resume/login
GET    /v1/ms-resume/public/info-page
GET    /v1/ms-resume/public/curriculum/{language}
POST   /v1/ms-resume/public/contact
GET    /v1/ms-resume/public/challenge
GET    /v1/ms-resume/public/blog
GET    /v1/ms-resume/public/blog/{id}
GET    /v1/ms-resume/public/blog-type
GET    /v1/ms-resume/public/blog-type/{id}
GET    /v1/ms-resume/public/templates
GET    /v1/ms-resume/public/languages
GET    /v1/ms-resume/public/courses
GET    /v1/ms-resume/public/certifications
GET    /v1/ms-resume/public/languages-data
GET    /v1/ms-resume/public/references
GET    /v1/ms-resume/public/custom-sections
```

### Protegidos (JWT Bearer)
```
GET/PUT  /v1/ms-resume/basic-data/{id}
GET/PUT  /v1/ms-resume/home/{id}
GET/POST /v1/ms-resume/label
GET/DELETE /v1/ms-resume/label/{id}
GET/POST /v1/ms-resume/image-url
GET/DELETE /v1/ms-resume/image-url/{id}
GET/POST /v1/ms-resume/video-url
GET/DELETE /v1/ms-resume/video-url/{id}
GET/POST /v1/ms-resume/blog
GET/PUT/DELETE /v1/ms-resume/blog/{id}
GET/POST /v1/ms-resume/blog-type
GET/PUT/DELETE /v1/ms-resume/blog-type/{id}
GET/POST /v1/ms-resume/skill-type
GET/PUT/DELETE /v1/ms-resume/skill-type/{id}
GET/POST /v1/ms-resume/skill
GET/PUT/DELETE /v1/ms-resume/skill/{id}
GET/POST /v1/ms-resume/skill-son
GET/PUT/DELETE /v1/ms-resume/skill-son/{id}
GET/POST /v1/ms-resume/experience
GET/PUT/DELETE /v1/ms-resume/experience/{id}
GET/POST /v1/ms-resume/education
GET/PUT/DELETE /v1/ms-resume/education/{id}
GET/POST /v1/ms-resume/futured-project
GET/PUT/DELETE /v1/ms-resume/futured-project/{id}
GET/POST /v1/ms-resume/course
GET/PUT/DELETE /v1/ms-resume/course/{id}
GET/POST /v1/ms-resume/certification
GET/PUT/DELETE /v1/ms-resume/certification/{id}
GET/POST /v1/ms-resume/language
GET/PUT/DELETE /v1/ms-resume/language/{id}
GET/POST /v1/ms-resume/reference
GET/PUT/DELETE /v1/ms-resume/reference/{id}
GET/POST /v1/ms-resume/custom-section
GET/PUT/DELETE /v1/ms-resume/custom-section/{id}
```

### Operacionales
```
GET /health
```

---

## 3. DTO/Schema Names

### Request DTOs
`LoginRequest`, `ContactRequest`, `BasicDataRequest`, `ImageUrlRequest`, `VideoUrlRequest`, `LabelRequest`, `HomeRequest`, `BlogRequest`, `BlogTypeRequest`, `SkillTypeRequest`, `SkillRequest`, `SkillSonRequest`, `ExperienceRequest`, `EducationRequest`, `FuturedProjectRequest`, `CourseRequest`, `CertificationRequest`, `LanguageRequest`, `ReferenceRequest`, `CustomSectionRequest`

### Response DTOs
`LoginResponse`, `AltchaChallengeResponse`, `TemplatesResponse`, `TemplateItem`, `LanguagesResponse`, `LanguageItem`, `BasicDataResponse`, `ImageUrlResponse`, `VideoUrlResponse`, `LabelResponse`, `HomeResponse`, `BlogResponse`, `BlogPageResponse`, `BlogTypeResponse`, `SkillTypeResponse`, `SkillResponse`, `SkillSonResponse`, `ExperienceResponse`, `EducationResponse`, `FuturedProjectResponse`, `InfoPageResponse`, `CourseResponse`, `CertificationResponse`, `LanguageResponse`, `ReferenceResponse`, `CustomSectionResponse`, `CreateResponse`, `HealthResponse`

### Error DTOs
`ApiErrorResponse`, `ApiErrorDetail`

---

## 4. DB Tables/Columns/Enums

### Tablas principales (24 total)
`basic_data`, `home`, `label`, `image_url`, `video_url`, `blog`, `blog_type`, `skill_type`, `skill`, `skill_son`, `experience`, `education`, `futured_project`, `course`, `certification`, `language`, `reference`, `custom_section`

### Tablas join (4)
`home_label`, `skill_type_skill`, `skill_skill_son`, `experience_skill_son`

### Tablas soporte (2)
`user_credentials`, `pdf_cache`

### Enums
- PDF languages: `english`, `spanish`
- PDF templates: `engineeringclassic`, `engineeringresumes`, `moderncv`, `sb2nov`
- Health status: `healthy`, `unhealthy`

---

## 5. Allowed Task Order (Fases)

Las tareas DEBE organizarse en estas fases secuenciales:

### Fase 1: Foundation
- Estructura de proyecto Go (go.mod, directorios)
- Configuración (config.go, env vars)
- Logging (slog setup)
- Health check endpoint
- Dockerfile multi-stage

### Fase 2: Database Layer
- Migraciones de esquema (001_initial_schema, 002_new_entities)
- SQLite repository base (conexión, CRUD genérico)
- Script de migración de datos

### Fase 3: Authentication & Security
- JWT middleware
- Altcha provider
- Login endpoint
- Auth middleware para endpoints protegidos
- Rate limiting middleware

### Fase 4: Existing Entities CRUD
- BasicData, Home, Label, ImageUrl, VideoUrl
- Blog, BlogType
- SkillType, Skill, SkillSon
- Experience, Education, FuturedProject
- Cada entidad: domain entity → port → service → handler → tests

### Fase 5: New Entities CRUD
- Course, Certification, Language, Reference, CustomSection
- Cada entidad: domain entity → port → service → handler → tests

### Fase 6: Public Endpoints
- InfoPage (aggregador de múltiples entidades)
- Blog público (paginado)
- BlogType público
- Nuevos endpoints públicos (courses, certifications, languages-data, references, custom-sections)
- Templates y Languages endpoints

### Fase 7: PDF Generation & Cache
- PDF service (hash calculation, RenderCV client, S3 integration)
- Curriculum endpoint con caché inteligente
- Fallbacks y error handling

### Fase 8: Contact & Telegram
- Contact endpoint
- Telegram client
- Altcha validation en contacto

### Fase 9: Image Processing
- Image upload validation (2MB limit)
- Image optimization (resize + EXIF removal)

### Fase 10: Integration & Polish
- CORS middleware
- Recovery middleware
- Logging middleware
- Request ID tracking
- Integration tests
- Documentation

---

## 6. Forbidden Stale Terms

Los siguientes términos NO deben aparecer en nombres de tareas, archivos o código:
- `ms-resume` (usar `hv-go-ms-resume`)
- `ErrorResponse` (usar `ApiErrorResponse`)
- `camelCase` en columnas SQL (usar `snake_case`)
- `JPA`, `Hibernate`, `Spring`, `@Entity`
- `Flyway` (usar `golang-migrate`)

---

## 7. Authoritative Files

| Tipo | Ruta |
|---|---|
| Master Spec | `docs/specs/increments/migration-to-go.md` |
| OpenAPI | `api/openapi.yaml` |
| Migration Contract | `docs/specs/increments/migration-to-go-migration.md` |
| Shared Context | `docs/specs/.working/migration-to-go-sdd-context.md` |
| Requirements Brief | `docs/specs/requirements/migration-to-go-requirements-brief.md` |
| ADR-001 | `../docs/architecture/decision-records/ADR-001-migration-to-go.md` |
| RenderCV OpenAPI | `../hv-py-ms-render-cv/docs/api/openapi.yaml` |

---

## 8. Dependency Graph

```
Fase 1 (Foundation)
    ↓
Fase 2 (Database)
    ↓
Fase 3 (Auth & Security)
    ↓
    ├──→ Fase 4 (Existing Entities CRUD) ──→ Fase 6 (Public Endpoints)
    ├──→ Fase 5 (New Entities CRUD) ────────→ Fase 6 (Public Endpoints)
    ↓
Fase 7 (PDF Generation) ← depende de Fase 4+5
    ↓
Fase 8 (Contact & Telegram)
    ↓
Fase 9 (Image Processing)
    ↓
Fase 10 (Integration & Polish)
```

---

## 9. Verification Steps por Fase

Cada fase debe incluir:
1. Unit tests para domain entities y services
2. Tests de handlers con mocks de ports
3. `go build` exitoso
4. `go vet` sin errores
5. `golangci-lint run` sin errores críticos
