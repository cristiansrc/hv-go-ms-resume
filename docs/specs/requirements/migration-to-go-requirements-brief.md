# Requirements Brief - Migración de ms-resume a Go

**Status**: `ready-for-planner`
**Fecha**: 2026-05-31
**Owner**: cristiansrc
**Proyecto**: `hv-go-ms-resume`
**Incremento**: `migration-to-go`

---

## 1. Objetivo

Migrar el microservicio `ms-resume` de Spring Boot (Java) a Go 1.22+ manteniendo:
- 100% de compatibilidad con el contrato API actual definido en `openapi.yml` (endpoints, payloads, formato de error para respuestas 200).
- Todos los flujos funcionales existentes (CRUD de datos de CV, autenticación con Altcha, generación de PDFs, formulario de contacto, gestión de multimedia).
- Persistencia en SQLite (archivo local, sin servicio de base de datos externo).
- Integraciones con servicios downstream sin cambios funcionales.

Adicionalmente, se incorporan **nuevas entidades funcionales** que no existían en la versión Spring Boot: Cursos, Certificaciones, Idiomas, Referencias y Secciones Personalizadas.

**Nota sobre respuestas de error**: Los formatos de error se estandarizarán según los skills configurados del proyecto, no se migra el formato de error actual de Spring Boot. Solo se mantiene compatibilidad en las respuestas 200.

---

## 2. Contexto

El ecosistema del portfolio/CV (`cristiansrc.com`) está en proceso de reestructuración arquitectónica. El servicio `ms-resume` original fue construido con Spring Boot 3.5 + Java 21 + SQLite y presentaba los siguientes problemas:

- Overhead de memoria innecesario (~512MB mínimo) para un proyecto personal.
- Startup time lento (segundos vs milisegundos).
- Acoplamiento de bounded contexts en los repositorios originales.
- Convención de nombres inconsistente.

La decisión arquitectónica (ADR-001) establece migrar a Go 1.22+ con Chi router, mantener arquitectura hexagonal, y separar el servicio de renderizado de PDFs (`hv-py-ms-render-cv`) en un repositorio independiente.

**Relación con el sistema actual**:
- Los frontends (`hv-rt-fr-portal` y `hv-rt-fr-admin`) están activos y consumiendo la API actual.
- La migración no debe romper ningún frontend existente en las respuestas 200.
- El servicio de renderizado de PDFs (`hv-py-ms-render-cv`) se comunica vía HTTP interno.

---

## 3. Actores y Permisos

| Actor | Rol | Permisos | Acceso |
|---|---|---|---|
| **Visitante público** | Usuario anónimo | Leer datos públicos del CV, descargar PDF del CV, enviar mensaje de contacto | Endpoints `/public/*` sin autenticación |
| **Administrador (cristiansrc)** | Owner único | CRUD completo de todas las entidades, gestión de multimedia, login con Altcha | Endpoints CRUD + `/login` con JWT Bearer Token |
| **Servicio Render CV** | Servicio interno | Recibir solicitudes de generación de PDF | HTTP interno (red Docker), sin auth |
| **Telegram Bot** | Servicio externo | Recibir notificaciones de mensajes de contacto | HTTP outbound desde ms-resume |
| **AWS S3** | Servicio externo | Almacenar y servir imágenes/videos | AWS SDK v2 con Access Key + Secret Key |

**Nota**: No hay multi-usuario ni roles diferenciados. El administrador es un único usuario (single-owner).

---

## 4. Alcance (Scope)

### 4.1 Entidades Existentes (Migrar desde OpenAPI actual)

El OpenAPI actual (`openapi.yml`) define las siguientes entidades con sus contratos de respuesta 200:

| # | Entidad | Endpoint Base | Descripción | Campos Clave (Response 200) |
|---|---|---|---|---|
| E1 | **BasicData** | `/basic-data/{id}` | Datos personales del titular del CV | id, firstName, othersName, firstSurName, othersSurName, dateBirth, located, locatedEng, startWorkingDate, greeting, greetingEng, email, instagram, linkedin, x, github, description, descriptionEng, descriptionPdf[], descriptionPdfEng[], wrapper[], wrapperEng[] |
| E2 | **Home** | `/home/{id}` | Configuración de la página de inicio | id, greeting, greetingEng, imageUrl (ref), buttonWorkLabel, buttonWorkLabelEng, buttonContactLabel, buttonContactLabelEng, labels[] (ref) |
| E3 | **Label** | `/label` | Etiquetas que se visualizan en el home | id, name, nameEng, **order** |
| E4 | **ImageUrl** | `/image-url` | URLs de imágenes utilizadas en el portal (referencias a S3) | id, name, nameEng, url |
| E5 | **VideoUrl** | `/video-url` | URLs de videos almacenados en YouTube | id, name, nameEng, url |
| E6 | **Blog** | `/blog` | Artículos del blog tecnológico | id, title, titleEng, cleanUrlTitle, descriptionShort, description, descriptionShortEng, descriptionEng, imageUrl (ref), videoUrl (ref), blogType (ref) |
| E7 | **BlogType** | `/blog-type` | Tipos/categorías de blog | id, name, nameEng, **order** |
| E8 | **SkillType** | `/skill-type` | Categorías de habilidades técnicas | id, name, nameEng, **order**, skills[] (ref) |
| E9 | **Skill** | `/skill` | Habilidades técnicas específicas | id, name, nameEng, **order**, skillSons[] (ref) |
| E10 | **SkillSon** | `/skill-son` | Especializaciones/tecnologías específicas | id, name, nameEng, **order** |
| E11 | **Experience** | `/experience` | Experiencias laborales profesionales | id, yearStart, yearEnd, company, location, locationEng, position, positionEng, summary, summaryEng, summaryPdf, summaryPdfEng, descriptionItemsPdf[], descriptionItemsPdfEng[], skillSons[] (ref) | **Sin campo `order`**. Se ordena siempre por fecha más reciente primero (`yearStart DESC`). |
| E12 | **Education** | `/education` | Formación académica | id, institution, area, areaEng, degree, degreeEng, startDate, endDate, location, locationEng, highlights[], highlightsEng[], **order** |
| E13 | **FuturedProject** | `/futured-project` | Proyectos destacados del portfolio | id, name, nameEng, descriptionShort, descriptionShortEng, description, descriptionEng, experience (ref), imageListUrl (ref), imageUrl (ref), **order** |

### 4.2 Endpoints Públicos (Migrar)

