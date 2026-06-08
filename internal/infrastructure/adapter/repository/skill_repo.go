package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.SkillRepository = (*SkillRepo)(nil)

type SkillRepo struct{ db *DB }

func NewSkillRepo(db *DB) *SkillRepo { return &SkillRepo{db: db} }

func (r *SkillRepo) List(ctx context.Context) ([]entity.Skill, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM skill WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Skill
	for rows.Next() {
		var e entity.Skill
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *SkillRepo) GetByID(ctx context.Context, id int64) (*entity.Skill, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM skill WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Skill
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("skill: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *SkillRepo) Create(ctx context.Context, data *entity.Skill) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO skill (name, name_eng, "order", created_at, updated_at)
		VALUES (?, ?, COALESCE((SELECT MAX("order") FROM skill WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *SkillRepo) Update(ctx context.Context, data *entity.Skill) error {
	query := `UPDATE skill SET name=?, name_eng=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, nowUTC(), data.ID)
	return err
}

func (r *SkillRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE skill SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
