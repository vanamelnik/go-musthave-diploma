package psql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/vanamelnik/gophermart/model"
	"github.com/vanamelnik/gophermart/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

// ProcessWithdraw implements Storage interface.
func (p Psql) ProcessWithdraw(ctx context.Context, withdraw *model.Withdrawal) error {
	// Try to create a new entry in the withdrawals_log table. If the order has already been processed, return an error.
	withdraw.Status = model.StatusProcessing
	if withdraw.Sum < 0 {
		withdraw.Status = model.StatusInvalid
	}

	if _, err := p.db.ExecContext(ctx, `INSERT INTO withdrawals_log (order_id, user_id, sum, status, processed_at)
	VALUES ($1, $2, $3, $4, $5);`, withdraw.OrderID, withdraw.UserID, withdraw.Sum, withdraw.Status, withdraw.ProcessedAt); err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return storage.ErrAlreadyProcessed
		}

		return err
	}
	if withdraw.Status == model.StatusInvalid {
		return storage.ErrInvalidInput
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// if the transaction has been rolled back, set status to 'INVALID'
		err := tx.Rollback()
		if !errors.Is(err, sql.ErrTxDone) { // if transaction was rejected, set status to INVALID
			log.Printf("transaction has been rolled back. Updating withdrawals_log set status=INVALID orderID=%v", withdraw.OrderID)
			//nolint:errcheck
			p.db.ExecContext(ctx, `UPDATE withdrawals_log SET status='INVALID' WHERE order_id=$1;`, withdraw.OrderID)
		}
	}()

	// Try to change user's balance. If there are no enough points, return an error.
	row := tx.QueryRowContext(ctx, `SELECT gpoints_balance FROM users WHERE id=$1 FOR UPDATE;`, withdraw.UserID)
	var balance float32
	if err := row.Scan(&balance); err != nil {
		return err
	}
	if balance < withdraw.Sum {
		return storage.ErrInsufficientPoints
	}
	if _, err := tx.ExecContext(ctx, `UPDATE users SET gpoints_balance = gpoints_balance - $1 WHERE id=$2;`,
		withdraw.Sum, withdraw.UserID); err != nil {
		return err
	}

	// All is OK, set status to 'processed'.
	if _, err := tx.ExecContext(ctx, `UPDATE withdrawals_log SET status='PROCESSED' WHERE order_id=$1;`, withdraw.OrderID); err != nil {
		return err
	}

	return tx.Commit()
}

// WithdrawalsByUserId implements Storage interface.
func (p Psql) WithdrawalsByUserID(ctx context.Context, id uuid.UUID) ([]model.Withdrawal, error) {
	withdrawals := make([]model.Withdrawal, 0)
	rows, err := p.db.QueryContext(ctx, `SELECT order_id, user_id, sum, status, processed_at
		FROM withdrawals_log WHERE user_id = $1 ORDER BY processed_at ASC;`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var w model.Withdrawal
		if err := rows.Scan(&w.OrderID, &w.UserID, &w.Sum, &w.Status, &w.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
