// Code generated by mockery v1.0.0. DO NOT EDIT.

package backup

import (
	mock "github.com/stretchr/testify/mock"

	models "github.com/percona/pmm-managed/models"

	time "time"
)

// mockJobsService is an autogenerated mock type for the jobsService type
type mockJobsService struct {
	mock.Mock
}

// StartMySQLBackupJob provides a mock function with given fields: id, pmmAgentID, timeout, name, dbConfig, locationConfig
func (_m *mockJobsService) StartMySQLBackupJob(id string, pmmAgentID string, timeout time.Duration, name string, dbConfig models.DBConfig, locationConfig models.BackupLocationConfig) error {
	ret := _m.Called(id, pmmAgentID, timeout, name, dbConfig, locationConfig)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, time.Duration, string, models.DBConfig, models.BackupLocationConfig) error); ok {
		r0 = rf(id, pmmAgentID, timeout, name, dbConfig, locationConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StopJob provides a mock function with given fields: jobID
func (_m *mockJobsService) StopJob(jobID string) error {
	ret := _m.Called(jobID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(jobID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}