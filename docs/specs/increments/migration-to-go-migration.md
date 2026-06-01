# Migration Contract - migration-to-go

**Status**: `awaiting-human-plan-approval`
**Fecha**: 2026-05-31
**Incremento**: `migration-to-go`

---

## 1. Objetivo

Definir el esquema de base de datos SQLite y el script de migración de datos desde la BD Spring Boot actual hacia la nueva estructura Go.

---

## 2. Migración de Esquema

### 2.1 Migración 001: Initial Schema

Archivo: `internal/infrastructure/adapter/repository/migrations/001_initial_schema.up.sql`

Crea todas las tablas de las 13 entidades existentes + tablas de soporte:

- `basic_data` (con campos `deleted_at`, `created_at`, `updated_at`)
- `home` (con FK a `image_url` y tabla join `home_label`)
- `label` (con campo `order`)
- `image_url`
- `video_url`
- `blog` (con FKs a `image_url`, `video_url`, `blog_type`)
- `blog_type` (con campo `order`)
- `skill_type` (con campo `order` y tabla join `skill_type_skill`)
- `skill` (con campo `order` y tabla join `skill_skill_son`)
- `skill_son` (con campo `order`)
- `experience` (con tabla join `experience_skill_son`)
- `education` (con campo `order`)
- `futured_project` (con campo `order` y FKs)
- `user_credentials`
- `pdf_cache`

### 2.2 Migración 002: New Entities

Archivo up: `internal/infrastructure/adapter/repository/migrations/002_new_entities.up.sql`
Archivo down: `internal/infrastructure/adapter/repository/migrations/002_new_entities.down.sql`

Crea las 5 nuevas entidades:

- `course`
- `certification`
- `language`
- `reference`
- `custom_section`

**Nota**: Cada migración debe tener su correspondiente `.down.sql` para permitir rollback.

---

## 3. Script de Migración de Datos

### 3.1 Ejecución

```bash
go run ./cmd/migrate --source /path/to/old.db --target /path/to/new.db
```

### 3.2 Pasos del Script

1. **Abrir BD fuente** (Spring Boot SQLite) y **BD destino** (nueva Go SQLite)
2. **Aplicar migraciones de esquema** en la BD destino (001 + 002)
3. **Migrar entidades básicas** (preservando IDs):
   - `basic_data`: copiar fila, mapear columnas camelCase → snake_case
   - `home`: copiar fila, mapear columnas
   - `label`: copiar filas, asignar `order = id` (orden existente por ID)
   - `image_url`: copiar filas
   - `video_url`: copiar filas
   - `blog`: copiar filas, mapear FKs
   - `blog_type`: copiar filas, asignar `order = id`
   - `skill_type`: copiar filas, asignar `order = id`
   - `skill`: copiar filas, asignar `order = id`
   - `skill_son`: copiar filas, asignar `order = id`
   - `experience`: copiar filas, asignar `order` basado en `year_start DESC` (más reciente = 1)
   - `education`: copiar filas, asignar `order = id`
   - `futured_project`: copiar filas, asignar `order = id`
4. **Migrar relaciones many-to-many**:
   - `home_label`: copiar desde tabla join existente
   - `skill_type_skill`: copiar desde tabla join existente
   - `skill_skill_son`: copiar desde tabla join existente
   - `experience_skill_son`: copiar desde tabla join existente
5. **Migrar credenciales de usuario**:
   - Leer username y password_hash de la tabla existente
   - Insertar en `user_credentials`
6. **Crear tablas nuevas vacías**:
   - `course`, `certification`, `language`, `reference`, `custom_section` (ya creadas por migración 002)
7. **Crear `pdf_cache` vacía** (ya creada por migración 001)
8. **Validar conteos**: comparar COUNT(*) entre fuente y destino para cada tabla
9. **Reportar resultado**: filas migradas, errores, warnings

### 3.3 Mapeo de Columnas (camelCase → snake_case)

