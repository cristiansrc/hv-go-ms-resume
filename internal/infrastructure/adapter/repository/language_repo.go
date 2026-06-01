package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.LanguageRepository = (*LanguageRepo)(nil)

type LanguageRepo struct{ db *DB }

func NewLanguageRepo(db *DB) *LanguageRepo { return &LanguageRepo{db: db} }

func (r *LanguageRepo) List(ctx context.Context) ([]entity.Language, error) {
	query := `SELECT id, language, language_eng, reading_level, writing_level, speaking_level, "order",
		created_at, updated_at, deleted_at FROM language WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Language
	for rows.Next() {
		var e entity.Language
		if err := rows.Scan(&e.ID, &e.Language, &e.LanguageEng, &e.ReadingLevel, &e.WritingLevel,
			&e.SpeakingLevel, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *LanguageRepo) GetByID(ctx context.Context, id int64) (*entity.Language, error) {
	query := `SELECT id, language, language_eng, reading_level, writing_level, speaking_level, "order",
		created_at, updated_at, deleted_at FROM language WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Language
	err := row.Scan(&e.ID, &e.Language, &e.LanguageEng, &e.ReadingLevel, &e.WritingLevel,
		&e.SpeakingLevel, &e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("language not found")
	}
	return &e, err
}

func (r *LanguageRepo) Create(ctx context.Context, data *entity.Language) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO language (language, language_eng, reading_level, writing_level, speaking_level, "order", created_at, updated_at)
		VALUES (?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM language WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Language, data.LanguageEng, data.ReadingLevel, data.WritingLevel, data.SpeakingLevel,
		now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *LanguageRepo) Update(ctx context.Context, data *entity.Language) error {
	query := `UPDATE language SET language=?, language_eng=?, reading_level=?, writing_level=?, speaking_level=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Language, data.LanguageEng, data.ReadingLevel, data.WritingLevel, data.SpeakingLevel,
		nowUTC(), data.ID)
	return err
}

func (r *LanguageRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE language SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
