package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.SkillTypeRepository = (*SkillTypeRepo)(nil)

type SkillTypeRepo struct{ db *DB }

func NewSkillTypeRepo(db *DB) *SkillTypeRepo { return &SkillTypeRepo{db: db} }

func (r *SkillTypeRepo) List(ctx context.Context) ([]entity.SkillType, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM skill_type WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.SkillType
	for rows.Next() {
		var e entity.SkillType
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *SkillTypeRepo) GetByID(ctx context.Context, id int64) (*entity.SkillType, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM skill_type WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.SkillType
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("skill_type: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *SkillTypeRepo) Create(ctx context.Context, data *entity.SkillType) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO skill_type (name, name_eng, "order", created_at, updated_at)
		VALUES (?, ?, COALESCE((SELECT MAX("order") FROM skill_type WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *SkillTypeRepo) Update(ctx context.Context, data *entity.SkillType) error {
	query := `UPDATE skill_type SET name=?, name_eng=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, nowUTC(), data.ID)
	return err
}

func (r *SkillTypeRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE skill_type SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
