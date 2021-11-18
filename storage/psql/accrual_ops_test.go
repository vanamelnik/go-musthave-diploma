package psql

import (
	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

func (ts *TestSuite) TestAccrual() {
	tt := []struct {
		name             string
		orderID          model.OrderID
		wantErr          error
		amount           float32 // the linter wants to save 8 bytes of memory by changing the order of items in this struct))
		checkIfProcessed bool
	}{
		{
			name:             "#1 Alice got 500$",
			orderID:          "117",
			amount:           500,
			wantErr:          nil,
			checkIfProcessed: true,
		},
		{
			name:             "#2 Bob got 1.25$",
			orderID:          "018",
			amount:           1.25,
			wantErr:          nil,
			checkIfProcessed: true,
		},
		{
			name:             "#3 try to repeat the same order",
			orderID:          "117",
			amount:           500,
			wantErr:          storage.ErrAlreadyProcessed,
			checkIfProcessed: false,
		},
		{
			name:             "#4 wrong order id",
			orderID:          "666",
			amount:           1000000000,
			wantErr:          storage.ErrNotFound,
			checkIfProcessed: false,
		},
		{
			name:             "#5 internal error - order already marked as PROCESSED",
			orderID:          "125",
			amount:           0.01,
			wantErr:          storage.ErrInvalidStatus,
			checkIfProcessed: false,
		},
		{
			name:             "#6 Bob got more 0.75$",
			orderID:          "026",
			amount:           0.75,
			wantErr:          nil,
			checkIfProcessed: true,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			err := ts.storage.CreateAccrual(ts.ctx, tc.orderID, tc.amount)
			ts.Assert().ErrorIs(err, tc.wantErr)
			if tc.checkIfProcessed {
				o, err := ts.storage.OrderByID(ts.ctx, tc.orderID)
				ts.Require().NoError(err)
				ts.Assert().Equal(model.StatusProcessed, o.Status)
			}
		})
	}
	ts.Run("#one_more UpdateBalance()", func() {
		n, err := ts.storage.UpdateBalance(ts.ctx)
		ts.Assert().NoError(err)
		ts.Assert().EqualValues(3, n)

		// Check users' balances
		alice, err := ts.storage.UserByLogin(ts.ctx, ts.alice.user.Login)
		ts.Require().NoError(err)
		bob, err := ts.storage.UserByLogin(ts.ctx, ts.bob.user.Login)
		ts.Require().NoError(err)
		ts.Assert().EqualValues(alice.GPointsBalance, 500.0)
		ts.Assert().EqualValues(bob.GPointsBalance, 1.25+0.75) // Bob got two accruals

		// Check orders AccrualPoints fields
		aliceOrder, err := ts.storage.OrderByID(ts.ctx, "117")
		ts.Require().NoError(err)
		bobOrder, err := ts.storage.OrderByID(ts.ctx, "018")
		ts.Require().NoError(err)
		ts.Assert().EqualValues(aliceOrder.AccrualPoints, 500)
		ts.Assert().EqualValues(bobOrder.AccrualPoints, 1.25)
	})
}
