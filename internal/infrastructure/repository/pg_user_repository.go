package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"guilhermefaleiros/go-service-template/internal/domain/entity"
)

type PGUserRepository struct {
	conn *pgxpool.Pool
}

func (r *PGUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := "SELECT id, name, email, status, phone, created_at, updated_at FROM users WHERE email = $1"
	row := r.conn.QueryRow(ctx, query, email)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PGUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, name, email, status, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.conn.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Status,
		user.Phone,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewPGUserRepository(conn *pgxpool.Pool) *PGUserRepository {
	return &PGUserRepository{conn: conn}
}
