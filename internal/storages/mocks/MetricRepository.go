// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MetricRepository is an autogenerated mock type for the MetricRepository type
type MetricRepository struct {
	mock.Mock
}

type MetricRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MetricRepository) EXPECT() *MetricRepository_Expecter {
	return &MetricRepository_Expecter{mock: &_m.Mock}
}

// AddCounter provides a mock function with given fields: metricName, value
func (_m *MetricRepository) AddCounter(metricName string, value int64) {
	_m.Called(metricName, value)
}

// MetricRepository_AddCounter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddCounter'
type MetricRepository_AddCounter_Call struct {
	*mock.Call
}

// AddCounter is a helper method to define mock.On call
//   - metricName string
//   - value int64
func (_e *MetricRepository_Expecter) AddCounter(metricName interface{}, value interface{}) *MetricRepository_AddCounter_Call {
	return &MetricRepository_AddCounter_Call{Call: _e.mock.On("AddCounter", metricName, value)}
}

func (_c *MetricRepository_AddCounter_Call) Run(run func(metricName string, value int64)) *MetricRepository_AddCounter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int64))
	})
	return _c
}

func (_c *MetricRepository_AddCounter_Call) Return() *MetricRepository_AddCounter_Call {
	_c.Call.Return()
	return _c
}

func (_c *MetricRepository_AddCounter_Call) RunAndReturn(run func(string, int64)) *MetricRepository_AddCounter_Call {
	_c.Call.Return(run)
	return _c
}

// GetCounter provides a mock function with given fields: metricName
func (_m *MetricRepository) GetCounter(metricName string) (int64, error) {
	ret := _m.Called(metricName)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(metricName)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(metricName)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(metricName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MetricRepository_GetCounter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCounter'
type MetricRepository_GetCounter_Call struct {
	*mock.Call
}

// GetCounter is a helper method to define mock.On call
//   - metricName string
func (_e *MetricRepository_Expecter) GetCounter(metricName interface{}) *MetricRepository_GetCounter_Call {
	return &MetricRepository_GetCounter_Call{Call: _e.mock.On("GetCounter", metricName)}
}

func (_c *MetricRepository_GetCounter_Call) Run(run func(metricName string)) *MetricRepository_GetCounter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MetricRepository_GetCounter_Call) Return(_a0 int64, _a1 error) *MetricRepository_GetCounter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MetricRepository_GetCounter_Call) RunAndReturn(run func(string) (int64, error)) *MetricRepository_GetCounter_Call {
	_c.Call.Return(run)
	return _c
}

// GetGauge provides a mock function with given fields: metricName
func (_m *MetricRepository) GetGauge(metricName string) (float64, error) {
	ret := _m.Called(metricName)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (float64, error)); ok {
		return rf(metricName)
	}
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(metricName)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(metricName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MetricRepository_GetGauge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGauge'
type MetricRepository_GetGauge_Call struct {
	*mock.Call
}

// GetGauge is a helper method to define mock.On call
//   - metricName string
func (_e *MetricRepository_Expecter) GetGauge(metricName interface{}) *MetricRepository_GetGauge_Call {
	return &MetricRepository_GetGauge_Call{Call: _e.mock.On("GetGauge", metricName)}
}

func (_c *MetricRepository_GetGauge_Call) Run(run func(metricName string)) *MetricRepository_GetGauge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MetricRepository_GetGauge_Call) Return(_a0 float64, _a1 error) *MetricRepository_GetGauge_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MetricRepository_GetGauge_Call) RunAndReturn(run func(string) (float64, error)) *MetricRepository_GetGauge_Call {
	_c.Call.Return(run)
	return _c
}

// GetMetrics provides a mock function with given fields:
func (_m *MetricRepository) GetMetrics() map[string]string {
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

// MetricRepository_GetMetrics_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMetrics'
type MetricRepository_GetMetrics_Call struct {
	*mock.Call
}

// GetMetrics is a helper method to define mock.On call
func (_e *MetricRepository_Expecter) GetMetrics() *MetricRepository_GetMetrics_Call {
	return &MetricRepository_GetMetrics_Call{Call: _e.mock.On("GetMetrics")}
}

func (_c *MetricRepository_GetMetrics_Call) Run(run func()) *MetricRepository_GetMetrics_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MetricRepository_GetMetrics_Call) Return(_a0 map[string]string) *MetricRepository_GetMetrics_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MetricRepository_GetMetrics_Call) RunAndReturn(run func() map[string]string) *MetricRepository_GetMetrics_Call {
	_c.Call.Return(run)
	return _c
}

// SetGauge provides a mock function with given fields: metricName, value
func (_m *MetricRepository) SetGauge(metricName string, value float64) {
	_m.Called(metricName, value)
}

// MetricRepository_SetGauge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetGauge'
type MetricRepository_SetGauge_Call struct {
	*mock.Call
}

// SetGauge is a helper method to define mock.On call
//   - metricName string
//   - value float64
func (_e *MetricRepository_Expecter) SetGauge(metricName interface{}, value interface{}) *MetricRepository_SetGauge_Call {
	return &MetricRepository_SetGauge_Call{Call: _e.mock.On("SetGauge", metricName, value)}
}

func (_c *MetricRepository_SetGauge_Call) Run(run func(metricName string, value float64)) *MetricRepository_SetGauge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(float64))
	})
	return _c
}

func (_c *MetricRepository_SetGauge_Call) Return() *MetricRepository_SetGauge_Call {
	_c.Call.Return()
	return _c
}

func (_c *MetricRepository_SetGauge_Call) RunAndReturn(run func(string, float64)) *MetricRepository_SetGauge_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMetricRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMetricRepository creates a new instance of MetricRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMetricRepository(t mockConstructorTestingTNewMetricRepository) *MetricRepository {
	mock := &MetricRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
