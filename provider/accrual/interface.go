//go:generate mockgen -source=interface.go -destination=./mock/mock_accrual.go -package=mockaccrual
package accrual

import (
	"context"
	"fmt"

	"github.com/vanamelnik/go-musthave-diploma/model"
)

// AccrualClient provides client requests to GopherAccuralService
type AccrualClient interface {
	Request(ctx context.Context, orderID model.OrderID) (*AccrualResponse, error)
}

// AccrualResponse represents the response of GopherAccrual service.
type AccrualResponse struct {
	Order   model.OrderID `json:"order"`
	Status  model.Status  `json:"status"`
	Accrual float32       `json:"accrual,omitempty"`
}

// ErrUnexpectedStatus is returned when the server returns an unexpected status code.
type ErrUnexpectedStatus struct {
	Code int
	Body string
}

func (e ErrUnexpectedStatus) Error() string {
	return fmt.Sprintf("accrual request: unexpected status code %d, body: %s", e.Code, e.Body)
}
