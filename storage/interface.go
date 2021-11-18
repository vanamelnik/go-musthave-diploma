//go:generate mockgen -source=interface.go -destination=./mock/mock_storage.go -package=mockstorage
package storage

import (
	"context"
	"errors"

	"github.com/vanamelnik/go-musthave-diploma/model"

	"github.com/google/uuid"
)

type Storage interface {
	// CreateUser adds user information to db. 'Password' field is ignored.
	CreateUser(ctx context.Context, user model.User) error
	// UserByLogin looks for a user with provided login.
	UserByLogin(ctx context.Context, login string) (*model.User, error)
	// UserByRemember lloks for a user with provided remember token.
	UserByRemember(ctx context.Context, remember string) (*model.User, error)
	// UpdateUser updates user information (login, password hash and remember token).
	UpdateUser(ctx context.Context, user model.User) error

	// CreateOrder creates a new entry in the orders table.
	CreateOrder(ctx context.Context, order *model.Order) error
	// UpdateOrderStatus sets the status of the order with orderId provided to the value provided.
	UpdateOrderStatus(ctx context.Context, orderID model.OrderID, status model.Status) error
	// UserOrders gets all orders made by the provided user.
	UserOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	// OrderByID searches for order with the provided order id.
	OrderByID(ctx context.Context, orderID model.OrderID) (*model.Order, error)
	// OrderByStatus returns all orders with specified status. If there aren't any, empty slice is returned.
	OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error)

	// CreateAccrual adds a new entry into the accruals_log table and updates an order status in orders table.
	// orderId must be unique.
	CreateAccrual(ctx context.Context, orderID model.OrderID, amount float32) error
	// UpdateBalance checks all unprocessed accruals in the table and adds the points to users' balances.
	// Flags 'processed' are set to true.
	UpdateBalance(ctx context.Context) (int, error)

	// CreateWithdraw creates a new entry in the withdrawals_log table and updates user's balance.
	// This function must update and check users's balance and return the error if the balance is less than the amount provided.
	// OrderId must be unique.
	CreateWithdraw(ctx context.Context, withdraw *model.Withdrawal) error
	// WithdrawalsByUserID fetches all withdrawals made by the provided user. If there aren't any, empty slice is returned.
	WithdrawalsByUserID(ctx context.Context, id uuid.UUID) ([]model.Withdrawal, error)

	// Close shuts the database down.
	Close() error
}

var (
	// ErrNotFound is returned when there's no data is available in the database.
	ErrNotFound           = errors.New("storage: not found")
	ErrLoginAlreadyExists = errors.New("storage: login (login) already occupied")

	ErrAlreadyProcessed = errors.New("storage: already processed")
	ErrInvalidStatus    = errors.New("storage: non-processed order has PROCESSED status")

	// ErrInsufficientPoints is threw when there are not enough G-Points on user's account to perform withdrawal operation.
	ErrInsufficientPoints = errors.New("storage: insufficient points to perform withdrawal operation")

	// ErrInvalidInput is threw when accrual or withdrawal amount less than zero.
	ErrInvalidInput = errors.New("storage: amount less than zero")
)
