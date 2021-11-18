package gophermart_test

import (
	"context"
	"testing"
	"time"

	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/pkg/bcrypt"
	appContext "github.com/vanamelnik/go-musthave-diploma/pkg/ctx"
	"github.com/vanamelnik/go-musthave-diploma/pkg/logging"
	"github.com/vanamelnik/go-musthave-diploma/service/gophermart"
	"github.com/vanamelnik/go-musthave-diploma/storage"
	mockstorage "github.com/vanamelnik/go-musthave-diploma/storage/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pepper = "custom pepper"

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := mockstorage.NewMockStorage(mockCtrl)
	ctx, s, err := initServices(db, pepper)
	require.NoError(t, err)

	db.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)                           // for test #5
	db.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(storage.ErrLoginAlreadyExists).Times(1) // for test #6

	tt := []struct {
		name     string
		login    string
		password string
		wantNil  bool
		wantErr  bool
	}{
		{
			name:     "#1 empty login and password",
			login:    "",
			password: "",
			wantNil:  true,
			wantErr:  true,
		},
		{
			name:     "#2 short login",
			login:    "ab",
			password: "sidelinatrube",
			wantNil:  true,
			wantErr:  true,
		},
		{
			name:     "#3 short password",
			login:    "abc@def.gh",
			password: "xyz",
			wantNil:  true,
			wantErr:  true,
		},
		{
			name:     "#4 normal case",
			login:    "harry@hogwarts.uk",
			password: "AvadaKedavra!",
			wantNil:  false,
			wantErr:  false,
		},
		{
			name:     "#5 login occupied",
			login:    "billgates@microsoft.com",
			password: "IL0veApples",
			wantNil:  true,
			wantErr:  true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			user, err := s.Create(ctx, tc.login, tc.password)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tc.wantNil {
				assert.Nil(t, user)
			} else {
				assert.NotNil(t, user)
			}
		})
	}
}
func TestAuthenticate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := mockstorage.NewMockStorage(mockCtrl)
	ctx, s, err := initServices(db, pepper)
	require.NoError(t, err)

	// generate hashes for mock responses
	billgatesRightPwdHash, err := bcrypt.BcryptPassword("L1nuxF0rever", pepper)
	require.NoError(t, err)
	hedgehogPwdHash, err := bcrypt.BcryptPassword("Ho0o0o0orse!", pepper)
	require.NoError(t, err)

	db.EXPECT().UserByLogin(gomock.Any(), gomock.Any()).Return(nil, storage.ErrNotFound).Times(1)
	db.EXPECT().UserByLogin(gomock.Any(), "billgates@microsoft.com").Return(&model.User{
		ID:             [16]byte{},
		Login:          "billgates@microsoft.com",
		Password:       "",
		PasswordHash:   billgatesRightPwdHash,
		CreatedAt:      time.Time{},
		GPointsBalance: 0,
	}, nil).Times(1)
	db.EXPECT().UserByLogin(gomock.Any(), "hedgehog@mist.ru").
		Return(&model.User{
			ID:             uuid.New(),
			Login:          "hedgehog@mist.ru",
			Password:       "",
			PasswordHash:   hedgehogPwdHash,
			CreatedAt:      time.Time{},
			GPointsBalance: 0,
		}, nil).Times(1)

	tt := []struct {
		name     string
		login    string
		password string
		wantNil  bool
		wantErr  bool
	}{
		{
			name:     "#1 Nonexistent login address",
			login:    "idont@exi.st",
			password: "nomatterwhatthepasswordis",
			wantNil:  true,
			wantErr:  true,
		},
		{
			name:     "#2 Wrong password",
			login:    "billgates@microsoft.com",
			password: "IL0veApples", // correct password is 'L1nuxF0rever'
			wantNil:  true,
			wantErr:  true,
		},
		{
			name:     "#3 Normal case",
			login:    "hedgehog@mist.ru",
			password: "Ho0o0o0orse!",
			wantNil:  false,
			wantErr:  false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			user, err := s.Authenticate(ctx, tc.login, tc.password)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tc.wantNil {
				assert.Nil(t, user)
			} else {
				assert.NotNil(t, user)
			}
		})
	}
}

func TestGetOrders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := mockstorage.NewMockStorage(mockCtrl)
	ctx, s, err := initServices(db, pepper)
	require.NoError(t, err)
	t.Run("#1 Not authenticated", func(t *testing.T) {
		orders, err := s.GetOrders(appContext.WithLogger(context.Background(),
			logging.NewLogger(logging.WithConsoleOutput(true),
				logging.WithLevel("trace"))))
		assert.ErrorIs(t, err, gophermart.ErrNotAuthenticated)
		assert.Nil(t, orders)
	})
	db.EXPECT().UserOrders(gomock.Any(), gomock.Any()).Return(nil, storage.ErrNotFound).Times(1)
	t.Run("#2 No orders found", func(t *testing.T) {
		orders, err := s.GetOrders(ctx)
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Nil(t, orders)
	})

	o := []model.Order{
		{
			ID:            "12345678",
			UserID:        [16]byte{},
			Status:        "PROCESSED",
			AccrualPoints: 1000,
			UploadedAt:    time.Time{},
		},
		{
			ID:            "87654321",
			UserID:        [16]byte{},
			Status:        "INVALID",
			AccrualPoints: 0,
			UploadedAt:    time.Time{},
		},
	}
	db.EXPECT().UserOrders(gomock.Any(), gomock.Any()).Return(o, nil).Times(1)
	t.Run("#3 Normal case", func(t *testing.T) {
		orders, err := s.GetOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, o, orders)
	})
}

// initService performs initialisation of services, needed for testing.
func initServices(mockdb storage.Storage, pepper string) (context.Context, *gophermart.GopherMart, error) {
	logger := logging.NewLogger(logging.WithConsoleOutput(true), logging.WithLevel("trace"))
	ctx := appContext.WithLogger(context.Background(), logger)
	ctx = appContext.WithUser(ctx, &model.User{
		ID:             uuid.New(),
		Login:          "frodo@hobbyton.shire.me",
		Password:       "TheRingIsM1ne!",
		CreatedAt:      time.Now(),
		GPointsBalance: 0,
	})
	s, err := gophermart.New(ctx, mockdb,
		gophermart.WithConfig(gophermart.Config{PasswordPepper: pepper}),
		gophermart.WithoutWorkers())
	if err != nil {
		return nil, nil, err
	}

	return ctx, s, nil
}
