// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Grabber is an autogenerated mock type for the Grabber type
type Grabber struct {
	mock.Mock
}

type Grabber_Expecter struct {
	mock *mock.Mock
}

func (_m *Grabber) EXPECT() *Grabber_Expecter {
	return &Grabber_Expecter{mock: &_m.Mock}
}

// GetSnapshot provides a mock function with given fields:
func (_m *Grabber) GetSnapshot() map[string]string {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// Grabber_GetSnapshot_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSnapshot'
type Grabber_GetSnapshot_Call struct {
	*mock.Call
}

// GetSnapshot is a helper method to define mock.On call
func (_e *Grabber_Expecter) GetSnapshot() *Grabber_GetSnapshot_Call {
	return &Grabber_GetSnapshot_Call{Call: _e.mock.On("GetSnapshot")}
}

func (_c *Grabber_GetSnapshot_Call) Run(run func()) *Grabber_GetSnapshot_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Grabber_GetSnapshot_Call) Return(_a0 map[string]string) *Grabber_GetSnapshot_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Grabber_GetSnapshot_Call) RunAndReturn(run func() map[string]string) *Grabber_GetSnapshot_Call {
	_c.Call.Return(run)
	return _c
}

// NewGrabber creates a new instance of Grabber. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGrabber(t interface {
	mock.TestingT
	Cleanup(func())
}) *Grabber {
	mock := &Grabber{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
