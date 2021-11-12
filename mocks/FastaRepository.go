// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	fasta "github.com/ganvoa/biopipe-tools/internal/fasta"
	mock "github.com/stretchr/testify/mock"
)

// FastaRepository is an autogenerated mock type for the FastaRepository type
type FastaRepository struct {
	mock.Mock
}

// GetByStrainId provides a mock function with given fields: strainId
func (_m *FastaRepository) GetByStrainId(strainId int) (*fasta.Strain, error) {
	ret := _m.Called(strainId)

	var r0 *fasta.Strain
	if rf, ok := ret.Get(0).(func(int) *fasta.Strain); ok {
		r0 = rf(strainId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fasta.Strain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(strainId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkAsDownloaded provides a mock function with given fields: strainId
func (_m *FastaRepository) MarkAsDownloaded(strainId int) error {
	ret := _m.Called(strainId)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(strainId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NotDownloaded provides a mock function with given fields: from
func (_m *FastaRepository) NotDownloaded(from int) ([]fasta.Strain, error) {
	ret := _m.Called(from)

	var r0 []fasta.Strain
	if rf, ok := ret.Get(0).(func(int) []fasta.Strain); ok {
		r0 = rf(from)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fasta.Strain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(from)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
