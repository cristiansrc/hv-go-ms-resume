# Task Board - migration-to-go

**Status**: `implemented`
**Incremento**: `migration-to-go`
**Fecha**: 2026-06-01
**Implementado**: 2026-06-01

---

## Fase 1: Foundation

### T1.1: Init Go module and project structure
- **Status**: `done`
- **Dependencies**: none
- **Verificación**: `go.mod` exists, directory structure matches Master Spec 2.2
- **Archivos**: `go.mod`, directories `cmd/server/`, `internal/domain/entity/`, `internal/application/port/input/`, `internal/application/port/output/`, `internal/application/service/`, `internal/application/dto/request/`, `internal/application/dto/response/`, `internal/infrastructure/adapter/http/handler/`, `internal/infrastructure/adapter/http/middleware/`, `internal/infrastructure/adapter/repository/migrations/`, `internal/infrastructure/adapter/client/`, `internal/infrastructure/adapter/altcha/`, `internal/infrastructure/config/`, `internal/infrastructure/mapper/`, `internal/migration/`, `pkg/hashutil/`

### T1.2: Configuration module
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: `config.go` reads all env vars from Master Spec section 9
- **Archivos**: `internal/infrastructure/config/config.go`

### T1.3: Logging setup
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: JSON structured logging with slog
- **Archivos**: `internal/infrastructure/config/logger.go`

### T1.4: Health check endpoint
- **Status**: `done`
- **Dependencies**: T1.2, T1.3
- **Verificación**: `GET /health` returns `{status, timestamp, version}`
- **Archivos**: `internal/infrastructure/adapter/http/handler/health_handler.go`

### T1.5: Dockerfile multi-stage
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: Multi-stage build, distroless target, ~10MB final image
- **Archivos**: `Dockerfile`

### T1.6: Env example file
- **Status**: `done`
- **Dependencies**: T1.2
- **Verificación**: All env vars from Master Spec section 9 documented
- **Archivos**: `.env.example`

---

## Fase 2: Database Layer

### T2.1: Domain entities (13 existing + 5 new)
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: All 18 entity structs with json tags, zero external dependencies
- **Archivos**: `internal/domain/entity/*.go`

### T2.2: Migration 001 - initial schema (up/down)
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: Contains all 13 existing entities + join tables + support tables
- **Archivos**: `internal/infrastructure/adapter/repository/migrations/001_initial_schema.up.sql`, `internal/infrastructure/adapter/repository/migrations/001_initial_schema.down.sql`

### T2.3: Migration 002 - new entities (up/down)
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: Contains 5 new entity tables
- **Archivos**: `internal/infrastructure/adapter/repository/migrations/002_new_entities.up.sql`, `internal/infrastructure/adapter/repository/migrations/002_new_entities.down.sql`

### T2.4: SQLite connection and migration runner
- **Status**: `done`
- **Dependencies**: T2.2, T2.3
- **Verificación**: Opens SQLite, runs migrations on startup
- **Archivos**: `internal/infrastructure/adapter/repository/sqlite_repository.go`

### T2.5: Output ports (repository interfaces)
- **Status**: `done`
- **Dependencies**: T2.1
- **Verificación**: Interface per entity with CRUD methods
- **Archivos**: `internal/application/port/output/repository_port.go`

### T2.6: Repository implementations
- **Status**: `done`
- **Dependencies**: T2.4, T2.5
- **Verificación**: Each entity has CRUD implementation using SQLite
- **Archivos**: `internal/infrastructure/adapter/repository/sqlite_repository.go` (CRUD methods for all entities)

### T2.7: Data migration script
- **Status**: `done`
- **Dependencies**: T2.6
- **Verificación**: CLI command that reads Spring Boot DB, writes to Go DB
- **Archivos**: `internal/migration/migrate_data.go`, `cmd/migrate/main.go`

---

## Fase 3: Authentication & Security

### T3.1: JWT service
- **Status**: `done`
- **Dependencies**: T2.1 (user_credentials domain entity)
- **Verificación**: HS256 signing, 24h expiration, claims validation
- **Archivos**: `internal/infrastructure/adapter/http/middleware/jwt_service.go`

### T3.2: Altcha provider
- **Status**: `done`
- **Dependencies**: T1.1
- **Verificación**: Challenge generation + validation
- **Archivos**: `internal/infrastructure/adapter/altcha/altcha_provider.go`, `internal/application/port/output/altcha_port.go`

### T3.3: Login endpoint
- **Status**: `done`
- **Dependencies**: T3.1, T3.2
- **Verificación**: POST `/v1/ms-resume/login` returns JWT on valid credentials + Altcha
- **Archivos**: `internal/infrastructure/adapter/http/handler/auth_handler.go`

### T3.4: Auth middleware
- **Status**: `done`
- **Dependencies**: T3.1
- **Verificación**: Validates JWT Bearer, sets claims in context, rejects without token
- **Archivos**: `internal/infrastructure/adapter/http/middleware/auth_middleware.go`

