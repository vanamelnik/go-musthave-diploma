package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

//NewUser implements Storage interface.
func (p Psql) NewUser(ctx context.Context, user *model.User) error {
	const query = `INSERT INTO users (id, login, password_hash, gpoints_balance, created_at)
	VALUES ($1, $2, $3, $4, $5);`
	_, err := p.db.ExecContext(ctx, query,
		user.ID, user.Login, user.PasswordHash, user.GPointsBalance, user.CreatedAt)
	if err != nil {
		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return storage.ErrAlreadyProcessed
		}

		return err
	}

	return nil
}

// UserByLogin implements Storage interface.
func (p Psql) UserByLogin(ctx context.Context, login string) (*model.User, error) {
	row := p.db.QueryRowContext(ctx, `SELECT id, password_hash, gpoints_balance, remember_token, created_at
	FROM users WHERE login=$1;`, login)
	u := &model.User{Login: login}
	err := row.Scan(&u.ID, &u.PasswordHash, &u.GPointsBalance, &u.RememberToken, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

// UserByRemember implements Storage interface.
func (p Psql) UserByRemember(ctx context.Context, remember string) (*model.User, error) {
	row := p.db.QueryRowContext(ctx, `SELECT id, password_hash, gpoints_balance, login, created_at
	FROM users WHERE remember_token=$1;`, remember)
	u := &model.User{RememberToken: remember}
	err := row.Scan(&u.ID, &u.PasswordHash, &u.GPointsBalance, &u.Login, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, err
	}

	return u, nil
}

// UpdateUser implements Storage interface.
func (p Psql) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := p.db.ExecContext(ctx, `UPDATE users SET
	login=$1, password_hash=$2, remember_token=$3
	WHERE id=$4;`, user.Login, user.PasswordHash, user.RememberToken, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrNotFound
		}

		return err
	}

	return nil
}
