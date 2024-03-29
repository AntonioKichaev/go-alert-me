// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Agent is an autogenerated mock type for the Agent type
type Agent struct {
	mock.Mock
}

type Agent_Expecter struct {
	mock *mock.Mock
}

func (_m *Agent) EXPECT() *Agent_Expecter {
	return &Agent_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with given fields:
func (_m *Agent) Run() {
	_m.Called()
}

// Agent_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type Agent_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
func (_e *Agent_Expecter) Run() *Agent_Run_Call {
	return &Agent_Run_Call{Call: _e.mock.On("Run")}
}

func (_c *Agent_Run_Call) Run(run func()) *Agent_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Agent_Run_Call) Return() *Agent_Run_Call {
	_c.Call.Return()
	return _c
}

func (_c *Agent_Run_Call) RunAndReturn(run func()) *Agent_Run_Call {
	_c.Call.Return(run)
	return _c
}

// NewAgent creates a new instance of Agent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAgent(t interface {
	mock.TestingT
	Cleanup(func())
}) *Agent {
	mock := &Agent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