### T3.5: Rate limiting middleware
- **Status**: `done`
- **Dependencies**: T1.2
- **Verificación**: Token bucket per IP/token, 429 on limit
- **Archivos**: `internal/infrastructure/adapter/http/middleware/rate_limit_middleware.go`

---

## Fase 4: Existing Entities CRUD

### T4.1: BasicData CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4
- **Archivos**: `internal/application/port/input/basic_data_usecase.go`, `internal/application/service/basic_data_service.go`, `internal/infrastructure/adapter/http/handler/basic_data_handler.go`, `internal/application/dto/request/basic_data.go`, `internal/application/dto/response/basic_data.go`

### T4.2: Home CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4
- **Archivos**: `internal/application/port/input/home_usecase.go`, `internal/application/service/home_service.go`, `internal/infrastructure/adapter/http/handler/home_handler.go`, `internal/application/dto/request/home.go`, `internal/application/dto/response/home.go`

### T4.3: Label CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4
- **Archivos**: Similar pattern

### T4.4: ImageUrl CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T4.5: VideoUrl CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T4.6: Blog & BlogType CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T4.7: SkillType, Skill, SkillSon CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T4.8: Experience CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T4.9: Education CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T4.10: FuturedProject CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

---

## Fase 5: New Entities CRUD

### T5.1: Course CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T5.2: Certification CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T5.3: Language CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T5.4: Reference CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

### T5.5: CustomSection CRUD
- **Status**: `done`
- **Dependencies**: T2.6, T3.4

---

## Fase 6: Public Endpoints

### T6.1: InfoPage aggregated endpoint
- **Status**: `done`
- **Dependencies**: T4.x, T5.x

### T6.2: Public Blog endpoints (paginated)
- **Status**: `done`
- **Dependencies**: T4.6

### T6.3: Public BlogType endpoints
- **Status**: `done`
- **Dependencies**: T4.6

### T6.4: New public read-only endpoints (courses, certifications, languages-data, references, custom-sections)
- **Status**: `done`
- **Dependencies**: T5.x

### T6.5: Templates and Languages endpoints
- **Status**: `done`
- **Dependencies**: T1.1

---

## Fase 7: PDF Generation & Cache

### T7.1: S3 client adapter
- **Status**: `done`
- **Dependencies**: T1.2
- **Archivos**: `internal/infrastructure/adapter/client/s3_client.go`, `internal/application/port/output/s3_port.go`

### T7.2: RenderCV client adapter
- **Status**: `done`
- **Dependencies**: T1.2
- **Archivos**: `internal/infrastructure/adapter/client/rendercv_client.go`, `internal/application/port/output/rendercv_port.go`

### T7.3: Hash utility
- **Status**: `done`
- **Dependencies**: T1.1
- **Archivos**: `pkg/hashutil/hash.go`

### T7.4: PDF service
- **Status**: `done`
- **Dependencies**: T7.1, T7.2, T7.3, T4.x, T5.x

### T7.5: Curriculum endpoint
- **Status**: `done`
- **Dependencies**: T7.4

---

## Fase 8: Contact & Telegram

### T8.1: Telegram client
- **Status**: `done`
- **Dependencies**: T1.2
- **Archivos**: `internal/infrastructure/adapter/client/telegram_client.go`, `internal/application/port/output/telegram_port.go`

### T8.2: Contact endpoint with Altcha
- **Status**: `done`
- **Dependencies**: T3.2, T8.1

---

## Fase 9: Image Processing

### T9.1: Image upload validation (2MB limit + optimization)
- **Status**: `done`
- **Dependencies**: T1.1

---

## Fase 10: Integration & Polish

### T10.1: CORS middleware
- **Status**: `done`
- **Dependencies**: T1.1

### T10.2: Recovery middleware
- **Status**: `done`
- **Dependencies**: T1.1

### T10.3: Logging middleware
- **Status**: `done`
- **Dependencies**: T1.1

### T10.4: Request ID middleware
- **Status**: `done`
- **Dependencies**: T1.1

### T10.5: Router wiring and error handler
- **Status**: `done`
- **Dependencies**: T10.1-T10.4, T3.4, T3.5

### T10.6: Main entry point
- **Status**: `done`
- **Dependencies**: T10.5, T1.2

### T10.7: Unit tests for domain entities
- **Status**: `done`
- **Dependencies**: T2.1
- **Archivos**: `internal/domain/entity/basic_data_test.go`, `internal/domain/entity/entities_test.go`

### T10.8: Handler tests with mocked ports
- **Status**: `done`
- **Dependencies**: T4.x, T5.x, T6.x
- **Archivos**: `internal/infrastructure/adapter/http/handler/handler_test.go`

### T10.9: DTO mappers
- **Status**: `done`
- **Dependencies**: T4.x, T5.x
- **Archivos**: `internal/infrastructure/mapper/dto_mapper.go`
