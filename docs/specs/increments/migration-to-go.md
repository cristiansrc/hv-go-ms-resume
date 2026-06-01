# Master Spec - migration-to-go

**Status**: `awaiting-human-plan-approval`
**Fecha**: 2026-05-31
**Owner**: cristiansrc
**Proyecto**: `hv-go-ms-resume`
**Incremento**: `migration-to-go`
**Tipo**: Migración + nuevas funcionalidades

---

## 1. Objetivo

Migrar el microservicio `ms-resume` de Spring Boot 3.5 (Java 21) a Go 1.22+ manteniendo 100% de compatibilidad con el contrato API actual (endpoints, métodos HTTP, formato de respuesta 200). Adicionalmente, se incorporan 5 nuevas entidades funcionales (Course, Certification, Language, Reference, CustomSection) y 2 nuevos endpoints públicos (`/public/templates`, `/public/languages`).

**Nota sobre respuestas de error**: Los formatos de error se estandarizan según los standards del proyecto Go (ApiErrorResponse con timestamp, status, error, code, message, path, trace_id, details). NO se migra el formato ErrorResponse de Spring Boot. Solo se mantiene compatibilidad en las respuestas 200.

---

## 2. Arquitectura

### 2.1 Stack Tecnológico

| Componente | Tecnología | Versión | Justificación |
|---|---|---|---|
| Lenguaje | Go | 1.22+ | Performance, bajo consumo de memoria, startup rápido |
| Router | Chi | latest | Ligero, idiomático Go, compatible con middleware estándar |
| Base de datos | SQLite | modernc.org/sqlite | Pure Go, sin CGO, archivo local |
| Migraciones | golang-migrate | latest | Soporte SQLite, migraciones incrementales |
| JWT | golang-jwt/jwt/v5 | latest | Estándar de facto para JWT en Go |
| Altcha | altcha-org/altcha-go | latest | Librería oficial Go para Altcha |
| AWS SDK | aws-sdk-go-v2 | latest | SDK oficial v2 para S3 |
| Validación | go-playground/validator/v10 | latest | Validación de structs con tags |
| Logging | log/slog | stdlib (Go 1.21+) | Structured logging nativo |
| Testing | testing + testify | stdlib + external | Unit tests + assertions |
| HTTP Client | net/http | stdlib | Para llamadas a RenderCV y Telegram |
| Image processing | disintegration/imaging | latest | Redimensionado y eliminación de metadatos EXIF |
| HTML Sanitization | N/A | N/A | Responsabilidad del frontend (PA12 = B) |
| Rate Limiting | golang.org/x/time/rate | stdlib | Token bucket para rate limiting |

### 2.2 Estructura de Paquetes (Arquitectura Hexagonal)

```
hv-go-ms-resume/
├── cmd/
│   └── server/
│       └── main.go              # Entry point: config, wiring, start
├── internal/
│   ├── domain/
│   │   ├── entity/              # Modelos de dominio puros
│   │   │   ├── basic_data.go
│   │   │   ├── home.go
│   │   │   ├── label.go
│   │   │   ├── image_url.go
│   │   │   ├── video_url.go
│   │   │   ├── blog.go
│   │   │   ├── blog_type.go
│   │   │   ├── skill_type.go
│   │   │   ├── skill.go
│   │   │   ├── skill_son.go
│   │   │   ├── experience.go
│   │   │   ├── education.go
│   │   │   ├── futured_project.go
│   │   │   ├── course.go
│   │   │   ├── certification.go
│   │   │   ├── language.go
│   │   │   ├── reference.go
│   │   │   └── custom_section.go
│   │   └── error.go             # Sentinel errors del dominio
│   ├── application/
│   │   ├── port/
│   │   │   ├── input/           # Input ports (interfaces de use cases)
│   │   │   │   ├── basic_data_usecase.go
│   │   │   │   ├── auth_usecase.go
│   │   │   │   ├── contact_usecase.go
│   │   │   │   ├── pdf_usecase.go
│   │   │   │   └── ... (una por entidad)
│   │   │   └── output/          # Output ports (interfaces de infraestructura)
│   │   │       ├── repository_port.go
│   │   │       ├── s3_port.go
│   │   │       ├── rendercv_port.go
│   │   │       ├── telegram_port.go
│   │   │       └── altcha_port.go
│   │   ├── service/             # Implementación de use cases
│   │   │   ├── basic_data_service.go
│   │   │   ├── auth_service.go
│   │   │   ├── contact_service.go
│   │   │   ├── pdf_service.go
│   │   │   └── ... (una por entidad)
│   │   └── dto/                 # DTOs de aplicación (command/query/result)
│   │       ├── request/
│   │       └── response/
│   └── infrastructure/
│       ├── adapter/
│       │   ├── http/            # Driving adapters (handlers + router)
│       │   │   ├── handler/
│       │   │   │   ├── basic_data_handler.go
│       │   │   │   ├── auth_handler.go
│       │   │   │   ├── public_handler.go
│       │   │   │   └── ... (un handler por bounded context)
│       │   │   ├── middleware/
│       │   │   │   ├── auth_middleware.go
│       │   │   │   ├── rate_limit_middleware.go
│       │   │   │   ├── logging_middleware.go
│       │   │   │   ├── recovery_middleware.go
│       │   │   │   └── cors_middleware.go
│       │   │   ├── router.go
│       │   │   └── error_handler.go  # Global error handler
│       │   ├── repository/        # Driven adapters (persistencia)
│       │   │   ├── sqlite_repository.go
│       │   │   └── migrations/
│       │   │       ├── 001_initial_schema.up.sql
│       │   │       ├── 001_initial_schema.down.sql
│       │   │       └── 002_new_entities.up.sql
│       │   │       └── 002_new_entities.down.sql
│       │   ├── client/            # Driven adapters (servicios externos)
│       │   │   ├── s3_client.go
│       │   │   ├── rendercv_client.go
│       │   │   └── telegram_client.go
│       │   └── altcha/
│       │       └── altcha_provider.go
│       ├── config/
│       │   └── config.go          # Configuración desde env vars
│       └── mapper/
│           └── dto_mapper.go      # Mapeo entre domain ↔ DTO ↔ response
├── api/
│   └── openapi.yaml               # Contrato OpenAPI canónico
├── pkg/
│   └── hashutil/                  # Utilidad para hash determinista de PDFs
│       └── hash.go
├── internal/migration/
│   └── migrate_data.go            # Script de migración de datos Spring Boot → Go
├── Dockerfile
├── go.mod
├── go.sum
└── .env.example
```

