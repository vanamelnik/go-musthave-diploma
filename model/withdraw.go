package model

import (
	"time"

	"github.com/google/uuid"
)

// Withdrawal represents information about a withdrawal transaction from user's bonus account.
// When a withdraw request is received, the service adds a new entry with the "PROCESSING" status to withdrawals log in the storage.
// Then the transaction begins and the service attempts to withdraw the amount provided from user's balance.
// If successful, the withdraw status is set to "PROCESSED", otherwise the transaction is rejected and status is set to "INVALID".
type Withdrawal struct {
	UserID  uuid.UUID
	OrderID OrderID `json:"order"`
	// Sum in G-Points
	Sum float32 `json:"sum"`

	Status Status `json:"status"`

	ProcessedAt time.Time `json:"processed_at"`
}
