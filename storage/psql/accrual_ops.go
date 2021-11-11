package psql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

// NewAccrual implements Storage interface.
func (p Psql) NewAccrual(ctx context.Context, orderID model.OrderID, amount float32) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer tx.Rollback()

	// Get user ID from orders table
	row := tx.QueryRowContext(ctx, `SELECT user_id, status FROM orders WHERE id=$1;`, orderID)
	var userID uuid.UUID
	var status model.Status
	if err := row.Scan(&userID, &status); err != nil {
		return storage.ErrNotFound
	}

	// Insert accrual information into accruals_log table. If the data has already been inserted
	// an unique violation error will be threw.
	if _, err = tx.ExecContext(ctx, `INSERT INTO accruals_log (order_id, user_id, sum) VALUES ($1, $2, $3);`,
		orderID, userID, amount); err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return storage.ErrAlreadyProcessed
		}

		return err
	}

	// If the order has not been processed yet, but has been marked as processed, that's an internal server error.
	if status == model.StatusProcessed {
		return storage.ErrInvalidStatus
	}

	// Update status and information about accrual points in order entry.
	if _, err = tx.ExecContext(ctx, `UPDATE orders SET status='PROCESSED', accrual_points=$1 WHERE id=$2;`,
		amount, orderID); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateBalance implements Storage interface.
func (p Psql) UpdateBalance(ctx context.Context) (int, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin: %v", err)
		return 0, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	// Collect information about unprocessed accruals.
	rows, err := tx.QueryContext(ctx, `SELECT order_id, user_id, sum FROM accruals_log WHERE NOT processed;`)
	if err != nil {
		// If there are no unprocessed entries, all is OK.
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		log.Printf("select: %v", err)

		return 0, err
	}
	defer rows.Close()

	type accrual struct {
		orderID string
		userID  uuid.UUID
		sum     float32
	}
	accruals := make([]accrual, 0)
	for rows.Next() {
		var a accrual
		if err := rows.Scan(&a.orderID, &a.userID, &a.sum); err != nil {
			log.Printf("select: scan: %v", err)

			return 0, err
		}
		accruals = append(accruals, a)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	// Add points to each account.
	stmtBalance, err := tx.PrepareContext(ctx, `UPDATE users SET gpoints_balance = gpoints_balance + $1 WHERE id = $2;`)
	if err != nil {
		log.Printf("prepare: users: %v", err)

		return 0, err
	}
	defer stmtBalance.Close()

	// Set flag 'processed' in each entry.
	stmtStatus, err := tx.PrepareContext(ctx, `UPDATE accruals_log SET processed = TRUE WHERE order_id = $1;`)
	if err != nil {
		log.Printf("prepare: accruals: %v", err)

		return 0, err
	}
	defer stmtStatus.Close()

	for _, a := range accruals {
		if _, err := stmtBalance.ExecContext(ctx, a.sum, a.userID); err != nil {
			log.Printf("process: users: %v", err)

			return 0, err
		}
		if _, err := stmtStatus.ExecContext(ctx, a.orderID); err != nil {
			log.Printf("process: accruals: %v", err)

			return 0, err
		}
	}

	return len(accruals), tx.Commit()
}
