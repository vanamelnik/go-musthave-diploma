package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vanamelnik/go-musthave-diploma/model"
)

var _ AccrualClient = (*Accrual)(nil)

const accrualRequestAPI = "/api/orders/"

// Accrual is implementation of api.AccrualClient interface.
type Accrual struct {
	accrualAPI string
}

// New creates a new instance of Accrual client.
func New(accrualSystemURL string) Accrual {
	return Accrual{
		accrualAPI: accrualSystemURL + accrualRequestAPI,
	}
}

// Request performs a request to GopherAccrualService.
func (a Accrual) Request(ctx context.Context, orderID model.OrderID) (*AccrualResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.accrualAPI+string(orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("client: AccrualRequest: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: AccrualRequest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &ErrUnexpectedStatus{
			Code: resp.StatusCode,
			Body: string(body),
		}
	}

	dec := json.NewDecoder(resp.Body)
	response := AccrualResponse{}
	if err := dec.Decode(&response); err != nil {
		return nil, fmt.Errorf("client: AccrualRequest: %w", err)
	}

	return &response, nil
}
