// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockHasher is an autogenerated mock type for the Hasher type
type MockHasher struct {
	mock.Mock
}

type MockHasher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHasher) EXPECT() *MockHasher_Expecter {
	return &MockHasher_Expecter{mock: &_m.Mock}
}

// Compare provides a mock function with given fields: givenPassword, hashedPassword
func (_m *MockHasher) Compare(givenPassword string, hashedPassword string) bool {
	ret := _m.Called(givenPassword, hashedPassword)

	if len(ret) == 0 {
		panic("no return value specified for Compare")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(givenPassword, hashedPassword)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockHasher_Compare_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Compare'
type MockHasher_Compare_Call struct {
	*mock.Call
}

// Compare is a helper method to define mock.On call
//   - givenPassword string
//   - hashedPassword string
func (_e *MockHasher_Expecter) Compare(givenPassword interface{}, hashedPassword interface{}) *MockHasher_Compare_Call {
	return &MockHasher_Compare_Call{Call: _e.mock.On("Compare", givenPassword, hashedPassword)}
}

func (_c *MockHasher_Compare_Call) Run(run func(givenPassword string, hashedPassword string)) *MockHasher_Compare_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockHasher_Compare_Call) Return(_a0 bool) *MockHasher_Compare_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHasher_Compare_Call) RunAndReturn(run func(string, string) bool) *MockHasher_Compare_Call {
	_c.Call.Return(run)
	return _c
}

// Hash provides a mock function with given fields: password
func (_m *MockHasher) Hash(password string) (string, error) {
	ret := _m.Called(password)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(password)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHasher_Hash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Hash'
type MockHasher_Hash_Call struct {
	*mock.Call
}

// Hash is a helper method to define mock.On call
//   - password string
func (_e *MockHasher_Expecter) Hash(password interface{}) *MockHasher_Hash_Call {
	return &MockHasher_Hash_Call{Call: _e.mock.On("Hash", password)}
}

func (_c *MockHasher_Hash_Call) Run(run func(password string)) *MockHasher_Hash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockHasher_Hash_Call) Return(_a0 string, _a1 error) *MockHasher_Hash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHasher_Hash_Call) RunAndReturn(run func(string) (string, error)) *MockHasher_Hash_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHasher creates a new instance of MockHasher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHasher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHasher {
	mock := &MockHasher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
