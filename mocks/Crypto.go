// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	sspvo "github.com/ftomza/go-sspvo"
	mock "github.com/stretchr/testify/mock"
)

// Crypto is an autogenerated mock type for the Crypto type
type Crypto struct {
	mock.Mock
}

// GetCert provides a mock function with given fields:
func (_m *Crypto) GetCert() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetVerifyCrypto provides a mock function with given fields: cert
func (_m *Crypto) GetVerifyCrypto(cert string) (sspvo.Crypto, error) {
	ret := _m.Called(cert)

	var r0 sspvo.Crypto
	if rf, ok := ret.Get(0).(func(string) sspvo.Crypto); ok {
		r0 = rf(cert)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sspvo.Crypto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(cert)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hash provides a mock function with given fields: data
func (_m *Crypto) Hash(data []byte) []byte {
	ret := _m.Called(data)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// Sign provides a mock function with given fields: digest
func (_m *Crypto) Sign(digest []byte) ([]byte, error) {
	ret := _m.Called(digest)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(digest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(digest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: sign, digest
func (_m *Crypto) Verify(sign []byte, digest []byte) (bool, error) {
	ret := _m.Called(sign, digest)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, []byte) bool); ok {
		r0 = rf(sign, digest)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, []byte) error); ok {
		r1 = rf(sign, digest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}