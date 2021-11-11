package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

// NewOrder implements Storage interface.
func (p Psql) NewOrder(ctx context.Context, order *model.Order) error {
	const query = `INSERT INTO orders (id, user_id, status, accrual_points, uploaded_at)
	VALUES ($1, $2, $3, $4, $5);`
	_, err := p.db.ExecContext(ctx, query,
		order.ID, order.UserID, order.Status, order.AccrualPoints, order.UploadedAt)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return storage.ErrAlreadyProcessed
		}

		return err
	}

	return nil
}

// OrderByID implements Storage interface.
func (p Psql) OrderByID(ctx context.Context, orderID model.OrderID) (*model.Order, error) {
	row := p.db.QueryRowContext(ctx, `SELECT id, user_id, status, accrual_points, uploaded_at
	FROM orders WHERE id=$1;`, orderID)
	var o model.Order
	if err := row.Scan(&o.ID,
		&o.UserID,
		&o.Status,
		&o.AccrualPoints,
		&o.UploadedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, err
	}

	return &o, nil
}

// OrderByStatus implements Storage interface.
func (p Psql) OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error) {
	orders := make([]model.Order, 0)
	rows, err := p.db.QueryContext(ctx, `SELECT id, user_id, status, accrual_points, uploaded_at
		FROM orders WHERE status = $1;`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.AccrualPoints, &o.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// UpdateOrderStatus implements Storage interface.
func (p Psql) UpdateOrderStatus(ctx context.Context, orderID model.OrderID, status model.Status) error {
	// check if the order exists in db
	row := p.db.QueryRowContext(ctx, `SELECT status	FROM orders WHERE id=$1;`, orderID)
	s := ""
	err := row.Scan(&s)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}

		return err
	}
	// update the status
	if _, err := p.db.ExecContext(ctx, `UPDATE orders SET status=$1 WHERE id=$2;`, status, orderID); err != nil {
		return err
	}

	return nil
}

// UserOrders implements Storage interface.
func (p Psql) UserOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT id, user_id, status, accrual_points, uploaded_at
	FROM orders WHERE user_id=$1 ORDER BY uploaded_at ASC;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID,
			&o.UserID,
			&o.Status,
			&o.AccrualPoints,
			&o.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
