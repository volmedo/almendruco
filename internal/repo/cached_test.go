package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPassword(t *testing.T) {
	repoMock := &MockRepo{}
	repoMock.On("GetPassword", "Some User").
		Return("s0m3p4ss", nil).
		Once()
	defer repoMock.AssertExpectations(t)

	cr := NewCachedRepo(repoMock)

	pass, err := cr.GetPassword("Some User")
	assert.NoError(t, err)
	assert.Equal(t, "s0m3p4ss", pass)

	// Call the method a second time. If it works properly, cached repo shouldn't call the base repo again.
	// If it does, the mock will return an error, as it is configured to expect only one call to GetPassword
	pass, err = cr.GetPassword("Some User")
	assert.NoError(t, err)
	assert.Equal(t, "s0m3p4ss", pass)
}

func TestGetLastNotifiedMessage(t *testing.T) {
	repoMock := &MockRepo{}
	repoMock.On("GetLastNotifiedMessage", "Some User").
		Return(uint64(123456), nil).
		Once()
	defer repoMock.AssertExpectations(t)

	cr := NewCachedRepo(repoMock)

	last, err := cr.GetLastNotifiedMessage("Some User")
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), last)

	// Call the method a second time. If it works properly, cached repo shouldn't call the base repo again.
	// If it does, the mock will return an error, as it is configured to expect only one call to GetPassword
	last, err = cr.GetLastNotifiedMessage("Some User")
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), last)
}

func TestSetLastNotifiedMessage(t *testing.T) {
	repoMock := &MockRepo{}
	repoMock.On("SetLastNotifiedMessage", "Some User", uint64(123456)).
		Return(nil).
		Once()
	defer repoMock.AssertExpectations(t)

	cr := NewCachedRepo(repoMock)

	err := cr.SetLastNotifiedMessage("Some User", uint64(123456))
	assert.NoError(t, err)

	// SetLastNotifiedMessage should cache the last notified message ID, so GetLastNotifiedMessage shouldn't
	// be called on the base repo when calling GetLastNotifiedMessage on the cached repo
	defer repoMock.AssertNotCalled(t, "GetLastNotifiedMessage", mock.Anything)
	last, err := cr.GetLastNotifiedMessage("Some User")
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), last)
}