| # | Endpoint | Método | Descripción | Response 200 |
|---|---|---|---|---|
| P1 | `/login` | POST | Autenticar usuario con credenciales + Altcha | `{ token: string }` (JWT válido 24h) |
| P2 | `/public/info-page` | GET | Información consolidada del portal | InfoPageResponse (home, basicData, skills[], experiences[], educations[], altchaChallenge) |
| P3 | `/public/curriculum/{language}` | GET | Descargar CV en PDF (english/spanish) | PDF binary con headers Content-Disposition, Content-Type, Content-Length |
| P4 | `/public/contact` | POST | Enviar mensaje de contacto con Altcha | 200 sin body |
| P5 | `/public/challenge` | GET | Obtener challenge Altcha | AltchaChallengeResponse (algorithm, challenge, salt, signature) |
| P6 | `/public/blog` | GET | Blog paginado | BlogPageResponse (content[], pageable, last, totalPages, totalElements, etc.) |
| P7 | `/public/blog/{id}` | GET | Artículo individual del blog | BlogResponse |
| P8 | `/public/blog-type` | GET | Todos los tipos de blog | BlogTypeResponse[] |
| P9 | `/public/blog-type/{id}` | GET | Tipo de blog por ID | BlogTypeResponse |

### 4.3 Nuevas Entidades (Agregar)

| # | Entidad | Descripción | Campos Clave |
|---|---|---|---|
| N1 | **Course** | Cursos completados | id, name, nameEng, institution, institutionEng, completionDate, description (HTML para portal), descriptionEng (HTML), summaryPdf (texto plano para PDF), summaryPdfEng, certificateUrl (opcional), **order** |
| N2 | **Certification** | Certificaciones profesionales | id, name, nameEng, issuingOrganization, issuingOrganizationEng, issueDate, expirationDate (opcional), verificationUrl, credentialId, description (HTML), descriptionEng (HTML), summaryPdf (texto plano), summaryPdfEng, **order** |
| N3 | **Language** | Idiomas conocidos | id, language, languageEng, readingLevel, writingLevel, speakingLevel, **order** |
| N4 | **Reference** | Referencias profesionales | id, fullName, position, company, companyEng, email, phone, relationship, relationshipEng, **order** |
| N5 | **CustomSection** | Secciones de contenido libre | id, title, titleEng, content (HTML para portal), contentEng (HTML), summaryPdf (texto plano para PDF), summaryPdfEng, **order**, visible |

### 4.4 Nuevos Endpoints Públicos

| # | Endpoint | Método | Descripción | Response 200 |
|---|---|---|---|---|
| NP1 | `/public/templates` | GET | Listar templates/temas disponibles para PDF | `{ templates: [{ id: "engineeringclassic", name: "Engineering Classic" }, { id: "engineeringresumes", name: "Engineering Resumes" }, { id: "moderncv", name: "Modern CV" }, { id: "sb2nov", name: "SB2Nov" }] }` |
| NP2 | `/public/languages` | GET | Listar idiomas disponibles para PDF | `{ languages: [{ code: "english", name: "English" }, { code: "spanish", name: "Español" }] }` |

---

## 5. No Objetivos (Out of Scope)

- **Multi-usuario o gestión de roles**: El sistema es single-owner. No se implementarán roles, permisos granulares ni gestión de usuarios múltiples.
- **Migración de base de datos**: Se mantiene SQLite como archivo local. No se migrará a PostgreSQL, MySQL u otro motor.
- **Funcionalidades de analytics o métricas**: No se implementará tracking de visitas, popularidad de contenido, ni dashboards de uso.
- **Webhooks o eventos asíncronos**: La comunicación con servicios downstream es síncrona (HTTP).
- **Versionado de CVs**: No se implementará historial de versiones del CV.
- **Notificaciones push o email**: Solo se usa Telegram para notificaciones de contacto.
- **Búsqueda full-text**: No se requiere búsqueda avanzada dentro del contenido del CV.
- **Internacionalización de la UI del admin**: El admin panel es un frontend separado; este servicio solo maneja datos multi-idioma para generación de PDFs.
- **Migración del formato de error**: Los formatos de error se estandarizarán según los skills configurados, no se replica el formato de Spring Boot.
- **Cualquier funcionalidad nueva no listada en las secciones 4.3 y 4.4**: Queda para iteraciones futuras.

---

## 6. Flujos de Usuario

### 6.1 Flujo: Login del Administrador

1. El administrador envía credenciales (user, password) + respuesta Altcha a `POST /v1/ms-resume/login`.
2. El servicio valida el challenge Altcha.
3. Si es válido, valida las credenciales contra el almacenamiento interno.
4. Si son válidas, retorna un JWT Bearer Token (validez 24 horas).
5. Si algo falla, retorna error 401.
6. El administrador usa el token en el header `Authorization: Bearer <token>` para todas las operaciones CRUD.

**Caso alterno**: Token expirado → el frontend debe solicitar nuevo login (no hay refresh token en esta iteración).

### 6.2 Flujo: Gestión de Entidades del CV (CRUD genérico)

1. El administrador autenticado envía una solicitud CRUD (GET, POST, PUT, DELETE) al endpoint correspondiente.
2. El servicio valida el token JWT.
3. El servicio valida los datos de entrada.
4. El servicio persiste los cambios en SQLite (soft delete para DELETE).
5. El servicio retorna la entidad creada/actualizada o confirmación de eliminación (204 No Content).

**Caso alterno**: Validación fallida → retorna error según skill de error response.
**Caso alterno**: Token inválido/expirado → retorna error 401.
**Caso alterno**: Entidad no encontrada → retorna error 404.

### 6.3 Flujo: Descarga de CV en PDF (con caché inteligente)

1. El usuario (público o administrador) solicita `GET /v1/ms-resume/public/curriculum/{language}` con parámetro opcional de template.
2. El servicio calcula el **hash de datos actual** de todas las entidades que componen el CV.
3. El servicio consulta en SQLite si existe un PDF cacheado para la combinación `{language, template}` con el **mismo hash**.
4. **Si existe y el hash coincide** (datos no cambiaron):
   - El servicio descarga el PDF desde AWS S3.
   - Retorna el PDF como respuesta binaria al usuario.
5. **Si no existe o el hash difiere** (datos cambiaron o es primera vez):
   - El servicio obtiene los datos del CV desde SQLite.
   - Transforma los datos al schema de RenderCV (`CvData`):
     - `name`: BasicData.firstName + BasicData.firstSurName
     - `email`: BasicData.email
     - `phone`: BasicData.phone (si existe)
     - `location`: BasicData.located o BasicData.locatedEng según idioma
     - `headline`: BasicData.wrapper o wrapperEng según idioma
     - `social_networks`: [{ network: "LinkedIn", username: BasicData.linkedin }, { network: "GitHub", username: BasicData.github }, ...]
     - `sections`: { experience: [...], education: [...], skills: [...], custom_sections: [...] }
     - `locale`: { language: language } (english/spanish)
     - `design`: { theme: template } (engineeringclassic/engineeringresumes/moderncv/sb2nov)
   - Envía los datos a `hv-py-ms-render-cv` vía `POST /render` (HTTP interno, JSON).
   - El servicio de renderizado genera el PDF y retorna `{ pdf_base64: string }`.
   - **Si existía un PDF anterior con hash diferente**: lo elimina de S3.
   - Guarda el nuevo PDF en S3 con la clave `{language}/{template}/{hash}.pdf`.
   - Registra en SQLite: `language`, `template`, `hash`, `s3_key`, `created_at`.
   - Retorna el PDF como respuesta binaria al usuario.

