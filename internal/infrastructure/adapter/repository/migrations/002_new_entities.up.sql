CREATE TABLE IF NOT EXISTS course (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    institution TEXT NOT NULL,
    institution_eng TEXT NOT NULL,
    completion_date TEXT NOT NULL,
    description TEXT,
    description_eng TEXT,
    summary_pdf TEXT,
    summary_pdf_eng TEXT,
    certificate_url TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS certification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    issuing_organization TEXT NOT NULL,
    issuing_organization_eng TEXT NOT NULL,
    issue_date TEXT NOT NULL,
    expiration_date TEXT,
    verification_url TEXT,
    credential_id TEXT,
    description TEXT,
    description_eng TEXT,
    summary_pdf TEXT,
    summary_pdf_eng TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS language (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    language TEXT NOT NULL,
    language_eng TEXT NOT NULL,
    reading_level TEXT NOT NULL,
    writing_level TEXT NOT NULL,
    speaking_level TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS reference (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    full_name TEXT NOT NULL,
    position TEXT NOT NULL,
    company TEXT,
    company_eng TEXT,
    email TEXT,
    phone TEXT,
    relationship TEXT,
    relationship_eng TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS custom_section (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    title_eng TEXT NOT NULL,
    content TEXT,
    content_eng TEXT,
    summary_pdf TEXT,
    summary_pdf_eng TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    visible INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);
