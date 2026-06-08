package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.ReferenceRepository = (*ReferenceRepo)(nil)

type ReferenceRepo struct{ db *DB }

func NewReferenceRepo(db *DB) *ReferenceRepo { return &ReferenceRepo{db: db} }

func (r *ReferenceRepo) List(ctx context.Context) ([]entity.Reference, error) {
	query := `SELECT id, full_name, position, company, company_eng, email, phone,
		relationship, relationship_eng, "order",
		created_at, updated_at, deleted_at FROM reference WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Reference
	for rows.Next() {
		var e entity.Reference
		if err := rows.Scan(&e.ID, &e.FullName, &e.Position, &e.Company, &e.CompanyEng,
			&e.Email, &e.Phone, &e.Relationship, &e.RelationshipEng,
			&e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *ReferenceRepo) GetByID(ctx context.Context, id int64) (*entity.Reference, error) {
	query := `SELECT id, full_name, position, company, company_eng, email, phone,
		relationship, relationship_eng, "order",
		created_at, updated_at, deleted_at FROM reference WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Reference
	err := row.Scan(&e.ID, &e.FullName, &e.Position, &e.Company, &e.CompanyEng,
		&e.Email, &e.Phone, &e.Relationship, &e.RelationshipEng,
		&e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("reference: %w", entity.ErrNotFound)
	}
	return &e, err
}

func (r *ReferenceRepo) Create(ctx context.Context, data *entity.Reference) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO reference (full_name, position, company, company_eng, email, phone,
		relationship, relationship_eng, "order", created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM reference WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.FullName, data.Position, data.Company, data.CompanyEng,
		data.Email, data.Phone, data.Relationship, data.RelationshipEng,
		now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ReferenceRepo) Update(ctx context.Context, data *entity.Reference) error {
	query := `UPDATE reference SET full_name=?, position=?, company=?, company_eng=?, email=?, phone=?,
		relationship=?, relationship_eng=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.FullName, data.Position, data.Company, data.CompanyEng,
		data.Email, data.Phone, data.Relationship, data.RelationshipEng,
		nowUTC(), data.ID)
	return err
}

func (r *ReferenceRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE reference SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