**Caso alterno**: Servicio de renderizado no disponible → retorna error 503.
**Caso alterno**: Idioma no soportado → retorna error 400.
**Caso alterno**: Template no encontrado → retorna error 400.
**Caso alterno**: S3 no disponible para leer cache → regenera el PDF (fallback).
**Caso alterno**: S3 no disponible para guardar cache → retorna el PDF generado sin cachear (no es error, solo no se cachea).

### 6.4 Flujo: Vista Previa de PDF

1. El administrador solicita una vista previa del PDF con un template/idioma específico.
2. El servicio sigue el mismo flujo que 6.3 (con caché inteligente).
3. El PDF se retorna como respuesta binaria para visualización en el frontend.

**Nota**: La vista previa también se cachea. Si el admin cambia datos y solicita vista previa de nuevo, se regenera.

### 6.5 Flujo: Formulario de Contacto

1. El visitante público completa el formulario de contacto en el portal.
2. El frontend obtiene un challenge Altcha desde `GET /v1/ms-resume/public/challenge`.
3. El visitante resuelve el challenge y envía el formulario (POST `/public/contact`) con name, email, message y la solución Altcha.
4. El servicio valida el challenge Altcha.
5. Si es válido, el servicio envía una notificación al Telegram Bot con los datos del mensaje.
6. El servicio confirma el envío exitoso (200 sin body).

**Caso alterno**: Challenge inválido → retorna error 400 (spam detectado).
**Caso alterno**: Telegram no disponible → el mensaje se registra internamente pero no se notifica (no se pierde el dato, no se retorna error al usuario).

### 6.6 Flujo: Gestión de Multimedia (Image/Video URLs)

1. El administrador autenticado envía un POST a `/image-url` o `/video-url` con los datos (name, nameEng, file/url).
2. El servicio valida los datos de entrada.
3. Para imágenes: el archivo ya debe estar en S3; el servicio registra la URL en SQLite.
4. Para videos: el servicio registra la URL de YouTube en SQLite.
5. El servicio retorna el ID del registro creado.
6. Para eliminar, el administrador envía DELETE con el identificador del registro.
7. El servicio marca el registro como eliminado (soft delete).

**Nota**: El servicio no sube archivos directamente a S3; registra URLs de archivos ya subidos. La subida a S3 se maneja externamente.

---

## 7. Entidades Funcionales

### 7.1 Entidades Existentes (del OpenAPI)

| Entidad | Descripción | Atributos Clave | Sensibilidad | Patrón Multi-idioma |
|---|---|---|---|---|
| **BasicData** | Datos personales del titular | firstName, othersName, firstSurName, othersSurName, dateBirth, located/locatedEng, startWorkingDate, greeting/greetingEng, email, instagram, linkedin, x, github, description/descriptionEng, descriptionPdf[]/descriptionPdfEng[], wrapper[]/wrapperEng[] | Email es dato personal | Campos duplicados: `*` y `*Eng` |
| **Home** | Configuración del home | greeting/greetingEng, imageUrl (ref), buttonWorkLabel/buttonWorkLabelEng, buttonContactLabel/buttonContactLabelEng, labels[] (ref) | Ninguna | Campos duplicados ES/EN |
| **Label** | Etiquetas del home | name, nameEng, **order** | Ninguna | Campos duplicados ES/EN |
| **ImageUrl** | Referencias a imágenes S3 | name, nameEng, url | Ninguna | Campos duplicados ES/EN |
| **VideoUrl** | Referencias a videos YouTube | name, nameEng, url | Ninguna | Campos duplicados ES/EN |
| **Blog** | Artículos del blog | title/titleEng, cleanUrlTitle, descriptionShort/descriptionShortEng, description/descriptionEng, imageUrl (ref), videoUrl (ref), blogType (ref) | Ninguna | Campos duplicados ES/EN |
| **BlogType** | Tipos de blog | name, nameEng, **order** | Ninguna | Campos duplicados ES/EN |
| **SkillType** | Categorías de habilidades | name, nameEng, **order**, skills[] (ref) | Ninguna | Campos duplicados ES/EN |
| **Skill** | Habilidades técnicas | name, nameEng, **order**, skillSons[] (ref) | Ninguna | Campos duplicados ES/EN |
| **SkillSon** | Especializaciones/tecnologías | name, nameEng, **order** | Ninguna | Campos duplicados ES/EN |
| **Experience** | Experiencia laboral | yearStart, yearEnd, company, location/locationEng, position/positionEng, summary/summaryEng, summaryPdf/summaryPdfEng, descriptionItemsPdf[]/descriptionItemsPdfEng[], skillSons[] (ref) | Ninguna | Campos duplicados ES/EN. **Sin campo `order`**: se ordena por `yearStart DESC`. |
| **Education** | Formación académica | institution, area/areaEng, degree/degreeEng, startDate, endDate, location/locationEng, highlights[]/highlightsEng[], **order** | Ninguna | Campos duplicados ES/EN |
| **FuturedProject** | Proyectos destacados | name/nameEng, descriptionShort/descriptionShortEng, description/descriptionEng, experience (ref), imageListUrl (ref), imageUrl (ref), **order** | Ninguna | Campos duplicados ES/EN |

### 7.2 Nuevas Entidades

| Entidad | Descripción | Atributos Clave | Sensibilidad | Patrón Multi-idioma |
|---|---|---|---|---|
| **Course** | Cursos completados | name/nameEng, institution/institutionEng, completionDate, description (HTML portal)/descriptionEng (HTML), summaryPdf (texto PDF)/summaryPdfEng, certificateUrl (opcional), **order** | Ninguna | Campos duplicados ES/EN + separación HTML/Texto |
| **Certification** | Certificaciones profesionales | name/nameEng, issuingOrganization/issuingOrganizationEng, issueDate, expirationDate (opcional), verificationUrl, credentialId, description (HTML)/descriptionEng (HTML), summaryPdf (texto)/summaryPdfEng, **order** | Ninguna | Campos duplicados ES/EN + separación HTML/Texto |
| **Language** | Idiomas conocidos | language/languageEng, readingLevel, writingLevel, speakingLevel, **order** | Ninguna | Campos duplicados ES/EN |
| **Reference** | Referencias profesionales | fullName, position, company/companyEng, email, phone, relationship/relationshipEng, **order** | Email y teléfono son datos de terceros | Campos duplicados ES/EN |
| **CustomSection** | Secciones de contenido libre | title/titleEng, content (HTML portal)/contentEng (HTML), summaryPdf (texto PDF)/summaryPdfEng, **order**, visible | Ninguna | Campos duplicados ES/EN + separación HTML/Texto |

### 7.3 Entidades de Soporte

