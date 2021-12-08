package psql

import (
	"time"

	"github.com/vanamelnik/gophermart/model"
	"github.com/vanamelnik/gophermart/storage"

	"github.com/google/uuid"
)

func (ts *TestSuite) TestNewOrder() {
	tt := []struct {
		name    string
		orderID model.OrderID
		userID  uuid.UUID
		status  model.Status

		wantErr bool
	}{
		{
			name:    "#1 One more order for Bob",
			orderID: "034",
			userID:  ts.bob.user.ID,
			status:  "INVALID",
			wantErr: false,
		},
		{
			name:    "#2 An order for non-existing user", // we don't need to use specified error
			orderID: "111111",                            // order number is not validated at storage level...
			userID:  uuid.New(),
			status:  "NEW",
			wantErr: true,
		},
		{
			name:    "#3 Try to process the same order",
			orderID: "034",
			userID:  ts.bob.user.ID,
			status:  "PROCESSING",
			wantErr: true,
		},
		{
			name:    "#4 Wrong status code",
			orderID: "133",
			userID:  ts.alice.user.ID,
			status:  "LONELY",
			wantErr: true,
		},
		{
			name:    "#4 Empty status code",
			orderID: "141",
			userID:  ts.alice.user.ID,
			status:  "",
			wantErr: true,
		},
	}

	for _, tc := range tt {
		ts.Run(tc.name, func() {
			err := ts.storage.CreateOrder(ts.ctx, &model.Order{
				ID:         tc.orderID,
				UserID:     tc.userID,
				Status:     tc.status,
				UploadedAt: time.Now(),
			})
			switch tc.wantErr {
			case true:
				ts.Assert().Error(err)
				ts.T().Logf("Returned error: %s", err)
			case false:
				ts.Assert().NoError(err)
			}
		})
	}
}

func (ts *TestSuite) TestOrderByID() {
	tt := []struct {
		name    string
		orderID model.OrderID
		wantErr error
	}{
		{
			name:    "#1 Fetch an existing order",
			orderID: "125",
			wantErr: nil,
		},
		{
			name:    "#2 Fetch an non-existing order",
			orderID: "414",
			wantErr: storage.ErrNotFound,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			_, err := ts.storage.OrderByID(ts.ctx, tc.orderID)
			ts.Assert().ErrorIs(err, tc.wantErr)
		})
	}
}

func (ts *TestSuite) TestUserOrders() {
	tt := []struct {
		name                  string
		id                    uuid.UUID
		lengthGreaterThanZero bool
		wantErr               error
	}{
		{
			name:                  "#1 Bob's orders",
			id:                    ts.bob.user.ID,
			lengthGreaterThanZero: true,
			wantErr:               nil,
		},
		{
			name:                  "#2 Alice's orders",
			id:                    ts.alice.user.ID,
			lengthGreaterThanZero: true,
			wantErr:               nil,
		},
		{
			name:                  "#2 non-existing user orders",
			id:                    uuid.New(),
			lengthGreaterThanZero: false,
			wantErr:               nil,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			orders, err := ts.storage.UserOrders(ts.ctx, tc.id)
			ts.T().Logf("Fetched %d orders", len(orders))
			ts.Assert().Condition(func() bool { return (len(orders) > 0) == tc.lengthGreaterThanZero })
			ts.Assert().ErrorIs(err, tc.wantErr)
		})
	}
}

func (ts *TestSuite) TestOrdersByStatus() {
	tt := []struct {
		name    string
		status  model.Status
		wantErr error
	}{
		{
			name:    "#1 PROCESSING",
			status:  "PROCESSING",
			wantErr: nil,
		},
		{
			name:    "#2 no such status loaded", // should return an empty list and no error
			status:  "NEW",
			wantErr: nil,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			orders, err := ts.storage.OrdersByStatus(ts.ctx, tc.status)
			ts.Assert().ErrorIs(err, tc.wantErr)
			ts.T().Logf("Fetched %d orders", len(orders))
		})
	}
}

func (ts *TestSuite) TestUpdateStatus() {
	const orderID model.OrderID = "067"
	ts.Require().NoError(ts.storage.CreateOrder(ts.ctx, &model.Order{
		ID:            orderID,
		UserID:        ts.bob.user.ID,
		Status:        "REGISTERED",
		AccrualPoints: 0,
		UploadedAt:    time.Now(),
	}))
	tt := []struct {
		name    string
		id      model.OrderID
		wantErr error
	}{
		{
			name:    "#1 Change status of existing order",
			id:      orderID,
			wantErr: nil,
		},
		{
			name:    "#2 Change status of non-existing order",
			id:      "1111",
			wantErr: storage.ErrNotFound,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			err := ts.storage.UpdateOrderStatus(ts.ctx, tc.id, model.StatusProcessed)
			ts.Assert().ErrorIs(err, tc.wantErr)
			o, err := ts.storage.OrderByID(ts.ctx, orderID)
			ts.Assert().NoError(err)
			if tc.wantErr == nil {
				ts.Assert().Equal(model.StatusProcessed, o.Status)
			}
		})
	}

}
