package gophermart

import (
	"context"

	"github.com/vanamelnik/go-musthave-diploma/model"
)

type (
	// Service defines model.Service operations.
	Service interface {
		// Create creates a new user, hashes his password and stores it in the storage
		Create(ctx context.Context, login, password string) (*model.User, error)
		// Authentcate checks whether a user with such login and password is in the storage.
		// If successful, the model.User object is saved in the ctx.
		Authenticate(ctx context.Context, login, password string) (*model.User, error)

		// The data of authenticated user is taken from the context.

		// GetOrders fetches all orders of authenticated user from the storage.
		GetOrders(ctx context.Context) ([]model.Order, error)
		// GetBalance returns information about authenticated user's bonus balance and total withdrawn amount.
		GetBalance(ctx context.Context) (*UserBalance, error)
		// GetWithdrawals returns information about all withdrawal transactions of authenticated user.
		GetWithdrawals(ctx context.Context) ([]model.Withdrawal, error)

		// ProcessOrder adds the order provided to the storage (marked as 'NEW')
		ProcessOrder(ctx context.Context, orderID model.OrderID) error
		// Withdraw adds a new entry to withdrawals log and subtracts the sum from authenticated user's bonus balance.
		Withdraw(ctx context.Context, orderID model.OrderID, sum float32) error

		// Close shuts down the service.
		Close()
	}

	// UserBalance is a struct returned by GetBalance.
	UserBalance struct {
		// Current is current amount of GopherPoints (1 GPoint = 1 RUB)
		Current float32 `json:"current"`
		// Withdrawn is total withdrawn amount.
		Withdrawn float32 `json:"withdrawn"`
	}
)
