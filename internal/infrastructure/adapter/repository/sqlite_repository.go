package repository

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB represents the SQLite database connection wrapper.
type DB struct {
	conn   *sql.DB
	logger *slog.Logger
}

// NewDB opens a SQLite database, applies migrations, and returns the DB wrapper.
func NewDB(dbPath string, logger *slog.Logger) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn, logger: logger}

	if err := db.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func (db *DB) runMigrations() error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	driver, err := sqlite3.WithInstance(db.conn, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	db.logger.Info("database migrations applied successfully")
	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying sql.DB.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// nowUTC returns current UTC timestamp in RFC3339 format.
func nowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// --- Repository implementations ---

// Compile-time interface checks.
var (
	_ output.BasicDataRepository        = (*BasicDataRepo)(nil)
	_ output.HomeRepository             = (*HomeRepo)(nil)
	_ output.LabelRepository            = (*LabelRepo)(nil)
	_ output.ImageUrlRepository         = (*ImageUrlRepo)(nil)
	_ output.VideoUrlRepository         = (*VideoUrlRepo)(nil)
	_ output.BlogRepository             = (*BlogRepo)(nil)
	_ output.BlogTypeRepository         = (*BlogTypeRepo)(nil)
	_ output.SkillTypeRepository        = (*SkillTypeRepo)(nil)
	_ output.SkillRepository            = (*SkillRepo)(nil)
	_ output.SkillSonRepository         = (*SkillSonRepo)(nil)
	_ output.ExperienceRepository       = (*ExperienceRepo)(nil)
	_ output.EducationRepository        = (*EducationRepo)(nil)
	_ output.FuturedProjectRepository   = (*FuturedProjectRepo)(nil)
	_ output.CourseRepository           = (*CourseRepo)(nil)
	_ output.CertificationRepository    = (*CertificationRepo)(nil)
	_ output.LanguageRepository         = (*LanguageRepo)(nil)
	_ output.ReferenceRepository        = (*ReferenceRepo)(nil)
	_ output.CustomSectionRepository    = (*CustomSectionRepo)(nil)
	_ output.UserCredentialsRepository  = (*UserCredentialsRepo)(nil)
	_ output.PdfCacheRepository         = (*PdfCacheRepo)(nil)
)

// --- BasicDataRepo ---

type BasicDataRepo struct{ db *DB }

func NewBasicDataRepo(db *DB) *BasicDataRepo { return &BasicDataRepo{db: db} }

func (r *BasicDataRepo) GetByID(ctx context.Context, id int64) (*entity.BasicData, error) {
	query := `SELECT id, first_name, others_name, first_surname, others_surname, date_birth,
		located, located_eng, start_working_date, greeting, greeting_eng, email,
		instagram, linkedin, x, github, description, description_eng,
		description_pdf, description_pdf_eng, wrapper, wrapper_eng,
		created_at, updated_at, deleted_at
		FROM basic_data WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	return scanBasicData(row)
}

func (r *BasicDataRepo) Update(ctx context.Context, data *entity.BasicData) error {
	query := `UPDATE basic_data SET first_name=?, others_name=?, first_surname=?, others_surname=?,
		date_birth=?, located=?, located_eng=?, start_working_date=?, greeting=?, greeting_eng=?,
		email=?, instagram=?, linkedin=?, x=?, github=?, description=?, description_eng=?,
		description_pdf=?, description_pdf_eng=?, wrapper=?, wrapper_eng=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.FirstName, data.OthersName, data.FirstSurname, data.OthersSurname,
		data.DateBirth, data.Located, data.LocatedEng, data.StartWorkingDate,
		data.Greeting, data.GreetingEng, data.Email, data.Instagram, data.Linkedin,
		data.X, data.Github, data.Description, data.DescriptionEng,
		data.DescriptionPdf, data.DescriptionPdfEng, data.Wrapper, data.WrapperEng,
		nowUTC(), data.ID)
	return err
}

func scanBasicData(row *sql.Row) (*entity.BasicData, error) {
	var e entity.BasicData
	err := row.Scan(&e.ID, &e.FirstName, &e.OthersName, &e.FirstSurname, &e.OthersSurname,
		&e.DateBirth, &e.Located, &e.LocatedEng, &e.StartWorkingDate,
		&e.Greeting, &e.GreetingEng, &e.Email, &e.Instagram, &e.Linkedin,
		&e.X, &e.Github, &e.Description, &e.DescriptionEng,
		&e.DescriptionPdf, &e.DescriptionPdfEng, &e.Wrapper, &e.WrapperEng,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("basic_data not found")
	}
	return &e, err
}

