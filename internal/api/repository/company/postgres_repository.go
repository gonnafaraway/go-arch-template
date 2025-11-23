package company

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-arch-template/internal/api/domain/company"
	postgres "go-arch-template/internal/api/storage/postgres"
)

type PostgresRepository struct {
	db *postgres.Client
}

func NewPostgresRepository(db *postgres.Client) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, c *company.Company) error {
	query := `
		INSERT INTO companies (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	if c.ID == "" {
		c.ID = fmt.Sprintf("company_%d", time.Now().UnixNano())
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.DB().ExecContext(ctx, query, c.ID, c.Name, c.Email, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*company.Company, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM companies
		WHERE id = $1
	`

	var c company.Company
	err := r.db.DB().QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]*company.Company, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM companies
		ORDER BY created_at DESC
	`

	rows, err := r.db.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*company.Company
	for rows.Next() {
		var c company.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, &c)
	}

	return companies, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, c *company.Company) error {
	query := `
		UPDATE companies
		SET name = $1, email = $2, updated_at = $3
		WHERE id = $4
	`

	c.UpdatedAt = time.Now()
	result, err := r.db.DB().ExecContext(ctx, query, c.Name, c.Email, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM companies WHERE id = $1`

	result, err := r.db.DB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

