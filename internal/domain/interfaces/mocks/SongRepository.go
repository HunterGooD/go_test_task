// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/HunterGooD/go_test_task/internal/domain/entity"
	interfaces "github.com/HunterGooD/go_test_task/internal/domain/interfaces"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// SongRepository is an autogenerated mock type for the SongRepository type
type SongRepository struct {
	mock.Mock
}

// DeleteForceByID provides a mock function with given fields: ctx, id
func (_m *SongRepository) DeleteForceByID(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteForceByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteForceByName provides a mock function with given fields: ctx, name
func (_m *SongRepository) DeleteForceByName(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for DeleteForceByName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSoftByGroupID provides a mock function with given fields: ctx, id
func (_m *SongRepository) DeleteSoftByGroupID(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSoftByGroupID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSoftByID provides a mock function with given fields: ctx, id
func (_m *SongRepository) DeleteSoftByID(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSoftByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSoftByName provides a mock function with given fields: ctx, name
func (_m *SongRepository) DeleteSoftByName(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSoftByName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSoftSong provides a mock function with given fields: ctx
func (_m *SongRepository) DeleteSoftSong(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSoftSong")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *SongRepository) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *entity.Song
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*entity.Song, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Song); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Song)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: ctx, name
func (_m *SongRepository) GetByName(ctx context.Context, name string) (*entity.Song, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

	var r0 *entity.Song
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.Song, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.Song); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Song)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetListSong provides a mock function with given fields: ctx, offset, limit, filters
func (_m *SongRepository) GetListSong(ctx context.Context, offset int, limit int, filters *entity.SongFilters) ([]entity.Song, error) {
	ret := _m.Called(ctx, offset, limit, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetListSong")
	}

	var r0 []entity.Song
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, *entity.SongFilters) ([]entity.Song, error)); ok {
		return rf(ctx, offset, limit, filters)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, *entity.SongFilters) []entity.Song); ok {
		r0 = rf(ctx, offset, limit, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Song)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, *entity.SongFilters) error); ok {
		r1 = rf(ctx, offset, limit, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongTextByID provides a mock function with given fields: ctx, songID
func (_m *SongRepository) GetSongTextByID(ctx context.Context, songID int64) (string, error) {
	ret := _m.Called(ctx, songID)

	if len(ret) == 0 {
		panic("no return value specified for GetSongTextByID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (string, error)); ok {
		return rf(ctx, songID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) string); ok {
		r0 = rf(ctx, songID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, songID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSongTextByName provides a mock function with given fields: ctx, name
func (_m *SongRepository) GetSongTextByName(ctx context.Context, name string) (*entity.Song, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetSongTextByName")
	}

	var r0 *entity.Song
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.Song, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.Song); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Song)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Total provides a mock function with given fields: ctx
func (_m *SongRepository) Total(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Total")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFromMap provides a mock function with given fields: ctx, fields
func (_m *SongRepository) UpdateFromMap(ctx context.Context, fields map[string]string) (*entity.Song, error) {
	ret := _m.Called(ctx, fields)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFromMap")
	}

	var r0 *entity.Song
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) (*entity.Song, error)); ok {
		return rf(ctx, fields)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string) *entity.Song); ok {
		r0 = rf(ctx, fields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Song)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]string) error); ok {
		r1 = rf(ctx, fields)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTransaction provides a mock function with given fields: tx
func (_m *SongRepository) WithTransaction(tx *sqlx.Tx) interfaces.SongRepository {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 interfaces.SongRepository
	if rf, ok := ret.Get(0).(func(*sqlx.Tx) interfaces.SongRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.SongRepository)
		}
	}

	return r0
}

// NewSongRepository creates a new instance of SongRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSongRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SongRepository {
	mock := &SongRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