### 2.3 Dirección de Dependencias

```
infrastructure → application → domain
     ↓                ↓
  adapters         ports
     ↓                ↓
  (Chi, SQLite,   (interfaces)
   S3, RenderCV,
   Telegram)
```

- **domain**: Cero dependencias externas. Solo tipos Go puros y sentinel errors.
- **application**: Depende de domain. Define interfaces (ports) y las implementa (services).
- **infrastructure**: Depende de application y domain. Implementa adapters concretos.

---

## 3. Modelo de Datos

### 3.1 Tablas SQLite

Todas las tablas de negocio incluyen: `created_at`, `updated_at`, `deleted_at` (soft delete).

#### 3.1.1 Tablas Existentes (migradas desde Spring Boot)

**basic_data** (1 fila máxima)
| Columna | Tipo | Nullable | Unique | Descripción |
|---|---|---|---|---|
| id | INTEGER PK | NO | AUTOINCREMENT | Identificador |
| first_name | TEXT | NO | - | Nombre principal |
| others_name | TEXT | YES | - | Nombres adicionales |
| first_surname | TEXT | NO | - | Primer apellido |
| others_surname | TEXT | YES | - | Apellidos adicionales |
| date_birth | TEXT | NO | - | Fecha nacimiento (YYYY-MM-DD) |
| located | TEXT | YES | - | Ubicación ES |
| located_eng | TEXT | YES | - | Ubicación EN |
| start_working_date | TEXT | YES | - | Fecha inicio laboral |
| greeting | TEXT | YES | - | Saludo ES |
| greeting_eng | TEXT | YES | - | Saludo EN |
| email | TEXT | NO | - | Email |
| instagram | TEXT | YES | - | Instagram |
| linkedin | TEXT | YES | - | LinkedIn |
| x | TEXT | YES | - | X/Twitter |
| github | TEXT | YES | - | GitHub |
| description | TEXT | YES | - | Descripción HTML ES |
| description_eng | TEXT | YES | - | Descripción HTML EN |
| description_pdf | TEXT | YES | - | JSON array texto plano ES |
| description_pdf_eng | TEXT | YES | - | JSON array texto plano EN |
| wrapper | TEXT | YES | - | JSON array roles ES |
| wrapper_eng | TEXT | YES | - | JSON array roles EN |
| created_at | TEXT | NO | - | Timestamp UTC |
| updated_at | TEXT | NO | - | Timestamp UTC |
| deleted_at | TEXT | YES | - | Soft delete |