// --- HomeRepo ---

type HomeRepo struct{ db *DB }

func NewHomeRepo(db *DB) *HomeRepo { return &HomeRepo{db: db} }

func (r *HomeRepo) GetByID(ctx context.Context, id int64) (*entity.Home, error) {
	query := `SELECT id, greeting, greeting_eng, image_url_id,
		button_work_label, button_work_label_eng, button_contact_label, button_contact_label_eng,
		created_at, updated_at, deleted_at FROM home WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Home
	err := row.Scan(&e.ID, &e.Greeting, &e.GreetingEng, &e.ImageURLID,
		&e.ButtonWorkLabel, &e.ButtonWorkLabelEng, &e.ButtonContactLabel, &e.ButtonContactLabelEng,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("home not found")
	}
	return &e, err
}

func (r *HomeRepo) Update(ctx context.Context, data *entity.Home) error {
	query := `UPDATE home SET greeting=?, greeting_eng=?, image_url_id=?,
		button_work_label=?, button_work_label_eng=?, button_contact_label=?, button_contact_label_eng=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Greeting, data.GreetingEng, data.ImageURLID,
		data.ButtonWorkLabel, data.ButtonWorkLabelEng, data.ButtonContactLabel, data.ButtonContactLabelEng,
		nowUTC(), data.ID)
	return err
}

// --- LabelRepo ---

type LabelRepo struct{ db *DB }

func NewLabelRepo(db *DB) *LabelRepo { return &LabelRepo{db: db} }

func (r *LabelRepo) List(ctx context.Context) ([]entity.Label, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM label WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLabels(rows)
}

func (r *LabelRepo) GetByID(ctx context.Context, id int64) (*entity.Label, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM label WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	return scanLabel(row)
}