| Spring Boot Column | Go Column |
|---|---|
| firstName | first_name |
| othersName | others_name |
| firstSurName | first_surname |
| othersSurName | others_surname |
| dateBirth | date_birth |
| startWorkingDate | start_working_date |
| greeting | greeting |
| greetingEng | greeting_eng |
| descriptionPdf | description_pdf |
| descriptionPdfEng | description_pdf_eng |
| imageUrlId | image_url_id |
| buttonWorkLabel | button_work_label |
| buttonWorkLabelEng | button_work_label_eng |
| buttonContactLabel | button_contact_label |
| buttonContactLabelEng | button_contact_label_eng |
| labelIds | (tabla join home_label) |
| cleanUrlTitle | clean_url_title |
| descriptionShort | description_short |
| descriptionShortEng | description_short_eng |
| imageUrlId | image_url_id |
| videoUrlId | video_url_id |
| blogTypeId | blog_type_id |
| skillIds | (tabla join skill_type_skill) |
| skillSonIds | (tabla join skill_skill_son) |
| yearStart | year_start |
| yearEnd | year_end |
| summaryPdf | summary_pdf |
| summaryPdfEng | summary_pdf_eng |
| descriptionItemsPdf | description_items_pdf |
| descriptionItemsPdfEng | description_items_pdf_eng |
| skillSonIds | (tabla join experience_skill_son) |
| startDate | start_date |
| endDate | end_date |
| experienceId | experience_id |
| imageListUrlId | image_list_url_id |
| imageUrlId | image_url_id |

### 3.4 Asignación de `order`

**Importante**: El campo `order` es interno de la base de datos. NO se expone en la API (ni en requests ni en responses). La lógica de ordenamiento se aplica en las consultas SQL (`ORDER BY order ASC` o según el criterio definido). Consultar la política de ordenamiento en la Master Spec para más detalles.

| Entidad | Fórmula de migración |
|---|---|
| Label | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |
| BlogType | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |
| SkillType | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |
| Skill | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |
| SkillSon | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |
| Experience | `order = ROW_NUMBER() OVER (ORDER BY year_start DESC)` |
| Education | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |
| FuturedProject | `order = ROW_NUMBER() OVER (ORDER BY id ASC)` |

**Reglas de auto-asignación post-migración**:
- En CREATE: `order = COALESCE((SELECT MAX(order) FROM table WHERE deleted_at IS NULL), 0) + 1`
- Soft-deleted entities se excluyen del cálculo
- El API NO expone el campo `order` al frontend; las respuestas se ordenan internamente antes de retornarse

### 3.5 Validación Post-Migración

El script debe verificar:
- COUNT(*) de cada tabla fuente == COUNT(*) de cada tabla destino
- Todos los IDs preservados
- Todas las FKs válidas en destino
- `order` asignado correctamente
- `user_credentials` tiene 1 fila
- `pdf_cache` tiene 0 filas

---

## 4. Rollback

Para rollback, el usuario debe:
1. Mantener la BD Spring Boot original como backup
2. Si la migración falla, el servicio Go no inicia y se puede revertir al servicio Spring Boot

---

## 5. Tablas de Soporte

### 5.1 user_credentials

| Columna | Tipo | Descripción |
|---|---|---|
| id | INTEGER PK | Identificador |
| username | TEXT UNIQUE | Username del admin |
| password_hash | TEXT | Bcrypt hash del password |
| created_at | TEXT | Timestamp UTC |
| updated_at | TEXT | Timestamp UTC |

### 5.2 pdf_cache

| Columna | Tipo | Descripción |
|---|---|---|
| id | INTEGER PK | Identificador |
| language | TEXT | english/spanish |
| template | TEXT | Template ID |
| data_hash | TEXT | SHA-256 hash |
| s3_key | TEXT | Clave S3 |
| created_at | TEXT | Timestamp UTC |
| file_size | INTEGER | Tamaño en bytes |

Índice único: `(language, template)`