**home** (1 fila máxima)
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| greeting | TEXT | NO | Saludo ES |
| greeting_eng | TEXT | NO | Saludo EN |
| image_url_id | INTEGER | YES | FK → image_url.id |
| button_work_label | TEXT | NO | Botón trabajo ES |
| button_work_label_eng | TEXT | NO | Botón trabajo EN |
| button_contact_label | TEXT | NO | Botón contacto ES |
| button_contact_label_eng | TEXT | NO | Botón contacto EN |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**label**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**image_url**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| url | TEXT | NO | URL de la imagen |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**video_url**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| url | TEXT | NO | URL del video |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**blog**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| title | TEXT | NO | Título ES |
| title_eng | TEXT | NO | Título EN |
| clean_url_title | TEXT | YES | URL slug |
| description_short | TEXT | NO | Descripción corta ES |
| description | TEXT | NO | Contenido HTML ES |
| description_short_eng | TEXT | NO | Descripción corta EN |
| description_eng | TEXT | NO | Contenido HTML EN |
| image_url_id | INTEGER | YES | FK → image_url.id |
| video_url_id | INTEGER | YES | FK → video_url.id |
| blog_type_id | INTEGER | YES | FK → blog_type.id |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**blog_type**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**skill_type**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**skill**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**skill_son**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**experience**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| year_start | TEXT | NO | Fecha inicio (YYYY-MM-DD) |
| year_end | TEXT | YES | Fecha fin (YYYY-MM-DD) |
| company | TEXT | NO | Empresa |
| location | TEXT | YES | Ubicación ES |
| location_eng | TEXT | YES | Ubicación EN |
| position | TEXT | YES | Posición ES |
| position_eng | TEXT | YES | Posición EN |
| summary | TEXT | YES | Resumen HTML ES |
| summary_eng | TEXT | YES | Resumen HTML EN |
| summary_pdf | TEXT | YES | Texto plano ES (PDF) |
| summary_pdf_eng | TEXT | YES | Texto plano EN (PDF) |
| description_items_pdf | TEXT | YES | JSON array texto ES |
| description_items_pdf_eng | TEXT | YES | JSON array texto EN |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**education**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| institution | TEXT | NO | Institución |
| area | TEXT | NO | Área ES |
| area_eng | TEXT | NO | Área EN |
| degree | TEXT | NO | Grado ES |
| degree_eng | TEXT | NO | Grado EN |
| start_date | TEXT | NO | Fecha inicio |
| end_date | TEXT | YES | Fecha fin |
| location | TEXT | NO | Ubicación ES |
| location_eng | TEXT | NO | Ubicación EN |
| highlights | TEXT | YES | JSON array ES |
| highlights_eng | TEXT | YES | JSON array EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**futured_project**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| description_short | TEXT | NO | Descripción corta ES |
| description | TEXT | NO | Descripción HTML ES |
| description_short_eng | TEXT | NO | Descripción corta EN |
| description_eng | TEXT | NO | Descripción HTML EN |
| experience_id | INTEGER | NO | FK → experience.id |
| image_list_url_id | INTEGER | YES | FK → image_url.id |
| image_url_id | INTEGER | YES | FK → image_url.id |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

#### 3.1.2 Tablas de Relaciones (many-to-many)

**experience_skill_son**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| experience_id | INTEGER | NO | FK → experience.id |
| skill_son_id | INTEGER | NO | FK → skill_son.id |

**home_label**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| home_id | INTEGER | NO | FK → home.id |
| label_id | INTEGER | NO | FK → label.id |

**skill_type_skill**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| skill_type_id | INTEGER | NO | FK → skill_type.id |
| skill_id | INTEGER | NO | FK → skill.id |

**skill_skill_son**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| skill_id | INTEGER | NO | FK → skill.id |
| skill_son_id | INTEGER | NO | FK → skill_son.id |

#### 3.1.3 Nuevas Tablas

**course**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| institution | TEXT | NO | Institución ES |
| institution_eng | TEXT | NO | Institución EN |
| completion_date | TEXT | NO | Fecha completado |
| description | TEXT | YES | HTML para portal ES |
| description_eng | TEXT | YES | HTML para portal EN |
| summary_pdf | TEXT | YES | Texto plano para PDF ES |
| summary_pdf_eng | TEXT | YES | Texto plano para PDF EN |
| certificate_url | TEXT | YES | URL certificado |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**certification**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| name | TEXT | NO | Nombre ES |
| name_eng | TEXT | NO | Nombre EN |
| issuing_organization | TEXT | NO | Organización emisora ES |
| issuing_organization_eng | TEXT | NO | Organización emisora EN |
| issue_date | TEXT | NO | Fecha emisión |
| expiration_date | TEXT | YES | Fecha expiración |
| verification_url | TEXT | YES | URL verificación |
| credential_id | TEXT | YES | ID credencial |
| description | TEXT | YES | HTML para portal ES |
| description_eng | TEXT | YES | HTML para portal EN |
| summary_pdf | TEXT | YES | Texto plano para PDF ES |
| summary_pdf_eng | TEXT | YES | Texto plano para PDF EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**language**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| language | TEXT | NO | Idioma ES |
| language_eng | TEXT | NO | Idioma EN |
| reading_level | TEXT | NO | Nivel lectura |
| writing_level | TEXT | NO | Nivel escritura |
| speaking_level | TEXT | NO | Nivel habla |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**reference**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| full_name | TEXT | NO | Nombre completo |
| position | TEXT | NO | Posición |
| company | TEXT | YES | Empresa ES |
| company_eng | TEXT | YES | Empresa EN |
| email | TEXT | YES | Email |
| phone | TEXT | YES | Teléfono |
| relationship | TEXT | YES | Relación ES |
| relationship_eng | TEXT | YES | Relación EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