func (r *LabelRepo) Create(ctx context.Context, data *entity.Label) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO label (name, name_eng, "order", created_at, updated_at)
		VALUES (?, ?, COALESCE((SELECT MAX("order") FROM label WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *LabelRepo) Update(ctx context.Context, data *entity.Label) error {
	query := `UPDATE label SET name=?, name_eng=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, nowUTC(), data.ID)
	return err
}

func (r *LabelRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE label SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}

func scanLabel(row *sql.Row) (*entity.Label, error) {
	var e entity.Label
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("label not found")
	}
	return &e, err
}

func scanLabels(rows *sql.Rows) ([]entity.Label, error) {
	var items []entity.Label
	for rows.Next() {
		var e entity.Label
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

// --- ImageUrlRepo ---

type ImageUrlRepo struct{ db *DB }

func NewImageUrlRepo(db *DB) *ImageUrlRepo { return &ImageUrlRepo{db: db} }

func (r *ImageUrlRepo) List(ctx context.Context) ([]entity.ImageUrl, error) {
	query := `SELECT id, name, name_eng, url, created_at, updated_at, deleted_at
		FROM image_url WHERE deleted_at IS NULL ORDER BY id ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.ImageUrl
	for rows.Next() {
		var e entity.ImageUrl
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.URL, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *ImageUrlRepo) GetByID(ctx context.Context, id int64) (*entity.ImageUrl, error) {
	query := `SELECT id, name, name_eng, url, created_at, updated_at, deleted_at
		FROM image_url WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.ImageUrl
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.URL, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("image_url not found")
	}
	return &e, err
}

func (r *ImageUrlRepo) Create(ctx context.Context, data *entity.ImageUrl) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO image_url (name, name_eng, url, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, data.URL, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ImageUrlRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE image_url SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}

// --- VideoUrlRepo ---

type VideoUrlRepo struct{ db *DB }

func NewVideoUrlRepo(db *DB) *VideoUrlRepo { return &VideoUrlRepo{db: db} }

func (r *VideoUrlRepo) List(ctx context.Context) ([]entity.VideoUrl, error) {
	query := `SELECT id, name, name_eng, url, created_at, updated_at, deleted_at
		FROM video_url WHERE deleted_at IS NULL ORDER BY id ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.VideoUrl
	for rows.Next() {
		var e entity.VideoUrl
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.URL, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *VideoUrlRepo) GetByID(ctx context.Context, id int64) (*entity.VideoUrl, error) {
	query := `SELECT id, name, name_eng, url, created_at, updated_at, deleted_at
		FROM video_url WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.VideoUrl
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.URL, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("video_url not found")
	}
	return &e, err
}

func (r *VideoUrlRepo) Create(ctx context.Context, data *entity.VideoUrl) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO video_url (name, name_eng, url, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, data.URL, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *VideoUrlRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE video_url SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}

// --- Helper: Generic CRUD repository for entities with full 5 operations ---

// crudRepo provides generic CRUD operations for entities with list, get, create, update, delete.
type crudRepo struct {
	db          *DB
	table       string
	columns     string
	orderBy     string
	scanFn      func(*sql.Row) (interface{}, error)
	scanListFn  func(*sql.Rows) (interface{}, error)
	insertCols  string
	insertVals  string
	insertArgs  func(interface{}) []interface{}
	updateSets  string
	updateArgs  func(interface{}) []interface{}
}

// --- BlogRepo ---

type BlogRepo struct{ db *DB }

func NewBlogRepo(db *DB) *BlogRepo { return &BlogRepo{db: db} }

func (r *BlogRepo) List(ctx context.Context) ([]entity.Blog, error) {
	query := `SELECT id, title, title_eng, clean_url_title, description_short, description,
		description_short_eng, description_eng, image_url_id, video_url_id, blog_type_id,
		created_at, updated_at, deleted_at FROM blog WHERE deleted_at IS NULL ORDER BY id DESC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Blog
	for rows.Next() {
		var e entity.Blog
		err := rows.Scan(&e.ID, &e.Title, &e.TitleEng, &e.CleanURLTitle,
			&e.DescriptionShort, &e.Description, &e.DescriptionShortEng, &e.DescriptionEng,
			&e.ImageURLID, &e.VideoURLID, &e.BlogTypeID,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *BlogRepo) GetByID(ctx context.Context, id int64) (*entity.Blog, error) {
	query := `SELECT id, title, title_eng, clean_url_title, description_short, description,
		description_short_eng, description_eng, image_url_id, video_url_id, blog_type_id,
		created_at, updated_at, deleted_at FROM blog WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Blog
	err := row.Scan(&e.ID, &e.Title, &e.TitleEng, &e.CleanURLTitle,
		&e.DescriptionShort, &e.Description, &e.DescriptionShortEng, &e.DescriptionEng,
		&e.ImageURLID, &e.VideoURLID, &e.BlogTypeID,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("blog not found")
	}
	return &e, err
}

func (r *BlogRepo) Create(ctx context.Context, data *entity.Blog) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO blog (title, title_eng, clean_url_title, description_short, description,
		description_short_eng, description_eng, image_url_id, video_url_id, blog_type_id,
		created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Title, data.TitleEng, data.CleanURLTitle,
		data.DescriptionShort, data.Description, data.DescriptionShortEng, data.DescriptionEng,
		data.ImageURLID, data.VideoURLID, data.BlogTypeID, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *BlogRepo) Update(ctx context.Context, data *entity.Blog) error {
	query := `UPDATE blog SET title=?, title_eng=?, clean_url_title=?, description_short=?, description=?,
		description_short_eng=?, description_eng=?, image_url_id=?, video_url_id=?, blog_type_id=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Title, data.TitleEng, data.CleanURLTitle,
		data.DescriptionShort, data.Description, data.DescriptionShortEng, data.DescriptionEng,
		data.ImageURLID, data.VideoURLID, data.BlogTypeID, nowUTC(), data.ID)
	return err
}

func (r *BlogRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE blog SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}

// --- Entity repositories are defined in separate files for maintainability ---
// See: blog_type_repo.go, skill_type_repo.go, skill_repo.go, skill_son_repo.go,
//      education_repo.go, futured_project_repo.go, course_repo.go,
//      certification_repo.go, language_repo.go, reference_repo.go, custom_section_repo.go

// --- ExperienceRepo ---

type ExperienceRepo struct{ db *DB }

func NewExperienceRepo(db *DB) *ExperienceRepo { return &ExperienceRepo{db: db} }

func (r *ExperienceRepo) List(ctx context.Context) ([]entity.Experience, error) {
	query := `SELECT id, year_start, year_end, company, location, location_eng, position, position_eng,
		summary, summary_eng, summary_pdf, summary_pdf_eng, description_items_pdf, description_items_pdf_eng,
		created_at, updated_at, deleted_at FROM experience WHERE deleted_at IS NULL ORDER BY year_start DESC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Experience
	for rows.Next() {
		var e entity.Experience
		err := rows.Scan(&e.ID, &e.YearStart, &e.YearEnd, &e.Company,
			&e.Location, &e.LocationEng, &e.Position, &e.PositionEng,
			&e.Summary, &e.SummaryEng, &e.SummaryPdf, &e.SummaryPdfEng,
			&e.DescriptionItemsPdf, &e.DescriptionItemsPdfEng,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *ExperienceRepo) GetByID(ctx context.Context, id int64) (*entity.Experience, error) {
	query := `SELECT id, year_start, year_end, company, location, location_eng, position, position_eng,
		summary, summary_eng, summary_pdf, summary_pdf_eng, description_items_pdf, description_items_pdf_eng,
		created_at, updated_at, deleted_at FROM experience WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Experience
	err := row.Scan(&e.ID, &e.YearStart, &e.YearEnd, &e.Company,
		&e.Location, &e.LocationEng, &e.Position, &e.PositionEng,
		&e.Summary, &e.SummaryEng, &e.SummaryPdf, &e.SummaryPdfEng,
		&e.DescriptionItemsPdf, &e.DescriptionItemsPdfEng,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("experience not found")
	}
	return &e, err
}

func (r *ExperienceRepo) Create(ctx context.Context, data *entity.Experience) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO experience (year_start, year_end, company, location, location_eng,
		position, position_eng, summary, summary_eng, summary_pdf, summary_pdf_eng,
		description_items_pdf, description_items_pdf_eng, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.YearStart, data.YearEnd, data.Company, data.Location, data.LocationEng,
		data.Position, data.PositionEng, data.Summary, data.SummaryEng,
		data.SummaryPdf, data.SummaryPdfEng, data.DescriptionItemsPdf, data.DescriptionItemsPdfEng,
		now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ExperienceRepo) Update(ctx context.Context, data *entity.Experience) error {
	query := `UPDATE experience SET year_start=?, year_end=?, company=?, location=?, location_eng=?,
		position=?, position_eng=?, summary=?, summary_eng=?, summary_pdf=?, summary_pdf_eng=?,
		description_items_pdf=?, description_items_pdf_eng=?, updated_at=?
		WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.YearStart, data.YearEnd, data.Company, data.Location, data.LocationEng,
		data.Position, data.PositionEng, data.Summary, data.SummaryEng,
		data.SummaryPdf, data.SummaryPdfEng, data.DescriptionItemsPdf, data.DescriptionItemsPdfEng,
		nowUTC(), data.ID)
	return err
}

func (r *ExperienceRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE experience SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}

// --- UserCredentialsRepo ---

type UserCredentialsRepo struct{ db *DB }

func NewUserCredentialsRepo(db *DB) *UserCredentialsRepo { return &UserCredentialsRepo{db: db} }

func (r *UserCredentialsRepo) GetByUsername(ctx context.Context, username string) (*entity.UserCredentials, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at FROM user_credentials WHERE username = ?`
	row := r.db.conn.QueryRowContext(ctx, query, username)
	var e entity.UserCredentials
	err := row.Scan(&e.ID, &e.Username, &e.PasswordHash, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}
	return &e, err
}

// --- PdfCacheRepo ---

type PdfCacheRepo struct{ db *DB }

func NewPdfCacheRepo(db *DB) *PdfCacheRepo { return &PdfCacheRepo{db: db} }

func (r *PdfCacheRepo) GetByLanguageAndTemplate(ctx context.Context, language, template string) (*entity.PdfCache, error) {
	query := `SELECT id, language, template, data_hash, s3_key, created_at, file_size
		FROM pdf_cache WHERE language = ? AND template = ?`
	row := r.db.conn.QueryRowContext(ctx, query, language, template)
	var e entity.PdfCache
	err := row.Scan(&e.ID, &e.Language, &e.Template, &e.DataHash, &e.S3Key, &e.CreatedAt, &e.FileSize)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &e, err
}

func (r *PdfCacheRepo) Upsert(ctx context.Context, data *entity.PdfCache) error {
	query := `INSERT INTO pdf_cache (language, template, data_hash, s3_key, created_at, file_size)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(language, template) DO UPDATE SET data_hash=excluded.data_hash, s3_key=excluded.s3_key, file_size=excluded.file_size`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Language, data.Template, data.DataHash, data.S3Key, nowUTC(), data.FileSize)
	return err
}

func (r *PdfCacheRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pdf_cache WHERE id = ?`
	_, err := r.db.conn.ExecContext(ctx, query, id)
	return err
}