| Entidad | Descripción | Sensibilidad |
|---|---|---|
| **ContactMessage** | Mensajes de contacto (efímero, no se persisten) | Email y mensaje son datos personales del remitente |
| **UserCredentials** | Credenciales de admin | **Crítico**: password debe estar hasheado (bcrypt) |
| **AltchaChallenge** | Challenge para anti-spam | Efímero, se genera por request |
| **PdfCache** | Metadata de PDFs cacheados en S3 | Ninguna |

#### PdfCache (detalle)

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | INTEGER PK | Identificador único |
| `language` | VARCHAR | Idioma del PDF (`english`, `spanish`) |
| `template` | VARCHAR | Template/tema usado para generar el PDF |
| `data_hash` | VARCHAR | Hash SHA-256 de los datos del CV usados para generar este PDF |
| `s3_key` | VARCHAR | Clave del archivo en S3 (ej: `pdfs/english/default/abc123.pdf`) |
| `created_at` | TIMESTAMP | Fecha de creación del cache |
| `file_size` | INTEGER | Tamaño del PDF en bytes |

**Índice único**: `(language, template)` → solo un PDF cacheado por combinación idioma+template.
**Invalidación**: Cuando el `data_hash` calculado difiere del almacenado, se elimina el PDF antiguo de S3 y se genera uno nuevo.

### 7.4 Política de Ordenamiento

Todas las entidades que se retornan como listas en las respuestas siguen esta política:

| Entidad | Criterio de Ordenamiento | Campo `order` |
|---|---|---|
| **Experience** | `yearStart DESC` (más reciente primero) | **No aplica** |
| **Label** | `order ASC` | Sí |
| **BlogType** | `order ASC` | Sí |
| **SkillType** | `order ASC` | Sí |
| **Skill** (dentro de SkillType) | `order ASC` | Sí |
| **SkillSon** (dentro de Skill) | `order ASC` | Sí |
| **Education** | `order ASC` | Sí |
| **FuturedProject** | `order ASC` | Sí |
| **Course** | `order ASC` | Sí |
| **Certification** | `order ASC` | Sí |
| **Language** | `order ASC` | Sí |
| **Reference** | `order ASC` | Sí |
| **CustomSection** | `order ASC` | Sí |
| **Blog** | `id DESC` (más reciente primero, ya usa paginación con sort) | No aplica |
| **ImageUrl** | `id ASC` | No aplica |
| **VideoUrl** | `id ASC` | No aplica |

**Reglas**:
- El campo `order` es un `INTEGER` en SQLite.
- Al crear una nueva entidad, el `order` se asigna automáticamente como `MAX(order) + 1` de las entidades no eliminadas.
- El administrador puede reordenar entidades mediante UPDATE del campo `order`.
- Las entidades con soft delete (`deleted_at IS NOT NULL`) se excluyen del cálculo de `order` y no aparecen en las listas.
- Las respuestas de listas retornan los elementos ordenados según la columna correspondiente. El frontend no reordena.

### 7.5 Estrategia de Caché de PDFs

Para optimizar la generación de PDFs y reducir la carga del servicio de renderizado, se implementa un sistema de caché inteligente con invalidación basada en cambios de datos.

#### Principio

```
Solicitud de PDF
       │
       ▼
┌─────────────────────────────┐
│ Calcular hash de datos CV   │
│ (SHA-256 de todas las       │
│  entidades que afectan PDF) │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│ ¿Existe cache con mismo     │
│ hash para {language,template}?│
└──────┬──────────────┬───────┘
       │ SÍ           │ NO
       ▼              ▼
┌──────────────┐  ┌──────────────────────────┐
│ Descargar    │  │ Generar PDF vía RenderCV │
│ desde S3     │  │ Guardar en S3            │
│ (cache hit)  │  │ Actualizar cache en SQLite│
└──────┬───────┘  └──────────┬───────────────┘
       │                     │
       ▼                     ▼
┌──────────────────────────────────────────┐
│ Retornar PDF binario al usuario          │
└──────────────────────────────────────────┘
```

#### Entidades que afectan el hash del PDF

Todas las entidades que contienen datos que se incluyen en el CV generado:

| Entidad | ¿Afecta hash? | Razón |
|---|---|---|
| **BasicData** | ✅ Sí | Datos personales, descripción, wrapper |
| **Experience** | ✅ Sí | Experiencia laboral en el CV |
| **Education** | ✅ Sí | Formación académica en el CV |
| **Skill / SkillSon** | ✅ Sí | Habilidades en el CV |
| **Course** | ✅ Sí | Cursos en el CV |
| **Certification** | ✅ Sí | Certificaciones en el CV |
| **Language** | ✅ Sí | Idiomas en el CV |
| **Reference** | ✅ Sí | Referencias en el CV |
| **CustomSection** | ✅ Sí | Secciones personalizadas en el CV |
| **Home** | ❌ No | Solo afecta el portal, no el PDF |
| **Label** | ❌ No | Solo afecta el portal |
| **Blog / BlogType** | ❌ No | Solo afecta el portal |
| **ImageUrl / VideoUrl** | ❌ No | Solo referencias para el portal |
| **FuturedProject** | ❌ No | Solo afecta el portal |

#### Cálculo del Hash

- Se serializan los datos relevantes de todas las entidades que afectan el PDF en un formato determinista (JSON ordenado).
- Se calcula SHA-256 sobre la serialización.
- El hash es único para cada combinación de datos + idioma + template.

#### Invalidación del Caché

| Evento | Acción |
|---|---|
| CREATE de entidad que afecta PDF | El próximo request de PDF tendrá hash diferente → regeneración automática |
| UPDATE de entidad que afecta PDF | El próximo request de PDF tendrá hash diferente → regeneración automática |
| DELETE (soft) de entidad que afecta PDF | El próximo request de PDF tendrá hash diferente → regeneración automática |
| Cambio de template | Nuevo hash → regeneración automática |
| Cambio de idioma | Nuevo hash → regeneración automática |

**Nota**: No se necesita invalidación explícita. El hash se calcula en cada request y la comparación con el cache determina si se regenera.

#### Clave de S3 para PDFs cacheados

Formato: `pdfs/{language}/{template}/{data_hash}.pdf`

Ejemplo: `pdfs/english/engineeringclassic/a1b2c3d4e5f6...pdf`

#### Fallbacks

| Escenario | Comportamiento |
|---|---|
| S3 no disponible para **leer** cache | Regenerar el PDF (no es error, solo no se usa cache) |
| S3 no disponible para **guardar** cache | Retornar el PDF generado sin cachear (no es error) |
| RenderCV no disponible | Retornar 503 (error real, no hay fallback) |
| Cache existe en SQLite pero PDF no existe en S3 | Regenerar el PDF (posible inconsistencia, se corrige automáticamente) |

---

## 8. Integraciones