**custom_section**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| title | TEXT | NO | Título ES |
| title_eng | TEXT | NO | Título EN |
| content | TEXT | YES | HTML para portal ES |
| content_eng | TEXT | YES | HTML para portal EN |
| summary_pdf | TEXT | YES | Texto plano para PDF ES |
| summary_pdf_eng | TEXT | YES | Texto plano para PDF EN |
| order | INTEGER | NO | Orden (auto-asignado) |
| visible | INTEGER | NO | 1=visible, 0=oculto |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |
| deleted_at | TEXT | YES | Soft delete |

> **Nota de mapeo `visible`**: En SQLite se almacena como `INTEGER` (0/1). En la API se expone como `boolean` (`true`/`false`). El adapter/mapper debe convertir: `1 → true`, `0 → false` en responses, y `true → 1`, `false → 0` en requests al persistir.

#### 3.1.4 Tablas de Soporte

**user_credentials** (1 fila máxima)
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| username | TEXT | NO | Username (único) |
| password_hash | TEXT | NO | Bcrypt hash |
| created_at | TEXT | NO | Timestamp UTC |
| updated_at | TEXT | NO | Timestamp UTC |

**pdf_cache**
| Columna | Tipo | Nullable | Descripción |
|---|---|---|---|
| id | INTEGER PK | NO | Identificador |
| language | TEXT | NO | english/spanish |
| template | TEXT | NO | Template ID |
| data_hash | TEXT | NO | SHA-256 hash de datos |
| s3_key | TEXT | NO | Clave en S3 |
| created_at | TEXT | NO | Timestamp UTC |
| file_size | INTEGER | YES | Tamaño en bytes |

Índice único: `(language, template)`

### 3.2 Política de Ordenamiento

Todas las respuestas de listas retornan los elementos **ya ordenados** según la columna correspondiente. El frontend no reordena. El campo `order` es **interno de la base de datos** y **no se expone** en la API (ni en requests ni en responses).

| Entidad | Criterio | Campo `order` en DB | Expuesto en API |
|---|---|---|---|
| Experience | `year_start DESC` | NO | — |
| Blog | `id DESC` (paginación) | NO | — |
| Label | `order ASC` | SÍ | NO |
| BlogType | `order ASC` | SÍ | NO |
| SkillType | `order ASC` | SÍ | NO |
| Skill | `order ASC` | SÍ | NO |
| SkillSon | `order ASC` | SÍ | NO |
| Education | `order ASC` | SÍ | NO |
| FuturedProject | `order ASC` | SÍ | NO |
| Course | `order ASC` | SÍ | NO |
| Certification | `order ASC` | SÍ | NO |
| Language | `order ASC` | SÍ | NO |
| Reference | `order ASC` | SÍ | NO |
| CustomSection | `order ASC` + `visible = 1` | SÍ | NO |
| ImageUrl | `id ASC` | NO | — |
| VideoUrl | `id ASC` | NO | — |

**Reglas de auto-asignación de `order`** (solo aplica a entidades con Campo `order` en DB = SÍ):
- En CREATE: `order = COALESCE((SELECT MAX(order) FROM table WHERE deleted_at IS NULL), 0) + 1`
- En UPDATE vía API: `order` no se expone en el request; el valor existente se mantiene.
- **Reordenamiento**: El administrador puede modificar el campo `order` directamente en la base de datos SQLite (acceso directo al archivo .db mediante herramienta SQL). El API no expone endpoints para reordenar; las respuestas siempre retornan los datos ya ordenados según el criterio definido.
- Soft-deleted entities se excluyen del cálculo de `order` y de las listas de resultados.
- **La lógica de ordenamiento es responsabilidad exclusiva del backend**. El API retorna los datos ya ordenados. El frontend no recibe ni envía el valor `order`.

### 3.3 Consistency Constraints

- `basic_data`: máximo 1 fila activa (`deleted_at IS NULL`)
- `home`: máximo 1 fila activa
- `user_credentials`: máximo 1 fila
- `pdf_cache`: único por `(language, template)`
- FKs: las relaciones many-to-many se validan al crear/actualizar
- Email en `reference`: formato email válido si se proporciona
- Phone en `reference`: formato teléfono válido si se proporciona

---

## 4. API Contract

El contrato OpenAPI canónico se encuentra en: `api/openapi.yaml`

### 4.1 Endpoints Públicos (sin autenticación)

