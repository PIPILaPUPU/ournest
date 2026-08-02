package dateideas

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("date idea not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

const selectColumns = `
	d.id,
	d.author_id,
	u.username,
	d.title,
	d.description,
	d.status,
	d.created_at,
	d.updated_at
`

func scanIdea(row pgx.Row) (DateIdea, error) {
	var idea DateIdea
	err := row.Scan(
		&idea.ID,
		&idea.AuthorID,
		&idea.AuthorUsername,
		&idea.Title,
		&idea.Description,
		&idea.Status,
		&idea.CreatedAt,
		&idea.UpdatedAt,
	)
	return idea, err
}

func (r *Repository) List(ctx context.Context) ([]DateIdea, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+selectColumns+`
		FROM date_ideas d
		JOIN users u ON u.id = d.author_id
		ORDER BY d.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query date ideas: %w", err)
	}
	defer rows.Close()

	ideas := make([]DateIdea, 0)
	for rows.Next() {
		idea, err := scanIdea(rows)
		if err != nil {
			return nil, fmt.Errorf("scan date idea: %w", err)
		}
		ideas = append(ideas, idea)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate date ideas: %w", err)
	}
	return ideas, nil
}

func (r *Repository) Create(ctx context.Context, authorID int, title string, description *string) (DateIdea, error) {
	row := r.db.QueryRow(ctx, `
		WITH inserted AS (
			INSERT INTO date_ideas (author_id, title, description)
			VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT
			i.id, i.author_id, u.username, i.title, i.description,
			i.status, i.created_at, i.updated_at
		FROM inserted i
		JOIN users u ON u.id = i.author_id`,
		authorID,
		title,
		description,
	)

	idea, err := scanIdea(row)
	if err != nil {
		return DateIdea{}, fmt.Errorf("insert date idea: %w", err)
	}
	return idea, nil
}

func (r *Repository) Update(ctx context.Context, id int, input UpdateRequest) (DateIdea, error) {
	row := r.db.QueryRow(ctx, `
		WITH updated AS (
			UPDATE date_ideas SET
				description = COALESCE($2, description),
				status = COALESCE($3::date_idea_status, status),
				updated_at = NOW()
			WHERE id = $1
			RETURNING *
		)
		SELECT
			d.id, d.author_id, u.username, d.title, d.description,
			d.status, d.created_at, d.updated_at
		FROM updated d
		JOIN users u ON u.id = d.author_id`,
		id,
		input.Description,
		input.Status,
	)

	idea, err := scanIdea(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DateIdea{}, ErrNotFound
		}
		return DateIdea{}, fmt.Errorf("update date idea: %w", err)
	}
	return idea, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM date_ideas WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete date idea: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