| # | Integración | Dirección | Protocolo | Propósito | Criticidad |
|---|---|---|---|---|---|
| I1 | **hv-rt-fr-portal → hv-go-ms-resume** | Inbound | HTTPS REST | Consumo de endpoints públicos (info-page, curriculum, contact, challenge, blog, blog-type) | Alta - el portal público depende de esto |
| I2 | **hv-rt-fr-admin → hv-go-ms-resume** | Inbound | HTTPS REST + JWT | CRUD completo y login | Alta - el admin panel depende de esto |
| I3 | **hv-go-ms-resume → hv-py-ms-render-cv** | Outbound | HTTP interno (Docker) | Generación de PDFs del CV (`POST /render`) + Health check (`GET /health`) | Alta - sin esto no hay descarga de CV |
| I4 | **hv-go-ms-resume → AWS S3** | Outbound | HTTPS (AWS SDK v2) | Almacenamiento de imágenes/videos Y caché de PDFs generados | Alta - sin esto no hay multimedia ni caché de PDFs |
| I5 | **hv-go-ms-resume → Telegram Bot API** | Outbound | HTTPS | Notificaciones de mensajes de contacto | Baja - es notificación, no bloqueante |

### 8.1 Contrato con hv-py-ms-render-cv (detalles del OpenAPI)

**Endpoint**: `POST /render`

**Request** (`RenderRequest`):
```json
{
  "cv": {
    "name": "string (required, 1-200 chars)",
    "email": "string (required, email format, max 254)",
    "phone": "string (optional, max 30)",
    "location": "string (optional, max 200)",
    "headline": "string (optional, max 200)",
    "photo": "string (optional, URL/path, max 2048)",
    "website": "string (optional, URI format, max 2048)",
    "social_networks": [
      { "network": "string (required, 1-50)", "username": "string (required, 1-100)" }
    ],
    "sections": {
      "experience": [{ "company", "position", "start_date", "end_date", "highlights[]" }],
      "education": [{ "institution", "area", "degree", "start_date", "end_date" }],
      "skills": [{ "skill", "level", "highlights[]" }],
      "custom_section_name": [{ ... }]
    },
    "locale": { "language": "english|spanish|french|..." },
    "design": { "theme": "engineeringclassic|engineeringresumes|moderncv|sb2nov" }
  }
}
```

**Response** (`RenderResponse`):
```json
{
  "pdf_base64": "string (required, base64-encoded PDF)"
}
```

**Idiomas soportados** (16): english, spanish, french, german, italian, portuguese, dutch, danish, russian, turkish, hindi, indonesian, japanese, korean, mandarin_chinese.

**Templates soportados** (4): engineeringclassic, engineeringresumes, moderncv, sb2nov.

**Health check**: `GET /health` → `{ status: "healthy|unhealthy", timestamp, version }`

**Formato de error de render-cv**: `ApiErrorResponse` con `{ timestamp, status, error, code, message, path, trace_id, details[] }`. Este formato es diferente al de ms-resume. El adapter de ms-resume debe traducir estos errores al formato de error estandarizado del proyecto Go.

**Timeout recomendado**: 30 segundos (la generación de PDF puede tardar).
**Retry**: 1 reintento en caso de fallo transitorio.

---

## 9. Seguridad y Restricciones

| Restricción | Descripción |
|---|---|
| **Autenticación JWT** | Todos los endpoints CRUD requieren JWT Bearer Token válido. Los endpoints `/public/*` y `/login` no requieren autenticación. |
| **Altcha en login** | El endpoint `/login` requiere validación de challenge Altcha además de credenciales. |
| **Altcha en contacto** | El endpoint `/public/contact` requiere validación de challenge Altcha. |
| **Password hashing** | La contraseña del administrador debe almacenarse hasheada (bcrypt o equivalente). Nunca en texto plano. |
| **Compatibilidad de respuestas 200** | Los paths, métodos HTTP y formato de respuesta 200 deben ser idénticos al OpenAPI actual para no romper los frontends. |
| **Formato de error** | Las respuestas de error se estandarizarán según los skills configurados del proyecto (no se replica el formato de Spring Boot). |
| **Sin exposición de credenciales** | AWS Access Key, Secret Key, Telegram Bot Token y JWT Secret deben provenir de variables de entorno, nunca hardcodeados. |
| **Validación de entrada** | Todos los datos de entrada deben validarse antes de persistir. Campos requeridos, formatos de email, URLs, fechas, longitudes máximas según OpenAPI. |
| **SQLite file permissions** | El archivo SQLite debe tener permisos restrictivos en el sistema de archivos (solo lectura/escritura para el proceso del servicio). |
| **Límite de tamaño multimedia** | Límite de 2MB para imágenes. Se recomienda optimización: redimensionar imágenes muy grandes, eliminar metadatos EXIF innecesarios para reducir peso. Si la optimización no es viable en esta iteración, al menos rechazar archivos > 2MB con error descriptivo. |

---

## 10. Edge Cases

| # | Escenario | Comportamiento Esperado |
|---|---|---|
| E1 | **Token JWT expirado** | Retornar 401 con mensaje claro. El frontend debe redirigir al login. |
| E2 | **Token JWT malformado** | Retornar 401. No exponer detalles internos del error. |
| E3 | **Altcha inválido en login** | Retornar 401 indicando que la validación anti-spam falló. |
| E4 | **Altcha inválido en contacto** | Retornar 400 (spam detectado). |
| E5 | **Altcha expirado** | Retornar 400 indicando que el challenge expiró y debe obtener uno nuevo. |
| E6 | **Entidad referenciada no existe** | Si se intenta asociar una entidad que no existe (ej: skillSonId inexistente), retornar 400. |
| E7 | **Duplicados** | Si el dominio lo permite (ej: dos experiencias en la misma empresa), permitir. Si no (ej: mismo idioma duplicado), retornar 409 Conflict. |
| E8 | **Servicio de renderizado caído** | Retornar 503 Service Unavailable con mensaje descriptivo. No reintentar indefinidamente. |
| E9 | **AWS S3 caído** | Las operaciones de CRUD que no involucren S3 deben seguir funcionando. Si se necesita una URL de S3 y no está disponible, retornar 503. |
| E10 | **Telegram caído** | El mensaje de contacto se procesa pero la notificación falla. No retornar error al usuario del formulario. |
| E11 | **Imagen > 2MB** | Retornar 400 indicando el límite de tamaño. Si se implementa optimización, intentar redimensionar/eliminar metadatos primero. |
| E12 | **Tipo de archivo no permitido** | Retornar 400 con mensaje indicando tipos permitidos. |
| E13 | **Concurrencia en actualización** | Si dos requests actualizan la misma entidad simultáneamente, el último gana (no se requiere locking optimista en esta iteración). |
| E14 | **Base de datos SQLite corrupta** | El servicio debe detectar la corrupción y retornar 500 con log detallado. |
| E15 | **Idioma no soportado para PDF** | Retornar 400 con lista de idiomas disponibles. |
| E16 | **Template no encontrado** | Retornar 400 con lista de templates disponibles. |
| E17 | **Soft delete: entidad "eliminada" accedida por ID** | Retornar 404 como si no existiera. No retornar en listas. |
| E18 | **Referencia con email/teléfono inválido** | Validar formato de email y teléfono. Si son inválidos, retornar 400. |
| E19 | **HTML malformado en campos de portal** | El servicio acepta HTML tal cual (no sanitiza). La sanitización es responsabilidad del frontend al renderizar. |