| # | Method | Path | Status | Descripción |
|---|---|---|---|---|
| P1 | POST | `/v1/ms-resume/login` | 200/401 | Login con credenciales + Altcha |
| P2 | GET | `/v1/ms-resume/public/info-page` | 200 | Info consolidada del portal |
| P3 | GET | `/v1/ms-resume/public/curriculum/{language}` | 200/400/503 | Descargar CV en PDF |
| P4 | POST | `/v1/ms-resume/public/contact` | 200/400 | Enviar mensaje de contacto |
| P5 | GET | `/v1/ms-resume/public/challenge` | 200 | Obtener challenge Altcha |
| P6 | GET | `/v1/ms-resume/public/blog` | 200/400 | Blog paginado |
| P7 | GET | `/v1/ms-resume/public/blog/{id}` | 200/404 | Artículo individual |
| P8 | GET | `/v1/ms-resume/public/blog-type` | 200 | Todos los tipos de blog |
| P9 | GET | `/v1/ms-resume/public/blog-type/{id}` | 200/404 | Tipo de blog por ID |
| NP1 | GET | `/v1/ms-resume/public/templates` | 200 | Templates disponibles |
| NP2 | GET | `/v1/ms-resume/public/languages` | 200 | Idiomas disponibles |
| NP3 | GET | `/v1/ms-resume/public/courses` | 200 | Cursos (solo lectura) |
| NP4 | GET | `/v1/ms-resume/public/certifications` | 200 | Certificaciones (solo lectura) |
| NP5 | GET | `/v1/ms-resume/public/languages-data` | 200 | Idiomas del CV (solo lectura) |
| NP6 | GET | `/v1/ms-resume/public/references` | 200 | Referencias (solo lectura) |
| NP7 | GET | `/v1/ms-resume/public/custom-sections` | 200 | Secciones custom (solo lectura) |

### 4.2 Endpoints Protegidos (JWT Bearer)

Cada entidad existente y nueva tiene CRUD completo:

| Entidad | GET (lista) | GET (por ID) | POST | PUT | DELETE |
|---|---|---|---|---|---|
| BasicData | - | `/{id}` | - | `/{id}` (204) | - |
| Home | - | `/{id}` | - | `/{id}` (204) | - |
| Label | `/label` | `/label/{id}` | `/label` (201) | `/label/{id}` (204) | `/label/{id}` (204) |
| ImageUrl | `/image-url` | `/image-url/{id}` | `/image-url` (201) | - | `/image-url/{id}` (204) |
| VideoUrl | `/video-url` | `/video-url/{id}` | `/video-url` (201) | - | `/video-url/{id}` (204) |
| Blog | `/blog` | `/blog/{id}` | `/blog` (201) | `/blog/{id}` (204) | `/blog/{id}` (204) |
| BlogType | `/blog-type` | `/blog-type/{id}` | `/blog-type` (201) | `/blog-type/{id}` (204) | `/blog-type/{id}` (204) |
| SkillType | `/skill-type` | `/skill-type/{id}` | `/skill-type` (201) | `/skill-type/{id}` (204) | `/skill-type/{id}` (204) |
| Skill | `/skill` | `/skill/{id}` | `/skill` (201) | `/skill/{id}` (204) | `/skill/{id}` (204) |
| SkillSon | `/skill-son` | `/skill-son/{id}` | `/skill-son` (201) | `/skill-son/{id}` (204) | `/skill-son/{id}` (204) |
| Experience | `/experience` | `/experience/{id}` | `/experience` (201) | `/experience/{id}` (204) | `/experience/{id}` (204) |
| Education | `/education` | `/education/{id}` | `/education` (201) | `/education/{id}` (204) | `/education/{id}` (204) |
| FuturedProject | `/futured-project` | `/futured-project/{id}` | `/futured-project` (201) | `/futured-project/{id}` (204) | `/futured-project/{id}` (204) |
| Course | `/course` | `/course/{id}` | `/course` (201) | `/course/{id}` (204) | `/course/{id}` (204) |
| Certification | `/certification` | `/certification/{id}` | `/certification` (201) | `/certification/{id}` (204) | `/certification/{id}` (204) |
| Language | `/language` | `/language/{id}` | `/language` (201) | `/language/{id}` (204) | `/language/{id}` (204) |
| Reference | `/reference` | `/reference/{id}` | `/reference` (201) | `/reference/{id}` (204) | `/reference/{id}` (204) |
| CustomSection | `/custom-section` | `/custom-section/{id}` | `/custom-section` (201) | `/custom-section/{id}` (204) | `/custom-section/{id}` (204) |

### 4.3 Formato de Error Estandarizado

```json
{
  "timestamp": "2026-05-31T12:00:00Z",
  "status": 400,
  "error": "Bad Request",
  "code": "VALIDATION_ERROR",
  "message": "The request contains invalid fields.",
  "path": "/v1/ms-resume/basic-data/1",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "details": [
    {
      "field": "first_name",
      "code": "FIELD_REQUIRED",
      "message": "First name is required"
    }
  ]
}
```

### 4.4 Rate Limiting

- **Endpoints públicos**: 100 requests/minuto por IP
- **Login**: 5 requests/minuto por IP
- **Endpoints protegidos**: 200 requests/minuto por token
- **Respuesta 429**: `ApiErrorResponse` con `code: "RATE_LIMIT_EXCEEDED"`

---

## 5. Integraciones

### 5.1 hv-go-ms-resume → hv-py-ms-render-cv

| Propiedad | Valor |
|---|---|
| Protocolo | HTTP interno (red Docker) |
| Endpoint | `POST /render` |
| Timeout | 30 segundos |
| Retry | 1 reintento con backoff exponencial (1s, 2s) |
| Auth | Ninguna (red interna) |
| Health check | `GET /health` cada 30 segundos |

