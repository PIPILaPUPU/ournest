package wishlist

import (
	"context"
	"errors"
	"fmt"

	"wishlistapp/internal/notifications"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("wish not found")

// WishRepository отвечает только за работу с БД — SQL, Scan, ошибки pgx.
// Handler не знает про запросы; repository не знает про HTTP.
type WishRepository struct {
	db       *pgxpool.Pool
	notifier *notifications.Repository
}

func NewWishRepository(db *pgxpool.Pool, notifier *notifications.Repository) *WishRepository {
	return &WishRepository{db: db, notifier: notifier}
}

func scanWish(row pgx.Row) (Wish, error) {
	var wish Wish
	err := row.Scan(
		&wish.ID,
		&wish.OwnerID,
		&wish.OwnerUsername,
		&wish.Title,
		&wish.Description,
		&wish.URL,
		&wish.Price,
		&wish.Status,
		&wish.GroupName,
		&wish.GroupColor,
		&wish.CreatedAt,
		&wish.UpdatedAt,
	)
	return wish, err
}

const wishSelectColumns = `
	w.id,
	w.owner_id,
	u.username,
	w.title,
	w.description,
	w.url,
	w.price,
	w.status,
	w.group_name,
	w.group_color,
	w.created_at,
	w.updated_at
`

func (r *WishRepository) List(ctx context.Context) ([]Wish, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+wishSelectColumns+`
		FROM wishes w
		JOIN users u ON u.id = w.owner_id
		ORDER BY w.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query wishes: %w", err)
	}
	defer rows.Close()

	wishes := make([]Wish, 0)
	for rows.Next() {
		var wish Wish
		if err := rows.Scan(
			&wish.ID,
			&wish.OwnerID,
			&wish.OwnerUsername,
			&wish.Title,
			&wish.Description,
			&wish.URL,
			&wish.Price,
			&wish.Status,
			&wish.GroupName,
			&wish.GroupColor,
			&wish.CreatedAt,
			&wish.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan wish: %w", err)
		}
		wishes = append(wishes, wish)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wishes: %w", err)
	}

	return wishes, nil
}

func (r *WishRepository) Create(ctx context.Context, input CreateWishInput) (Wish, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Wish{}, fmt.Errorf("begin create wish: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	row := tx.QueryRow(ctx, `
		WITH inserted AS (
			INSERT INTO wishes (owner_id, title, description, url, price, status, group_name, group_color)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING *
		)
		SELECT
			i.id,
			i.owner_id,
			u.username,
			i.title,
			i.description,
			i.url,
			i.price,
			i.status,
			i.group_name,
			i.group_color,
			i.created_at,
			i.updated_at
		FROM inserted i
		JOIN users u ON u.id = i.owner_id`,
		input.OwnerID,
		input.Title,
		input.Description,
		input.URL,
		input.Price,
		input.Status,
		input.GroupName,
		input.GroupColor,
	)

	wish, err := scanWish(row)
	if err != nil {
		return Wish{}, fmt.Errorf("insert wish: %w", err)
	}

	if err := r.notifier.EnqueueForPartner(
		ctx,
		tx,
		notifications.EventWishCreated,
		wish.OwnerID,
		notifications.WishCreatedPayload(wish.OwnerUsername, wish.Title),
	); err != nil {
		return Wish{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Wish{}, fmt.Errorf("commit create wish: %w", err)
	}
	return wish, nil
}

func (r *WishRepository) Update(ctx context.Context, id int, input UpdateWishInput) (Wish, error) {
	row := r.db.QueryRow(ctx, `
		WITH updated AS (
			UPDATE wishes SET
				description = COALESCE($2, description),
				url         = COALESCE($3, url),
				price       = COALESCE($4, price),
				updated_at  = NOW()
			WHERE id = $1
			RETURNING *
		)
		SELECT
			up.id,
			up.owner_id,
			u.username,
			up.title,
			up.description,
			up.url,
			up.price,
			up.status,
			up.group_name,
			up.group_color,
			up.created_at,
			up.updated_at
		FROM updated up
		JOIN users u ON u.id = up.owner_id`,
		id,
		input.Description,
		input.URL,
		input.Price,
	)

	wish, err := scanWish(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Wish{}, ErrNotFound
		}
		return Wish{}, fmt.Errorf("update wish: %w", err)
	}

	return wish, nil
}

func (r *WishRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM wishes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete wish: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
