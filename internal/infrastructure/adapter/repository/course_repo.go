package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.CourseRepository = (*CourseRepo)(nil)

type CourseRepo struct{ db *DB }

func NewCourseRepo(db *DB) *CourseRepo { return &CourseRepo{db: db} }

func (r *CourseRepo) List(ctx context.Context) ([]entity.Course, error) {
	query := `SELECT id, name, name_eng, institution, institution_eng, completion_date,
		description, description_eng, summary_pdf, summary_pdf_eng, certificate_url, "order",
		created_at, updated_at, deleted_at FROM course WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Course
	for rows.Next() {
		var e entity.Course
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Institution, &e.InstitutionEng,
			&e.CompletionDate, &e.Description, &e.DescriptionEng, &e.SummaryPdf, &e.SummaryPdfEng,
			&e.CertificateURL, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *CourseRepo) GetByID(ctx context.Context, id int64) (*entity.Course, error) {
	query := `SELECT id, name, name_eng, institution, institution_eng, completion_date,
		description, description_eng, summary_pdf, summary_pdf_eng, certificate_url, "order",
		created_at, updated_at, deleted_at FROM course WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Course
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.Institution, &e.InstitutionEng,
		&e.CompletionDate, &e.Description, &e.DescriptionEng, &e.SummaryPdf, &e.SummaryPdfEng,
		&e.CertificateURL, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("course: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *CourseRepo) Create(ctx context.Context, data *entity.Course) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO course (name, name_eng, institution, institution_eng, completion_date,
		description, description_eng, summary_pdf, summary_pdf_eng, certificate_url, "order", created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM course WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Name, data.NameEng, data.Institution, data.InstitutionEng, data.CompletionDate,
		data.Description, data.DescriptionEng, data.SummaryPdf, data.SummaryPdfEng, data.CertificateURL,
		now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CourseRepo) Update(ctx context.Context, data *entity.Course) error {
	query := `UPDATE course SET name=?, name_eng=?, institution=?, institution_eng=?, completion_date=?,
		description=?, description_eng=?, summary_pdf=?, summary_pdf_eng=?, certificate_url=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Name, data.NameEng, data.Institution, data.InstitutionEng, data.CompletionDate,
		data.Description, data.DescriptionEng, data.SummaryPdf, data.SummaryPdfEng, data.CertificateURL,
		nowUTC(), data.ID)
	return err
}

func (r *CourseRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE course SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