**Contrato**: Ver `api/openapi.yaml` del servicio render-cv.

**Transformación de datos**: El `pdf_service` transforma entidades del dominio al schema `CvData` de RenderCV:

```
BasicData.firstName + BasicData.firstSurName → CvData.name
BasicData.email → CvData.email
BasicData.located/locatedEng → CvData.location (según idioma)
BasicData.wrapper/wrapperEng → CvData.headline (primer elemento)
BasicData.linkedin → SocialNetwork{network: "LinkedIn", username: ...}
BasicData.github → SocialNetwork{network: "GitHub", username: ...}
Experience[] → sections.experience[]
Education[] → sections.education[]
SkillType[].skills[] → sections.skills[]
Course[] → sections.custom_sections["Courses"][]
Certification[] → sections.custom_sections["Certifications"][]
Language[] → sections.custom_sections["Languages"][]
Reference[] → sections.custom_sections["References"][]
CustomSection[].visible=1 → sections.custom_sections[title][]
```

### 5.2 hv-go-ms-resume → AWS S3

| Propiedad | Valor |
|---|---|
| Protocolo | HTTPS (AWS SDK v2) |
| Timeout | 15 segundos |
| Auth | Access Key + Secret Key (env vars) |
| Bucket | Configurado por variable de entorno |
| Prefix multimedia | `images/`, `videos/` |
| Prefix PDF cache | `pdfs/` |

**Operaciones**:
- `GetObject`: Descargar PDF cacheado o imagen
- `PutObject`: Guardar PDF cacheado
- `DeleteObject`: Eliminar PDF cacheado obsoleto

### 5.3 hv-go-ms-resume → Telegram Bot

| Propiedad | Valor |
|---|---|
| Protocolo | HTTPS (Telegram Bot API) |
| Endpoint | `POST https://api.telegram.org/bot{token}/sendMessage` |
| Timeout | 10 segundos |
| Auth | Bot Token (env var) |
| Retry | 1 reintento |
| Fallback | Si falla, se loggea pero NO retorna error al usuario |

---

## 6. Flujos de Usuario

### 6.1 Login del Administrador
1. POST `/v1/ms-resume/login` con `{user, password, altcha}`
2. Validar Altcha challenge
3. Validar credenciales contra `user_credentials` (bcrypt compare)
4. Generar JWT (24h validez)
5. Retornar `{token: string}`
6. Fallo → 401

### 6.2 CRUD Genérico
1. Validar JWT (middleware)
2. Validar request body (validator)
3. Ejecutar use case
4. Persistir en SQLite
5. Retornar respuesta (200/201/204)

### 6.3 Descarga de PDF con Caché Inteligente
1. GET `/v1/ms-resume/public/curriculum/{language}?template=xxx`
2. Calcular hash SHA-256 determinista de todas las entidades que afectan PDF
3. Consultar `pdf_cache` por `(language, template)`
4. Si hash coincide → descargar desde S3 → retornar PDF binario
5. Si hash difiere o no existe:
   a. Obtener datos de SQLite
   b. Transformar a schema RenderCV
   c. POST `/render` a hv-py-ms-render-cv
   d. Si existía PDF anterior con hash diferente → eliminar de S3
   e. Guardar nuevo PDF en S3: `pdfs/{language}/{template}/{hash}.pdf`
   f. Actualizar `pdf_cache`
   g. Retornar PDF binario

**Concurrencia**: Como el caché se valida por hash SHA-256 de los datos, si dos requests simultáneos detectan cache miss, ambos pueden regenerar el PDF. La operación es idempotente: mismo hash → mismo PDF. El último en guardar en S3 simplemente sobrescribe. No se requiere locking.

### 6.4 Formulario de Contacto
1. POST `/v1/ms-resume/public/contact` con `{name, email, message, altcha}`
2. Validar Altcha
3. Enviar notificación a Telegram
4. Si Telegram falla → loggear, NO retornar error
5. Retornar 200 sin body

---

## 7. Seguridad

| Aspecto | Implementación |
|---|---|
| Autenticación | JWT Bearer Token (HS256, secret desde env var) |
| Password | Bcrypt (cost 12) |
| Altcha | Validación en login y contacto |
| Rate Limiting | Token bucket por IP/token |
| CORS | Configurado para dominios de frontend |
| Env vars | JWT_SECRET, ADMIN_USERNAME, ADMIN_PASSWORD_HASH, AWS_ACCESS_KEY, AWS_SECRET_KEY, AWS_BUCKET, AWS_REGION, RENDER_CV_URL, TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID |
| SQLite permissions | 0600 (solo owner read/write) |
| Input validation | go-playground/validator en todos los requests |

---

## 8. Observabilidad

