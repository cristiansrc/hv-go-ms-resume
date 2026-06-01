CREATE TABLE IF NOT EXISTS basic_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    others_name TEXT,
    first_surname TEXT NOT NULL,
    others_surname TEXT,
    date_birth TEXT NOT NULL,
    located TEXT,
    located_eng TEXT,
    start_working_date TEXT,
    greeting TEXT,
    greeting_eng TEXT,
    email TEXT NOT NULL,
    instagram TEXT,
    linkedin TEXT,
    x TEXT,
    github TEXT,
    description TEXT,
    description_eng TEXT,
    description_pdf TEXT,
    description_pdf_eng TEXT,
    wrapper TEXT,
    wrapper_eng TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS home (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    greeting TEXT NOT NULL,
    greeting_eng TEXT NOT NULL,
    image_url_id INTEGER REFERENCES image_url(id),
    button_work_label TEXT NOT NULL,
    button_work_label_eng TEXT NOT NULL,
    button_contact_label TEXT NOT NULL,
    button_contact_label_eng TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS label (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS image_url (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS video_url (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS blog_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS blog (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    title_eng TEXT NOT NULL,
    clean_url_title TEXT,
    description_short TEXT NOT NULL,
    description TEXT NOT NULL,
    description_short_eng TEXT NOT NULL,
    description_eng TEXT NOT NULL,
    image_url_id INTEGER REFERENCES image_url(id),
    video_url_id INTEGER REFERENCES video_url(id),
    blog_type_id INTEGER REFERENCES blog_type(id),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS skill_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS skill (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS skill_son (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS experience (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year_start TEXT NOT NULL,
    year_end TEXT,
    company TEXT NOT NULL,
    location TEXT,
    location_eng TEXT,
    position TEXT,
    position_eng TEXT,
    summary TEXT,
    summary_eng TEXT,
    summary_pdf TEXT,
    summary_pdf_eng TEXT,
    description_items_pdf TEXT,
    description_items_pdf_eng TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS education (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    institution TEXT NOT NULL,
    area TEXT NOT NULL,
    area_eng TEXT NOT NULL,
    degree TEXT NOT NULL,
    degree_eng TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT,
    location TEXT NOT NULL,
    location_eng TEXT NOT NULL,
    highlights TEXT,
    highlights_eng TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS futured_project (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    name_eng TEXT NOT NULL,
    description_short TEXT NOT NULL,
    description TEXT NOT NULL,
    description_short_eng TEXT NOT NULL,
    description_eng TEXT NOT NULL,
    experience_id INTEGER NOT NULL REFERENCES experience(id),
    image_list_url_id INTEGER REFERENCES image_url(id),
    image_url_id INTEGER REFERENCES image_url(id),
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    deleted_at TEXT
);

-- Join tables
CREATE TABLE IF NOT EXISTS home_label (
    home_id INTEGER NOT NULL REFERENCES home(id),
    label_id INTEGER NOT NULL REFERENCES label(id),
    PRIMARY KEY (home_id, label_id)
);

CREATE TABLE IF NOT EXISTS skill_type_skill (
    skill_type_id INTEGER NOT NULL REFERENCES skill_type(id),
    skill_id INTEGER NOT NULL REFERENCES skill(id),
    PRIMARY KEY (skill_type_id, skill_id)
);

CREATE TABLE IF NOT EXISTS skill_skill_son (
    skill_id INTEGER NOT NULL REFERENCES skill(id),
    skill_son_id INTEGER NOT NULL REFERENCES skill_son(id),
    PRIMARY KEY (skill_id, skill_son_id)
);

CREATE TABLE IF NOT EXISTS experience_skill_son (
    experience_id INTEGER NOT NULL REFERENCES experience(id),
    skill_son_id INTEGER NOT NULL REFERENCES skill_son(id),
    PRIMARY KEY (experience_id, skill_son_id)
);

-- Support tables
CREATE TABLE IF NOT EXISTS user_credentials (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);

CREATE TABLE IF NOT EXISTS pdf_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    language TEXT NOT NULL,
    template TEXT NOT NULL,
    data_hash TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    file_size INTEGER,
    UNIQUE(language, template)
);
