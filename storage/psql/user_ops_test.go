package psql

import (
	"time"

	"github.com/vanamelnik/go-musthave-diploma/model"
	"github.com/vanamelnik/go-musthave-diploma/pkg/bcrypt"
	"github.com/vanamelnik/go-musthave-diploma/storage"

	"github.com/google/uuid"
)

func (ts *TestSuite) TestCreateUsers() {
	tt := []struct {
		name    string
		login   string
		wantErr error
	}{
		{
			name:    "#1 Create Charlie",
			login:   "charlieparker@yahoo.md",
			wantErr: nil,
		},
		{
			name:    "#2 Create Bob again",
			login:   "bobmarley@rambler.ru",
			wantErr: storage.ErrAlreadyProcessed,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			hash, err := bcrypt.BcryptPassword(tc.login, "") // all our fake people use their login as a password!..
			ts.Require().NoError(err)
			c := &model.User{
				ID:             uuid.New(),
				Login:          tc.login,
				PasswordHash:   hash,
				CreatedAt:      time.Now(),
				GPointsBalance: 0,
			}
			err = ts.storage.NewUser(ts.ctx, c)
			ts.Assert().ErrorIs(err, tc.wantErr)
		})
	}
}

func (ts *TestSuite) TestUserByLogin() {
	tt := []struct {
		name    string
		login   string
		wantErr error
	}{
		{
			name:    "#1 Fetch Bob",
			login:   "bobmarley@rambler.ru",
			wantErr: nil,
		},
		{
			name:    "#2 Fetch Alice",
			login:   "alicecooper@yandex.cn",
			wantErr: nil,
		},
		{
			name:    "#3 Fetch non-existed user",
			login:   "paulmccartney@mail.ua",
			wantErr: storage.ErrNotFound,
		},
	}
	for _, tc := range tt {
		ts.Run(tc.name, func() {
			_, err := ts.storage.UserByLogin(ts.ctx, tc.login)
			ts.Assert().ErrorIs(err, tc.wantErr)
		})
	}
}