| Señal | Implementación |
|---|---|
| Logs | Structured logging con slog (JSON format en producción) |
| Health check | `GET /health` → `{status, timestamp, version}` |
| Request ID | Header `X-Request-ID` o generado UUID |
| Métricas | Log de duración de requests, cache hit/miss de PDFs |
| Error tracking | Log con stack trace para 5xx |

---

## 9. Configuración

Variables de entorno obligatorias:

| Variable | Requerida | Descripción |
|---|---|---|
| `PORT` | NO | Puerto del servicio (default: 8080) |
| `JWT_SECRET` | SÍ | Secret para firmar JWT |
| `ADMIN_USERNAME` | SÍ | Username del admin |
| `ADMIN_PASSWORD_HASH` | SÍ | Bcrypt hash del password |
| `JWT_EXPIRATION_HOURS` | NO | Expiración JWT (default: 24) |
| `DATABASE_PATH` | NO | Path del archivo SQLite (default: ./data/resume.db) |
| `AWS_ACCESS_KEY_ID` | SÍ | AWS Access Key |
| `AWS_SECRET_ACCESS_KEY` | SÍ | AWS Secret Key |
| `AWS_REGION` | SÍ | AWS Region |
| `AWS_BUCKET` | SÍ | S3 Bucket name |
| `RENDER_CV_URL` | SÍ | URL del servicio RenderCV |
| `TELEGRAM_BOT_TOKEN` | SÍ | Telegram Bot Token |
| `TELEGRAM_CHAT_ID` | SÍ | Telegram Chat ID |
| `LOG_LEVEL` | NO | Nivel de log (default: info) |
| `RATE_LIMIT_PUBLIC` | NO | Rate limit público req/min (default: 100) |
| `RATE_LIMIT_LOGIN` | NO | Rate limit login req/min (default: 5) |
| `RATE_LIMIT_AUTHENTICATED` | NO | Rate limit autenticado req/min (default: 200) |

---

## 10. Migración de Datos

### 10.1 Script de Migración

El script `internal/migration/migrate_data.go`:

1. Lee la BD SQLite actual (Spring Boot) desde un path configurable
2. Crea las nuevas tablas con campos adicionales (`order`, `deleted_at`, `created_at`, `updated_at`)
3. Migra los datos existentes preservando IDs y relaciones
4. Asigna valores de `order` automáticos:
   - Para Experience: `order` basado en `yearStart DESC` (el más reciente = order 1)
   - Para las demás: `order` basado en `id ASC`
5. NO migra datos a las nuevas entidades (Course, Certification, Language, Reference, CustomSection)
6. Crea la tabla `pdf_cache` vacía
7. Crea/migra la tabla `user_credentials` con el password hasheado actual
8. Crea las tablas de relaciones many-to-many y migra las relaciones existentes

### 10.2 Ejecución

El script se ejecuta como un comando separado:
```bash
go run ./cmd/migrate --source /path/to/old.db --target /path/to/new.db
```

---

## 11. Criterios de Aceptación

### CA1: Compatibilidad de Respuestas 200
- [ ] Todos los endpoints del OpenAPI actual están disponibles con mismo path, método y response 200
- [ ] Los schemas de respuesta son idénticos al OpenAPI actual
- [ ] Los frontends funcionan sin cambios de código contra la nueva API

### CA2: Autenticación
- [ ] Login con credenciales válidas + Altcha válido retorna JWT (24h)
- [ ] Login con credenciales inválidas retorna 401
- [ ] Login con Altcha inválido retorna 401
- [ ] Endpoints CRUD sin token retornan 401
- [ ] Endpoints CRUD con token expirado retornan 401
- [ ] Endpoints `/public/*` accesibles sin token

### CA3: CRUD de Entidades Existentes
- [ ] Las 13 entidades soportan CREATE, READ, UPDATE, DELETE
- [ ] Validación de campos requeridos según OpenAPI
- [ ] READ retorna datos en formato esperado
- [ ] DELETE aplica soft delete

### CA4: Nuevas Entidades
- [ ] Course, Certification, Language, Reference, CustomSection soportan CRUD
- [ ] Validación de campos requeridos
- [ ] Patrón multi-idioma (campos ES/EN)
- [ ] Separación HTML (portal) / Texto plano (PDF)
- [ ] Accesibles vía endpoints públicos de solo lectura

### CA5: Multimedia
- [ ] Registro de URL de imagen/video retorna ID
- [ ] Soft delete de registros
- [ ] Imágenes > 2MB rechazadas con error descriptivo
- [ ] Optimización de imágenes: redimensionar + eliminar metadatos EXIF

### CA6: Generación de PDF y Caché
- [ ] Descarga de CV en PDF (english/spanish)
- [ ] Selección de template
- [ ] Caché inteligente con invalidación por hash
- [ ] RenderCV no disponible → 503
- [ ] Idiomas/templates no válidos → 400
- [ ] S3 no disponible para leer → regenerar sin error
- [ ] S3 no disponible para guardar → retornar sin cachear
- [ ] Hash determinista correcto
- [ ] Clave S3: `pdfs/{language}/{template}/{hash}.pdf`

