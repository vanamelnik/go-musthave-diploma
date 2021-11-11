package model

import (
	"time"

	"github.com/vanamelnik/go-musthave-diploma/pkg/luhn"

	"github.com/google/uuid"
)

/* Номером заказа является последовательность цифр произвольной длины.
Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.*/

const (
	StatusNew        Status = "NEW"        // the order was uploaded, but was not processed
	StatusRegistered Status = "REGISTERED" // the order is registered, but the accural is not claculated
	StatusInvalid    Status = "INVALID"    // the order is not accepted and the accural is not calculated
	StatusProcessing Status = "PROCESSING" // reward for the order is being calculated and reward
	StatusProcessed  Status = "PROCESSED"  // calculating of accural is complete.
)

type (
	// Order represents information of the order received by UserService.
	Order struct {
		ID            OrderID   `json:"number"`
		UserID        uuid.UUID `json:"-"`
		Status        Status    `json:"status"`
		AccrualPoints float32   `json:"accrual,omitempty"`
		UploadedAt    time.Time `json:"uploaded_at"`
	}

	// OrderID is a sequence of numbers of arbitrary length.
	// The number must satisfy Luhn algorithm.
	OrderID string

	// Status represents the status of the order.
	Status string
)

// Valid validades the order status string.
func (s Status) Valid() bool {
	return s == StatusNew ||
		s == StatusRegistered ||
		s == StatusInvalid ||
		s == StatusProcessing ||
		s == StatusProcessed
}

// Valid validates the order ID.
func (id OrderID) Valid() bool {
	if len(id) < 2 {
		return false
	}

	return luhn.Validate(string(id))
}

// String implements Stringer interface.
func (id OrderID) String() string {
	return string(id)
}