---

## 11. Criterios de Aceptación

### CA1: Compatibilidad de Respuestas 200
- [ ] Todos los endpoints del OpenAPI actual están disponibles en la versión Go con el mismo path, método HTTP y formato de response 200.
- [ ] Los schemas de respuesta (BasicDataResponse, HomeResponse, ImageUrlResponse, VideoUrlResponse, LabelResponse, BlogResponse, BlogPageResponse, BlogTypeResponse, SkillTypeResponse, SkillResponse, SkillSonResponse, ExperienceResponse, EducationResponse, FuturedProjectResponse, InfoPageResponse, AltchaChallengeResponse) son idénticos al OpenAPI actual.
- [ ] Los frontends (`hv-rt-fr-portal` y `hv-rt-fr-admin`) funcionan sin cambios de código contra la nueva API.

### CA2: Autenticación
- [ ] El login con credenciales válidas + Altcha válido retorna un JWT válido (24h).
- [ ] El login con credenciales inválidas retorna 401.
- [ ] El login con Altcha inválido retorna 401.
- [ ] Los endpoints CRUD sin token retornan 401.
- [ ] Los endpoints CRUD con token expirado retornan 401.
- [ ] Los endpoints `/public/*` son accesibles sin token.

### CA3: CRUD de Entidades Existentes
- [ ] BasicData, Home, Label, ImageUrl, VideoUrl, Blog, BlogType, SkillType, Skill, SkillSon, Experience, Education, FuturedProject soportan operaciones CREATE, READ, UPDATE, DELETE.
- [ ] Las operaciones CREATE y UPDATE validan los campos requeridos según OpenAPI (minLength, maxLength, format, required).
- [ ] READ retorna los datos en el formato esperado por los frontends.
- [ ] DELETE aplica soft delete: la entidad ya no es accesible vía GET ni aparece en listas.

### CA4: Nuevas Entidades
- [ ] Course, Certification, Language, Reference, CustomSection soportan operaciones CREATE, READ, UPDATE, DELETE.
- [ ] Cada entidad tiene validación de campos requeridos.
- [ ] Las nuevas entidades respetan el patrón multi-idioma (campos ES/EN).
- [ ] Las nuevas entidades con contenido para portal usan campos HTML; las de PDF usan campos de texto plano.
- [ ] Las nuevas entidades son accesibles vía endpoints públicos de solo lectura.

### CA5: Multimedia
- [ ] Se puede registrar una URL de imagen y se retorna el ID.
- [ ] Se puede registrar una URL de video y se retorna el ID.
- [ ] Se puede eliminar un registro de imagen/video (soft delete).
- [ ] Imágenes > 2MB son rechazadas con error descriptivo.
- [ ] Si se implementa optimización de imágenes, se redimensionan y eliminan metadatos EXIF.

### CA6: Generación de PDF y Caché
- [ ] Se puede descargar el CV en PDF en un idioma seleccionado (english/spanish).
- [ ] Se puede seleccionar un template/tema para el PDF.
- [ ] Se puede obtener una vista previa del PDF.
- [ ] Si el servicio de renderizado no está disponible, se retorna 503.
- [ ] Idiomas y templates no válidos retornan 400.
- [ ] El primer request de un PDF (language+template) genera el PDF y lo guarda en S3.
- [ ] El segundo request con los mismos datos y mismo language+template sirve el PDF desde S3 (cache hit).
- [ ] Si se modifica cualquier entidad que afecta el PDF (BasicData, Experience, Education, Skill, Course, Certification, Language, Reference, CustomSection), el próximo request regenera el PDF.
- [ ] Si S3 no está disponible para leer cache, se regenera el PDF sin error.
- [ ] Si S3 no está disponible para guardar cache, se retorna el PDF generado sin cachear.
- [ ] El hash de datos se calcula correctamente y es determinista para los mismos datos.
- [ ] La clave de S3 sigue el formato `pdfs/{language}/{template}/{data_hash}.pdf`.

### CA7: Formulario de Contacto
- [ ] El challenge Altcha se genera correctamente (GET `/public/challenge`).
- [ ] Un mensaje con challenge válido se procesa y notifica a Telegram.
- [ ] Un mensaje con challenge inválido es rechazado con 400.
- [ ] Si Telegram no está disponible, el mensaje se procesa pero no se retorna error al usuario.

### CA8: Base de Datos
- [ ] El servicio utiliza SQLite como archivo local.
- [ ] No requiere un servicio de base de datos externo.
- [ ] Las migraciones de esquema se aplican al iniciar el servicio.
- [ ] Los datos persisten entre reinicios del servicio.
- [ ] Soft delete implementado: campo `deleted_at` o equivalente en todas las entidades.

### CA9: Seguridad
- [ ] La contraseña del administrador está hasheada en la base de datos.
- [ ] Las credenciales de AWS, Telegram y JWT Secret provienen de variables de entorno.
- [ ] El archivo SQLite tiene permisos restrictivos.

### CA10: Performance
- [ ] El servicio inicia en menos de 2 segundos.
- [ ] El consumo de memoria en reposo es menor a 50MB.
- [ ] Los endpoints públicos responden en menos de 500ms (p95).

### CA11: Nuevos Endpoints Públicos
- [ ] `GET /public/templates` retorna la lista de templates disponibles.
- [ ] `GET /public/languages` retorna la lista de idiomas disponibles.

---

## 12. Preguntas Abiertas

### Resueltas

| # | Pregunta | Respuesta |
|---|---|---|
| PA1 | ¿Cuáles son exactamente los campos de cada entidad? | **Resuelta**: Se extrajeron del OpenAPI actual (`openapi.yml`). Ver sección 7.1. |
| PA3 | ¿Las credenciales de admin están en SQLite? | **Resuelta**: Sí, se asume que están en SQLite con password hasheado. Planner debe verificar la tabla actual. |
| PA6 | ¿Límite de tamaño para multimedia? | **Resuelta**: 2MB. Se recomienda optimización (redimensionar, eliminar metadatos EXIF) si es viable. |
| PA7 | ¿Endpoint para listar templates/idiomas? | **Resuelta**: Sí, se agregan `GET /public/templates` y `GET /public/languages`. |
| PA8 | ¿HTML para portal, texto plano para PDF? | **Resuelta**: Sí. El OpenAPI actual ya usa este patrón (ej: `description` para portal HTML, `descriptionPdf` para texto plano). Las nuevas entidades seguirán el mismo patrón. |
| PA9 | ¿Soft delete? | **Resuelta**: Sí, todas las entidades usarán soft delete. |
| PA10 | ¿El orden de entidades importa? | **Resuelta**: Sí. Todas las entidades (excepto Experience) llevan campo `order` (INTEGER) en la base de datos. Las respuestas retornan las listas ordenadas por `order ASC`. Experience se ordena siempre por `yearStart DESC` (más reciente primero) sin campo `order`. |

