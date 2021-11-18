// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mockstorage is a generated GoMock package.
package mockstorage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "github.com/vanamelnik/go-musthave-diploma/model"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// CreateAccrual mocks base method.
func (m *MockStorage) CreateAccrual(ctx context.Context, orderID model.OrderID, amount float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccrual", ctx, orderID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccrual indicates an expected call of CreateAccrual.
func (mr *MockStorageMockRecorder) CreateAccrual(ctx, orderID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccrual", reflect.TypeOf((*MockStorage)(nil).CreateAccrual), ctx, orderID, amount)
}

// CreateOrder mocks base method.
func (m *MockStorage) CreateOrder(ctx context.Context, order *model.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockStorageMockRecorder) CreateOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockStorage)(nil).CreateOrder), ctx, order)
}

// CreateUser mocks base method.
func (m *MockStorage) CreateUser(ctx context.Context, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStorageMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorage)(nil).CreateUser), ctx, user)
}

// CreateWithdraw mocks base method.
func (m *MockStorage) CreateWithdraw(ctx context.Context, withdraw *model.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithdraw", ctx, withdraw)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWithdraw indicates an expected call of CreateWithdraw.
func (mr *MockStorageMockRecorder) CreateWithdraw(ctx, withdraw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithdraw", reflect.TypeOf((*MockStorage)(nil).CreateWithdraw), ctx, withdraw)
}

// OrderByID mocks base method.
func (m *MockStorage) OrderByID(ctx context.Context, orderID model.OrderID) (*model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrderByID", ctx, orderID)
	ret0, _ := ret[0].(*model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrderByID indicates an expected call of OrderByID.
func (mr *MockStorageMockRecorder) OrderByID(ctx, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrderByID", reflect.TypeOf((*MockStorage)(nil).OrderByID), ctx, orderID)
}

// OrdersByStatus mocks base method.
func (m *MockStorage) OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OrdersByStatus", ctx, status)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OrdersByStatus indicates an expected call of OrdersByStatus.
func (mr *MockStorageMockRecorder) OrdersByStatus(ctx, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OrdersByStatus", reflect.TypeOf((*MockStorage)(nil).OrdersByStatus), ctx, status)
}

// UpdateBalance mocks base method.
func (m *MockStorage) UpdateBalance(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockStorageMockRecorder) UpdateBalance(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockStorage)(nil).UpdateBalance), ctx)
}

// UpdateOrderStatus mocks base method.
func (m *MockStorage) UpdateOrderStatus(ctx context.Context, orderID model.OrderID, status model.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", ctx, orderID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockStorageMockRecorder) UpdateOrderStatus(ctx, orderID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockStorage)(nil).UpdateOrderStatus), ctx, orderID, status)
}

// UpdateUser mocks base method.
func (m *MockStorage) UpdateUser(ctx context.Context, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStorageMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStorage)(nil).UpdateUser), ctx, user)
}

// UserByLogin mocks base method.
func (m *MockStorage) UserByLogin(ctx context.Context, login string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByLogin", ctx, login)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByLogin indicates an expected call of UserByLogin.
func (mr *MockStorageMockRecorder) UserByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByLogin", reflect.TypeOf((*MockStorage)(nil).UserByLogin), ctx, login)
}

// UserByRemember mocks base method.
func (m *MockStorage) UserByRemember(ctx context.Context, remember string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByRemember", ctx, remember)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByRemember indicates an expected call of UserByRemember.
func (mr *MockStorageMockRecorder) UserByRemember(ctx, remember interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByRemember", reflect.TypeOf((*MockStorage)(nil).UserByRemember), ctx, remember)
}

// UserOrders mocks base method.
func (m *MockStorage) UserOrders(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserOrders", ctx, userID)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserOrders indicates an expected call of UserOrders.
func (mr *MockStorageMockRecorder) UserOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserOrders", reflect.TypeOf((*MockStorage)(nil).UserOrders), ctx, userID)
}

// WithdrawalsByUserID mocks base method.
func (m *MockStorage) WithdrawalsByUserID(ctx context.Context, id uuid.UUID) ([]model.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawalsByUserID", ctx, id)
	ret0, _ := ret[0].([]model.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WithdrawalsByUserID indicates an expected call of WithdrawalsByUserID.
func (mr *MockStorageMockRecorder) WithdrawalsByUserID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawalsByUserID", reflect.TypeOf((*MockStorage)(nil).WithdrawalsByUserID), ctx, id)
}
