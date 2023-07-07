// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Random is an autogenerated mock type for the Random type
type Random struct {
	mock.Mock
}

type Random_Expecter struct {
	mock *mock.Mock
}

func (_m *Random) EXPECT() *Random_Expecter {
	return &Random_Expecter{mock: &_m.Mock}
}

// Int provides a mock function with given fields:
func (_m *Random) Int() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Random_Int_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Int'
type Random_Int_Call struct {
	*mock.Call
}

// Int is a helper method to define mock.On call
func (_e *Random_Expecter) Int() *Random_Int_Call {
	return &Random_Int_Call{Call: _e.mock.On("Int")}
}

func (_c *Random_Int_Call) Run(run func()) *Random_Int_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Random_Int_Call) Return(_a0 int) *Random_Int_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Random_Int_Call) RunAndReturn(run func() int) *Random_Int_Call {
	_c.Call.Return(run)
	return _c
}

// NewRandom creates a new instance of Random. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRandom(t interface {
	mock.TestingT
	Cleanup(func())
}) *Random {
	mock := &Random{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