### Críticas (bloquean handoff a Planner)

_No hay preguntas críticas abiertas. Todas resueltas._

### No Críticas (Planner puede decidir)

| # | Pregunta | Impacto |
|---|---|---|
| PA5 | ¿Se debe implementar rate limiting en esta iteración o queda como mejora futura? | Afecta complejidad pero no bloquea la migración. |
| PA11 | ¿La optimización de imágenes (redimensionar, eliminar metadatos) se implementa en esta iteración o queda como mejora? | Si se implementa, se necesita una librería de procesamiento de imágenes en Go. Si no, solo se rechazan archivos > 2MB. |
| PA12 | ¿Los campos HTML del portal deben ser sanitizados en el backend o la sanitización es responsabilidad exclusiva del frontend? | Afecta seguridad (XSS). Si el backend sanitiza, se necesita una librería de sanitización HTML. |

---

## 12.1 Datos Actuales en Producción (para migración)

Análisis de la respuesta actual de `GET /v1/ms-resume/public/info-page` en producción (`api.cristiansrc.com`):

| Entidad | Cantidad | Observaciones |
|---|---|---|
| **Home** | 1 | Con greeting ES/EN, botones, 1 imageUrl ref, 3 labels |
| **BasicData** | 1 | Con todos los campos poblados, description (HTML), descriptionPdf (array texto), wrapper (array) |
| **Label** | 3 | "Líder Técnico", "Desarrollador", "Arquitecto" |
| **ImageUrl** | 1+ | Al menos 1 (foto del home). Posiblemente más para blog/projects |
| **VideoUrl** | 0+ | No visible en info-page, pero existe el endpoint |
| **SkillType** | 6 | "Lenguajes y frameworks", "Base de datos", "Programación en la nube", "Habilidades Blandas", "Arquitectura", "Estrategia de IA" |
| **Skill** | ~30+ skills distribuidos en los 6 SkillTypes | Cada skill tiene 1+ skillSons |
| **SkillSon** | ~40+ especializaciones | Ej: "Java 21", "React", "PostgreSQL", "Hexagonal", "Agentes", etc. |
| **Experience** | 8 | Ordenadas por fecha más reciente primero (2026 → 2015). Todas con summaryPdf y descriptionItemsPdf |
| **Education** | 1 | UNINPAHU, Ingeniería de software, BS, 2012-2017 |
| **Blog** | 0+ | No visible en info-page, pero existe el endpoint |
| **BlogType** | 0+ | No visible en info-page |
| **FuturedProject** | 0+ | No visible en info-page, pero existe el endpoint |

### Observaciones para la migración:

1. **Orden actual de experiencias**: Ya están ordenadas por `yearStart DESC` en la respuesta actual. Esto confirma que no necesitan campo `order`.
2. **Campos HTML vs Texto**: Los campos `summary` y `summaryEng` contienen HTML en producción. Los campos `summaryPdf` y `summaryPdfEng` son texto plano. Esto confirma el patrón de separación.
3. **Campos vacíos**: Algunos `summaryPdf` tienen strings vacíos o solo whitespace (ej: experience id 6). La migración debe manejar esto correctamente.
4. **skillSons en Experience**: Cada experiencia tiene una lista de skillSons asociados (tecnologías usadas). Esta relación many-to-many debe preservarse.
5. **InfoPageResponse**: El endpoint público retorna `skills: SkillResponse[]` (lista plana de skills con sus skillSons), NO `SkillTypeResponse[]`. La jerarquía SkillType → Skill → SkillSon existe en los endpoints CRUD pero se aplana en la respuesta pública.
6. **AltchaChallenge**: Se incluye en cada respuesta de info-page. El challenge se genera dinámicamente por request.

### Script de migración requerido:

El script debe:
1. Leer la base de datos SQLite actual (Spring Boot).
2. Crear las nuevas tablas con campos adicionales (`order`, `deleted_at`).
3. Migrar los datos existentes preservando IDs y relaciones.
4. Asignar valores de `order` automáticos basados en el orden actual (ej: para experiencias, el order se puede derivar del orden por fecha).
5. No migrar datos a las nuevas entidades (Course, Certification, Language, Reference, CustomSection) ya que no existen en la BD actual.
6. Crear la tabla `PdfCache` vacía.
7. Crear la tabla de credenciales de admin si no existe, migrando el usuario/password actual.

---

## 13. Supuestos

| # | Supuesto | Validación Requerida |
|---|---|---|
| S1 | Los frontends actuales no requieren cambios de código para consumir las respuestas 200 de la nueva API. | Verificar con tests de integración contra los frontends. |
| S2 | El servicio de renderizado de PDFs (`hv-py-ms-render-cv`) estará disponible en la red interna Docker. | Verificar configuración de red Docker. |
| S3 | Un único usuario administrador es suficiente para el alcance actual. | Confirmado por el owner. |
| S4 | SQLite es suficiente para el volumen de datos de un portfolio personal. | Confirmado por el owner. |
| S5 | Las credenciales de AWS S3 y Telegram Bot están disponibles y activas. | Verificar acceso a las cuentas. |
| S6 | No se requiere backup automático de la base de datos SQLite en esta iteración. | Confirmar con el owner. |
| S7 | El patrón multi-idioma con campos duplicados (`*` y `*Eng`) se mantiene en la nueva versión. | Confirmado por el OpenAPI actual. |
| S8 | El patrón de separación HTML/Texto (`description` para portal, `descriptionPdf` para PDF) se mantiene y se replica en nuevas entidades. | Confirmado por el owner. |
| S9 | Las URLs de imágenes/videos ya existen en S3/YouTube antes de registrarse en el servicio. El servicio solo registra la referencia. | Confirmado por el OpenAPI actual (ImageUrlRequest recibe `file` como string, no como multipart). |
| S10 | El login requiere Altcha además de credenciales. | Confirmado por el OpenAPI actual. |

---

## 14. Handoff para Planner

### Resumen Funcional

Se requiere migrar el microservicio `ms-resume` de Spring Boot a Go 1.22+ manteniendo compatibilidad total con las respuestas 200 del contrato OpenAPI existente. El servicio gestiona el CRUD completo de datos de CV (13 entidades existentes + 5 nuevas), autenticación JWT con Altcha para un único administrador, generación de PDFs con **caché inteligente en S3** (invalidación por hash de datos), formulario de contacto con anti-spam, y gestión de referencias de multimedia. La persistencia es SQLite local con soft delete en todas las entidades.

