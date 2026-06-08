package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.BlogTypeRepository = (*BlogTypeRepo)(nil)

type BlogTypeRepo struct{ db *DB }

func NewBlogTypeRepo(db *DB) *BlogTypeRepo { return &BlogTypeRepo{db: db} }

func (r *BlogTypeRepo) List(ctx context.Context) ([]entity.BlogType, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM blog_type WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.BlogType
	for rows.Next() {
		var e entity.BlogType
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *BlogTypeRepo) GetByID(ctx context.Context, id int64) (*entity.BlogType, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM blog_type WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.BlogType
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("blog_type: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *BlogTypeRepo) Create(ctx context.Context, data *entity.BlogType) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO blog_type (name, name_eng, "order", created_at, updated_at)
		VALUES (?, ?, COALESCE((SELECT MAX("order") FROM blog_type WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *BlogTypeRepo) Update(ctx context.Context, data *entity.BlogType) error {
	query := `UPDATE blog_type SET name=?, name_eng=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, nowUTC(), data.ID)
	return err
}

func (r *BlogTypeRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE blog_type SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
