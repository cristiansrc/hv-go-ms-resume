package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.EducationRepository = (*EducationRepo)(nil)

type EducationRepo struct{ db *DB }

func NewEducationRepo(db *DB) *EducationRepo { return &EducationRepo{db: db} }

func (r *EducationRepo) List(ctx context.Context) ([]entity.Education, error) {
	query := `SELECT id, institution, area, area_eng, degree, degree_eng, start_date, end_date,
		location, location_eng, highlights, highlights_eng, "order",
		created_at, updated_at, deleted_at FROM education WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Education
	for rows.Next() {
		var e entity.Education
		if err := rows.Scan(&e.ID, &e.Institution, &e.Area, &e.AreaEng, &e.Degree, &e.DegreeEng,
			&e.StartDate, &e.EndDate, &e.Location, &e.LocationEng, &e.Highlights, &e.HighlightsEng,
			&e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *EducationRepo) GetByID(ctx context.Context, id int64) (*entity.Education, error) {
	query := `SELECT id, institution, area, area_eng, degree, degree_eng, start_date, end_date,
		location, location_eng, highlights, highlights_eng, "order",
		created_at, updated_at, deleted_at FROM education WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Education
	err := row.Scan(&e.ID, &e.Institution, &e.Area, &e.AreaEng, &e.Degree, &e.DegreeEng,
		&e.StartDate, &e.EndDate, &e.Location, &e.LocationEng, &e.Highlights, &e.HighlightsEng,
		&e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("education: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *EducationRepo) Create(ctx context.Context, data *entity.Education) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO education (institution, area, area_eng, degree, degree_eng, start_date, end_date,
		location, location_eng, highlights, highlights_eng, "order", created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM education WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Institution, data.Area, data.AreaEng, data.Degree, data.DegreeEng,
		data.StartDate, data.EndDate, data.Location, data.LocationEng,
		data.Highlights, data.HighlightsEng, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *EducationRepo) Update(ctx context.Context, data *entity.Education) error {
	query := `UPDATE education SET institution=?, area=?, area_eng=?, degree=?, degree_eng=?,
		start_date=?, end_date=?, location=?, location_eng=?, highlights=?, highlights_eng=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Institution, data.Area, data.AreaEng, data.Degree, data.DegreeEng,
		data.StartDate, data.EndDate, data.Location, data.LocationEng,
		data.Highlights, data.HighlightsEng, nowUTC(), data.ID)
	return err
}

func (r *EducationRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE education SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
