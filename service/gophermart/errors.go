package gophermart

import "errors"

var (
	// ErrWrongPassword is returned by Authenticate when provided password doesn't mismatch with user's password from the storage.
	ErrWrongPassword = errors.New("service: wrong password")
	// ErrNotAuthenticated is returned when authenticated user data is not found in the provided context.
	ErrNotAuthenticated = errors.New("service: no authenticted user found in the context")

	// ErrWrongOrderNumber is returned when the provided order number wasn't found in the storage or order with that number was excuted
	// by another user.
	ErrWrongOrderNumber = errors.New("service: wrong order number")
	// ErrInvalidOrderNumber is threw when the provided order number is empty or hasn't passed Luhn's check.
	ErrInvalidOrderNumber = errors.New("service: invalid order number")

	// ErrOrderExecutedBySameUser is threw when an order with the provided number already executed by the same user.
	ErrOrderExecutedBySameUser = errors.New("service: order already executed by the same user")
	// ErrOrderExecutedByAnotherUser is threw when an order with the provided number already executed by another user.
	ErrOrderExecutedByAnotherUser = errors.New("service: order already executed by another user")
)
