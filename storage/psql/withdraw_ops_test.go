package psql

import (
	"time"

	"github.com/google/uuid"
	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/storage"
)

func (ts *TestSuite) TestWithdraw() {
	dave := &model.User{
		ID:             uuid.New(),
		Login:          "davidbowie@mail.su",
		PasswordHash:   "qWeRtYuIoP",
		CreatedAt:      time.Now(),
		GPointsBalance: 300,
	}
	ts.Require().NoError(ts.storage.NewUser(ts.ctx, dave))

	tt := []struct {
		name    string
		userID  uuid.UUID
		orderID model.OrderID
		sum     float32
		wantErr error
	}{
		{
			name:    "#1 Dave withdraws 200",
			userID:  dave.ID,
			orderID: "1016",
			sum:     200,
			wantErr: nil,
		},
		{
			name:    "#2 Bob tries to withdraw 200",
			userID:  ts.bob.user.ID,
			orderID: "2022",
			sum:     200,
			wantErr: storage.ErrInsufficientPoints,
		},
		{
			name:    "#3 The same order as #1",
			userID:  dave.ID,
			orderID: "1016",
			sum:     200,
			wantErr: storage.ErrAlreadyProcessed,
		},
		{
			name:    "#4 Dave withdraws another 200",
			userID:  dave.ID,
			orderID: "3037",
			sum:     200,
			wantErr: storage.ErrInsufficientPoints,
		},
		{
			name:    "#5 sum < 0",
			userID:  ts.alice.user.ID,
			orderID: "4044",
			sum:     -100,
			wantErr: storage.ErrInvalidInput,
		},
		{
			name:    "#6 Dave withdraws 99.99",
			userID:  dave.ID,
			orderID: "5058",
			sum:     99.5,
			wantErr: nil,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			err := ts.storage.NewWithdraw(ts.ctx, &model.Withdrawal{
				UserID:      tc.userID,
				OrderID:     tc.orderID,
				Sum:         tc.sum,
				ProcessedAt: time.Now(),
			})
			ts.Assert().ErrorIs(err, tc.wantErr)
		})
	}
	ts.Run("#7: check Dave's withdrawals log and balance", func() {
		daveW, err := ts.storage.WithdrawalsByUserID(ts.ctx, dave.ID)
		ts.Assert().NoError(err)
		ts.Assert().Equal(3, len(daveW))
		ts.Assert().Equal(model.StatusProcessed, daveW[0].Status)
		ts.Assert().Equal(model.StatusInvalid, daveW[1].Status)
		ts.Assert().Equal(model.StatusProcessed, daveW[0].Status)
		dave, err = ts.storage.UserByLogin(ts.ctx, dave.Login)
		ts.Require().NoError(err)
		ts.Assert().EqualValues(300-200-99.5, dave.GPointsBalance)
	})
	ts.Run("#8: check Bob's withdrawals log", func() {
		bobW, err := ts.storage.WithdrawalsByUserID(ts.ctx, ts.bob.user.ID)
		ts.Assert().NoError(err)
		ts.Assert().Equal(1, len(bobW))
		ts.Assert().Equal(model.StatusInvalid, bobW[0].Status)
	})
	ts.Run("#9: check Alice's withdrawals log", func() {
		aliceW, err := ts.storage.WithdrawalsByUserID(ts.ctx, ts.alice.user.ID)
		ts.Assert().NoError(err)
		ts.Assert().Equal(1, len(aliceW))
		ts.Assert().Equal(model.StatusInvalid, aliceW[0].Status)
	})
}
