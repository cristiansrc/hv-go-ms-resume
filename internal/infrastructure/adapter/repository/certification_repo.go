package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cristiansrc/hv-go-ms-resume/internal/application/port/output"
	"github.com/cristiansrc/hv-go-ms-resume/internal/domain/entity"
)

var _ output.CertificationRepository = (*CertificationRepo)(nil)

type CertificationRepo struct{ db *DB }

func NewCertificationRepo(db *DB) *CertificationRepo { return &CertificationRepo{db: db} }

func (r *CertificationRepo) List(ctx context.Context) ([]entity.Certification, error) {
	query := `SELECT id, name, name_eng, issuing_organization, issuing_organization_eng,
		issue_date, expiration_date, verification_url, credential_id,
		description, description_eng, summary_pdf, summary_pdf_eng, "order",
		created_at, updated_at, deleted_at FROM certification WHERE deleted_at IS NULL ORDER BY "order" ASC`
	rows, err := r.db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entity.Certification
	for rows.Next() {
		var e entity.Certification
		if err := rows.Scan(&e.ID, &e.Name, &e.NameEng, &e.IssuingOrganization, &e.IssuingOrganizationEng,
			&e.IssueDate, &e.ExpirationDate, &e.VerificationURL, &e.CredentialID,
			&e.Description, &e.DescriptionEng, &e.SummaryPdf, &e.SummaryPdfEng,
			&e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *CertificationRepo) GetByID(ctx context.Context, id int64) (*entity.Certification, error) {
	query := `SELECT id, name, name_eng, issuing_organization, issuing_organization_eng,
		issue_date, expiration_date, verification_url, credential_id,
		description, description_eng, summary_pdf, summary_pdf_eng, "order",
		created_at, updated_at, deleted_at FROM certification WHERE id = ? AND deleted_at IS NULL`
	row := r.db.conn.QueryRowContext(ctx, query, id)
	var e entity.Certification
	err := row.Scan(&e.ID, &e.Name, &e.NameEng, &e.IssuingOrganization, &e.IssuingOrganizationEng,
		&e.IssueDate, &e.ExpirationDate, &e.VerificationURL, &e.CredentialID,
		&e.Description, &e.DescriptionEng, &e.SummaryPdf, &e.SummaryPdfEng,
		&e.Order, &e.CreatedAt, &e.UpdatedAt, &e.DeletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("certification not found")
	}
	return &e, err
}

func (r *CertificationRepo) Create(ctx context.Context, data *entity.Certification) (int64, error) {
	now := nowUTC()
	query := `INSERT INTO certification (name, name_eng, issuing_organization, issuing_organization_eng,
		issue_date, expiration_date, verification_url, credential_id,
		description, description_eng, summary_pdf, summary_pdf_eng, "order", created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		COALESCE((SELECT MAX("order") FROM certification WHERE deleted_at IS NULL), 0) + 1, ?, ?)`
	result, err := r.db.conn.ExecContext(ctx, query,
		data.Name, data.NameEng, data.IssuingOrganization, data.IssuingOrganizationEng,
		data.IssueDate, data.ExpirationDate, data.VerificationURL, data.CredentialID,
		data.Description, data.DescriptionEng, data.SummaryPdf, data.SummaryPdfEng,
		now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CertificationRepo) Update(ctx context.Context, data *entity.Certification) error {
	query := `UPDATE certification SET name=?, name_eng=?, issuing_organization=?, issuing_organization_eng=?,
		issue_date=?, expiration_date=?, verification_url=?, credential_id=?,
		description=?, description_eng=?, summary_pdf=?, summary_pdf_eng=?,
		updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query,
		data.Name, data.NameEng, data.IssuingOrganization, data.IssuingOrganizationEng,
		data.IssueDate, data.ExpirationDate, data.VerificationURL, data.CredentialID,
		data.Description, data.DescriptionEng, data.SummaryPdf, data.SummaryPdfEng,
		nowUTC(), data.ID)
	return err
}

func (r *CertificationRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE certification SET deleted_at=?, updated_at=? WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.conn.ExecContext(ctx, query, nowUTC(), nowUTC(), id)
	return err
}