**Caché de PDFs**: Los PDFs generados se almacenan en S3. Cada combinación `{language, template}` tiene un cache asociado a un hash SHA-256 de los datos del CV. Si los datos no cambian, se sirve el PDF desde S3. Si los datos cambian, se regenera el PDF, se elimina el anterior y se guarda el nuevo.

### Scope y Out of Scope

**Incluye**:
- Migración de 13 entidades existentes: BasicData, Home, Label, ImageUrl, VideoUrl, Blog, BlogType, SkillType, Skill, SkillSon, Experience, Education, FuturedProject.
- Migración de 9 endpoints públicos: login, info-page, curriculum, contact, challenge, blog (paginado), blog/{id}, blog-type, blog-type/{id}.
- 5 nuevas entidades: Course, Certification, Language, Reference, CustomSection.
- 2 nuevos endpoints públicos: `/public/templates`, `/public/languages`.
- Soft delete en todas las entidades.
- Límite de 2MB para imágenes con posible optimización.
- Patrón multi-idioma (campos ES/EN) y separación HTML/Texto.

**Excluye**: Multi-usuario, analytics, webhooks, versionado de CVs, migración de motor de base de datos, migración del formato de error.

### Actores, Permisos y Restricciones

- **Visitante público**: Solo lectura de endpoints `/public/*` + formulario de contacto.
- **Administrador (single-owner)**: CRUD completo vía JWT + Altcha en login.
- **Restricciones**: Respuestas 200 idénticas al OpenAPI, formato de error según skills configurados, credenciales en variables de entorno, password hasheado, soft delete universal.

### Flujos Principales y Alternos

1. Login (con Altcha) → JWT → CRUD autenticado.
2. Solicitud pública → Datos de CV → Respuesta JSON.
3. Descarga PDF → Transformación a RenderCV → HTTP interno → PDF binario.
4. Formulario contacto → Validación Altcha → Notificación Telegram.
5. Registro de multimedia → SQLite (URLs ya existentes en S3/YouTube).

### Entidades Funcionales y Datos Sensibles

- 13 entidades existentes con relaciones jerárquicas (SkillType → Skill → SkillSon, Blog → BlogType, FuturedProject → Experience, Home → ImageUrl + Label).
- 5 nuevas entidades siguiendo el mismo patrón multi-idioma y separación HTML/Texto.
- Datos sensibles: password del admin (hash), email del titular, email/teléfono de referencias.

### Integraciones Esperadas y Criticidad

- hv-rt-fr-portal (alta), hv-rt-fr-admin (alta), hv-py-ms-render-cv (alta), AWS S3 (**alta** - multimedia + caché de PDFs), Telegram Bot (baja).

### Criterios de Aceptación

11 grupos de criterios (CA1-CA11) cubriendo compatibilidad 200, autenticación, CRUD existentes, nuevas entidades, multimedia, **PDFs con caché inteligente**, contacto, base de datos, seguridad, performance y nuevos endpoints públicos.

### Preguntas Abiertas

- **Críticas**: _Ninguna. Todas las 9 preguntas críticas fueron resueltas._
- **No críticas (PA5, PA11, PA12)**: Planner puede tomar decisiones razonables y documentarlas como supuestos.

### Riesgos Funcionales que Planner Debe Resolver

1. **Compatibilidad de contrato**: Planner debe verificar que cada endpoint Go produce respuestas 200 idénticas al OpenAPI actual.
2. **Migración de datos**: Confirmado. Hay datos reales en producción (1 Home, 1 BasicData, 3 Labels, 6 SkillTypes, ~30 Skills, ~40 SkillSons, 8 Experiences, 1 Education, etc.). Se necesita un script de migración SQLite → SQLite que preserve IDs, relaciones y asigne campos `order` automáticamente. Ver sección 12.1 para detalles de los datos actuales.
3. **Transformación a schema RenderCV**: Planner debe definir el mapeo exacto de las 13+5 entidades del dominio al schema `CvData` de RenderCV (name, email, phone, location, headline, social_networks, sections, locale, design). El contrato está definido en `hv-py-ms-render-cv/docs/api/openapi.yaml`.
4. **Manejo de errores en integraciones**: Planner debe definir la estrategia de retry, timeout y fallback para cada integración.
5. **Orden de entidades**: Resuelto (PA10). Todas las entidades excepto Experience llevan campo `order`. Experience se ordena por `yearStart DESC`. Planner debe implementar la lógica de auto-asignación de `order` en CREATE y permitir reordenamiento vía UPDATE.
6. **Optimización de imágenes**: Si se implementa (PA11), Planner debe seleccionar la librería de procesamiento de imágenes en Go.
7. **Sanitización HTML**: Si el backend sanitiza (PA12), Planner debe seleccionar la librería de sanitización.
8. **Caché de PDFs**: Planner debe definir:
   - Algoritmo de serialización determinista para el hash (orden de campos, manejo de nulls, etc.).
   - Estrategia de concurrencia: si dos requests simultáneos detectan cache miss, ¿ambos generan el PDF o uno espera al otro?
   - Limpieza de PDFs huérfanos en S3 (cache en SQLite pero PDF eliminado de S3).
   - Bucket/prefix de S3 para PDFs (separado de multimedia o mismo bucket con prefijo diferente).

### Decisiones Pendientes para Planner

| Decisión | Contexto |
|---|---|
| Estructura de paquetes Go | Planner debe definir la estructura hexagonal (domain, ports, adapters, config, etc.) |
| Framework de migraciones | golang-migrate u otra herramienta compatible con SQLite |
| Validación de entrada | go-playground/validator u equivalente |
| Logging | Structured logging con zap, slog u otro |
| Testing strategy | Unit tests, integration tests, coverage target |
| Dockerfile | Multi-stage build, imagen mínima |
| Formato de error | Seguir el skill de error response configurado para Go |
| Soft delete implementation | Campo `deleted_at` (timestamp nullable) vs campo `is_deleted` (boolean) |
| Campo `order` | Todas las entidades excepto Experience llevan campo `order` INTEGER. Auto-asignación en CREATE (`MAX(order) + 1`). Reordenamiento vía UPDATE. Experience usa `yearStart DESC`. |
| Caché de PDFs | Hash SHA-256 de datos serializados. Clave S3: `pdfs/{language}/{template}/{hash}.pdf`. Entidad PdfCache en SQLite. Fallbacks si S3 no está disponible. Concurrencia en cache miss. |

---

## 15. Estado del Brief

| Campo | Valor |
|---|---|
| **Status** | `ready-for-planner` ✅ |
| **Preguntas críticas abiertas** | **0** - Todas resueltas |
| **Preguntas no críticas abiertas** | 3 (PA5, PA11, PA12) - Planner puede decidir |
| **Preguntas resueltas** | **7** (PA1, PA3, PA6, PA7, PA8, PA9, PA10) |
| **Próximo paso** | Planner genera specs SDD |