### CA7: Formulario de Contacto
- [ ] Challenge Altcha se genera correctamente
- [ ] Mensaje con challenge válido → notifica Telegram
- [ ] Challenge inválido → 400
- [ ] Telegram no disponible → no retorna error al usuario

### CA8: Base de Datos
- [ ] SQLite como archivo local
- [ ] Sin servicio de BD externo
- [ ] Migraciones aplicadas al iniciar
- [ ] Datos persisten entre reinicios
- [ ] Soft delete en todas las entidades

### CA9: Seguridad
- [ ] Password hasheado (bcrypt)
- [ ] Credenciales desde variables de entorno
- [ ] SQLite file permissions 0600

### CA10: Performance
- [ ] Startup < 2 segundos
- [ ] Memoria en reposo < 50MB
- [ ] Endpoints públicos p95 < 500ms

### CA11: Nuevos Endpoints Públicos
- [ ] `GET /public/templates` retorna templates disponibles
- [ ] `GET /public/languages` retorna idiomas disponibles

### CA12: Rate Limiting
- [ ] Rate limit en endpoints públicos (100 req/min por IP)
- [ ] Rate limit en login (5 req/min por IP)
- [ ] Rate limit en endpoints protegidos (200 req/min por token)
- [ ] Respuesta 429 con ApiErrorResponse

### CA13: Migración de Datos
- [ ] Script migra datos desde BD Spring Boot preservando IDs
- [ ] Relaciones many-to-many preservadas
- [ ] Campo `order` auto-asignado correctamente
- [ ] User credentials migradas con password hasheado
- [ ] Nuevas tablas creadas vacías

---

## 12. Decisiones Técnicas

| Decisión | Valor | Justificación |
|---|---|---|
| Router | Chi | Ligero, idiomático, middleware estándar |
| SQLite driver | modernc.org/sqlite | Pure Go, sin CGO required |
| Migraciones | golang-migrate | Soporte SQLite, estándar |
| JWT library | golang-jwt/jwt/v5 | Estándar de facto |
| Validation | go-playground/validator | Tags en structs, ampliamente usado |
| Logging | log/slog | Nativo Go 1.21+, structured |
| Image processing | disintegration/imaging | Simple, soporta redimensionado y EXIF |
| Rate limiting | golang.org/x/time/rate | Token bucket, stdlib |
| HTML sanitization | N/A (frontend) | PA12 = B: responsabilidad del frontend |
| Soft delete | `deleted_at` TIMESTAMP | Permite auditoría y recuperación |
| JSON arrays en SQLite | TEXT column | Serialización JSON, estándar en SQLite |
| IDs | INTEGER (auto-increment) | Compatibilidad con Spring Boot actual |

---

## 13. Riesgos y Mitigación

| Riesgo | Impacto | Mitigación |
|---|---|---|
| Drift de contrato API | Frontends se rompen | OpenAPI-first, tests de integración contra frontends |
| Migración de datos corrupta | Pérdida de datos | Script con validación, backup previo, dry-run mode |
| RenderCV no disponible | No se generan PDFs | Health check, 503 descriptivo, retry con backoff |
| S3 no disponible | No hay multimedia ni caché | Fallbacks documentados, errores descriptivos |
| Concurrencia en PDF cache | Duplicación de generación | Idempotente por hash, último gana |

---

## 14. Out of Scope

- Multi-usuario o gestión de roles
- Analytics o métricas de uso
- Webhooks o eventos asíncronos
- Versionado de CVs
- Notificaciones push o email (solo Telegram)
- Búsqueda full-text
- Internacionalización de la UI del admin
- Backup automático de SQLite

---

## 15. Preguntas Abiertas

### Resueltas
| # | Pregunta | Respuesta |
|---|---|---|
| PA1 | ¿Cuáles son los campos de cada entidad? | Extraídos del OpenAPI actual |
| PA3 | ¿Credenciales de admin en SQLite? | Sí, con password hasheado |
| PA5 | ¿Rate limiting? | SÍ, en esta iteración |
| PA6 | ¿Límite multimedia? | 2MB con optimización |
| PA7 | ¿Endpoints templates/languages? | SÍ, nuevos |
| PA8 | ¿HTML portal / texto PDF? | SÍ, patrón confirmado |
| PA9 | ¿Soft delete? | SÍ, universal |
| PA10 | ¿Orden de entidades? | Todas con `order` excepto Experience |
| PA11 | ¿Optimización de imágenes? | SÍ, redimensionar + eliminar EXIF |
| PA12 | ¿Sanitización HTML? | NO, responsabilidad del frontend |

### Críticas
_Ninguna abierta._

---

## 16. Handoff para Spec Validator

Artefactos a validar:
1. Esta Master Spec (`docs/specs/increments/migration-to-go.md`)
2. OpenAPI contract (`api/openapi.yaml`)
3. Migration contract (`docs/specs/increments/migration-to-go-migration.md`)
4. Shared context (`docs/specs/.working/migration-to-go-sdd-context.md`)
