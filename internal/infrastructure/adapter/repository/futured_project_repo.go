package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.FuturedProjectRepository = (*FuturedProjectRepo)(nil)

type FuturedProjectRepo struct{ db *DB }

func NewFuturedProjectRepo(db *DB) *FuturedProjectRepo { return &FuturedProjectRepo{db: db} }

func (r *FuturedProjectRepo) List(ctx context.Context) ([]entity.FuturedProject, error) {
	query := `SELECT id, name, name_eng, description_short, description, description_short_eng, description_eng,
		experience_id, image_list_url_id, image_url_id, "order",
		created_at, updated_at, deleted_at FROM futured_project WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.FuturedProject
	for rows.Next() {
		var e entity.FuturedProject
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.DescriptionShort, &e.Description,
			&e.DescriptionShortEng, &e.DescriptionEng, &e.ExperienceID,
			&e.ImageListURLID, &e.ImageURLID, &e.Order,
			&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *FuturedProjectRepo) GetByID(ctx context.Context, id int64) (*entity.FuturedProject, error) {
	query := `SELECT id, name, name_eng, description_short, description, description_short_eng, description_eng,
		experience_id, image_list_url_id, image_url_id, "order",
		created_at, updated_at, deleted_at FROM futured_project WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.FuturedProject
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.DescriptionShort, &e.Description,
		&e.DescriptionShortEng, &e.DescriptionEng, &e.ExperienceID,
		&e.ImageListURLID, &e.ImageURLID, &e.Order,
		&e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("futured_project not found")
	}
	return &e, err
}

func (r *FuturedProjectRepo) Create(ctx context.Context, data *entity.FuturedProject) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO futured_project (name, name_eng, description_short, description, description_short_eng, description_eng,
		experience_id, image_list_url_id, image_url_id, "order", created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM futured_project WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Name, data.NameEng, data.DescriptionShort, data.Description,
		data.DescriptionShortEng, data.DescriptionEng, data.ExperienceID,
		data.ImageListURLID, data.ImageURLID, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *FuturedProjectRepo) Update(ctx context.Context, data *entity.FuturedProject) error {
	query := `UPDATE futured_project SET name=?, name_eng=?, description_short=?, description=?,
		description_short_eng=?, description_eng=?, experience_id=?, image_list_url_id=?, image_url_id=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Name, data.NameEng, data.DescriptionShort, data.Description,
		data.DescriptionShortEng, data.DescriptionEng, data.ExperienceID,
		data.ImageListURLID, data.ImageURLID, nowUTC(), data.ID)
	return err
}

func (r *FuturedProjectRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE futured_project SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
