package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.CustomSectionRepository = (*CustomSectionRepo)(nil)

type CustomSectionRepo struct{ db *DB }

func NewCustomSectionRepo(db *DB) *CustomSectionRepo { return &CustomSectionRepo{db: db} }

func (r *CustomSectionRepo) List(ctx context.Context) ([]entity.CustomSection, error) {
	query := `SELECT id, title, title_eng, content, content_eng, summary_pdf, summary_pdf_eng,
		"order", visible, created_at, updated_at, deleted_at
		FROM custom_section WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.CustomSection
	for rows.Next() {
		var e entity.CustomSection
		var visible int
		if err := rows.Scan(&e.ID, &e.Title, &e.TitleEng, &e.Content, &e.ContentEng,
			&e.SummaryPdf, &e.SummaryPdfEng, &e.Order, &visible,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		e.Visible = visible == 1
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *CustomSectionRepo) GetByID(ctx context.Context, id int64) (*entity.CustomSection, error) {
	query := `SELECT id, title, title_eng, content, content_eng, summary_pdf, summary_pdf_eng,
		"order", visible, created_at, updated_at, deleted_at
		FROM custom_section WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.CustomSection
	var visible int
	err := row.Scan(&e.ID, &e.Title, &e.TitleEng, &e.Content, &e.ContentEng,
		&e.SummaryPdf, &e.SummaryPdfEng, &e.Order, &visible,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("custom_section: %w", entity.ErrNotFound)
	}
	e.Visible = visible == 1
	return &e, err
}

func (r *CustomSectionRepo) Create(ctx context.Context, data *entity.CustomSection) (int64, error) {
	now := nowUTC()
	visible := 0
	if data.Visible {
		visible = 1
	}
	query := `INSERT INTO custom_section (title, title_eng, content, content_eng, summary_pdf, summary_pdf_eng,
		"order", visible, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM custom_section WHERE deleted_at IS NULL), 0) + 1, ?, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Title, data.TitleEng, data.Content, data.ContentEng,
		data.SummaryPdf, data.SummaryPdfEng, visible, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CustomSectionRepo) Update(ctx context.Context, data *entity.CustomSection) error {
	visible := 0
	if data.Visible {
		visible = 1
	}
	query := `UPDATE custom_section SET title=?, title_eng=?, content=?, content_eng=?,
		summary_pdf=?, summary_pdf_eng=?, visible=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Title, data.TitleEng, data.Content, data.ContentEng,
		data.SummaryPdf, data.SummaryPdfEng, visible, nowUTC(), data.ID)
	return err
}

func (r *CustomSectionRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE custom_section SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
