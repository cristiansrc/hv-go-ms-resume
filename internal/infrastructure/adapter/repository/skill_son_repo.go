package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.SkillSonRepository = (*SkillSonRepo)(nil)

type SkillSonRepo struct{ db *DB }

func NewSkillSonRepo(db *DB) *SkillSonRepo { return &SkillSonRepo{db: db} }

func (r *SkillSonRepo) List(ctx context.Context) ([]entity.SkillSon, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM skill_son WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.SkillSon
	for rows.Next() {
		var e entity.SkillSon
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *SkillSonRepo) ListBySkillID(ctx context.Context, skillID int64) ([]entity.SkillSon, error) {
	query := `SELECT ss.id, ss.name, ss.name_eng, ss."order", ss.created_at, ss.updated_at, ss.deleted_at
		FROM skill_son ss
		JOIN skill_skill_son sss ON ss.id = sss.skill_son_id
		WHERE sss.skill_id = ? AND ss.deleted_at IS NULL
		ORDER BY ss."order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query, skillID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.SkillSon
	for rows.Next() {
		var e entity.SkillSon
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *SkillSonRepo) GetByID(ctx context.Context, id int64) (*entity.SkillSon, error) {
	query := `SELECT id, name, name_eng, "order", created_at, updated_at, deleted_at
		FROM skill_son WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.SkillSon
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("skill_son: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *SkillSonRepo) Create(ctx context.Context, data *entity.SkillSon) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO skill_son (name, name_eng, "order", created_at, updated_at)
		VALUES (?, ?, COALESCE((SELECT MAX("order") FROM skill_son WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *SkillSonRepo) Update(ctx context.Context, data *entity.SkillSon) error {
	query := `UPDATE skill_son SET name=?, name_eng=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, data.Name, data.NameEng, nowUTC(), data.ID)
	return err
}

func (r *SkillSonRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE skill_son SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
