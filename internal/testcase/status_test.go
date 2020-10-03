package testcase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_MarshalJSON(t *testing.T) {
	statusIn := Accepted
	x, err := statusIn.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"Accepted\""), x)

	statusIn = TimeLimitExceeded
	x, err = statusIn.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"TimeLimitExceeded\""), x)

	statusIn = RuntimeError
	x, err = statusIn.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"RuntimeError\""), x)

	statusIn = InternalError
	x, err = statusIn.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"InternalError\""), x)

	statusIn = WrongAnswer
	x, err = statusIn.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"WrongAnswer\""), x)
}

func TestStatus_UnmarshalJSON(t *testing.T) {
	var statusOut Status

	err := statusOut.UnmarshalJSON([]byte("\"SomeInvalidData\""))
	assert.Error(t, err)

	err = statusOut.UnmarshalJSON([]byte("\"Accepted\""))
	assert.NoError(t, err)
	assert.Equal(t, Accepted, statusOut)

	err = statusOut.UnmarshalJSON([]byte("\"WrongAnswer\""))
	assert.NoError(t, err)
	assert.Equal(t, WrongAnswer, statusOut)

	err = statusOut.UnmarshalJSON([]byte("\"InternalError\""))
	assert.NoError(t, err)
	assert.Equal(t, InternalError, statusOut)

	err = statusOut.UnmarshalJSON([]byte("\"TimeLimitExceeded\""))
	assert.NoError(t, err)
	assert.Equal(t, TimeLimitExceeded, statusOut)

	err = statusOut.UnmarshalJSON([]byte("\"RuntimeError\""))
	assert.NoError(t, err)
	assert.Equal(t, RuntimeError, statusOut)
}
