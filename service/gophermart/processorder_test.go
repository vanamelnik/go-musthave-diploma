package gophermart_test

import (
	"errors"
	"testing"
	"time"

	"github.com/vanamelnik/gophermart/model"
	appContext "github.com/vanamelnik/gophermart/pkg/ctx"
	"github.com/vanamelnik/gophermart/service/gophermart"
	mockstorage "github.com/vanamelnik/gophermart/storage/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := mockstorage.NewMockStorage(mockCtrl)
	ctx, s, err := initServices(db, pepper)
	require.NoError(t, err)
	user := appContext.User(ctx)
	user2 := uuid.New()

	// Perform tests with calls to the mock storage
	db.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	tt1 := []struct {
		name       string
		orderID    model.OrderID
		mockArg    model.OrderID
		mockReturn struct {
			order *model.Order
			err   error
		}
		wantErrSpecified bool // if false, we wait for an unspecified error and ignore wantErr filed.
		wantErr          error
	}{
		{
			name:    "#2 Order already executed by the same user",
			orderID: "00",
			mockArg: "00",
			mockReturn: struct {
				order *model.Order
				err   error
			}{
				order: &model.Order{
					ID:            "00",
					UserID:        user.ID,
					Status:        model.StatusProcessing,
					AccrualPoints: 0,
					UploadedAt:    time.Time{},
				},
				err: nil,
			},
			wantErrSpecified: true,
			wantErr:          gophermart.ErrOrderExecutedBySameUser,
		},
		{
			name:    "#3 Order already executed by another user",
			orderID: "26",
			mockArg: "26",
			mockReturn: struct {
				order *model.Order
				err   error
			}{
				order: &model.Order{
					ID:            "26",
					UserID:        user2,
					Status:        model.StatusProcessing,
					AccrualPoints: 0,
					UploadedAt:    time.Time{},
				},
				err: nil,
			},

			wantErrSpecified: true,
			wantErr:          gophermart.ErrOrderExecutedByAnotherUser,
		},
		{
			name:    "#4 DB order search: unspecified storage error",
			orderID: "18",
			mockArg: "18",
			mockReturn: struct {
				order *model.Order
				err   error
			}{
				order: nil,
				err:   errors.New("some unspecified error"),
			},
			wantErrSpecified: false,
			wantErr:          nil,
		},
	}
	for _, tc := range tt1 {
		t.Run(tc.name, func(t *testing.T) {
			db.EXPECT().
				OrderByID(gomock.Any(), tc.mockArg).
				Return(tc.mockReturn.order, tc.mockReturn.err).
				Times(1)
			err := s.ProcessOrder(ctx, tc.orderID)
			if !tc.wantErrSpecified {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}
