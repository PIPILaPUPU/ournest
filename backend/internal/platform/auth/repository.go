package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindCredentialsByUsername(ctx context.Context, username string) (UserCredentials, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		username,
	)

	var creds UserCredentials
	if err := row.Scan(&creds.ID, &creds.Username, &creds.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserCredentials{}, ErrInvalidCredentials
		}
		return UserCredentials{}, fmt.Errorf("find user by username: %w", err)
	}

	return creds, nil
}
