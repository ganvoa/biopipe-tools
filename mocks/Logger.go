// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: message
func (_m *Logger) Debug(message string) {
	_m.Called(message)
}

// Debugf provides a mock function with given fields: message, args
func (_m *Logger) Debugf(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Fatal provides a mock function with given fields: err
func (_m *Logger) Fatal(err error) {
	_m.Called(err)
}

// Fatalf provides a mock function with given fields: message, args
func (_m *Logger) Fatalf(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: message
func (_m *Logger) Info(message string) {
	_m.Called(message)
}

// Infof provides a mock function with given fields: message, args
func (_m *Logger) Infof(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: message
func (_m *Logger) Warn(message string) {
	_m.Called(message)
}

// Warnf provides a mock function with given fields: message, args
func (_m *Logger) Warnf(message string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}
