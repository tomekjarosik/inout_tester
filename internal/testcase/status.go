package testcase

import (
	"encoding/json"
	"errors"
)

//go:generate stringer -type=Status

// Status status of the
type Status int

const (
	// NotRunYet test is waiting to be run
	NotRunYet Status = iota + 1
	// InternalError something unexpected went wrong
	InternalError
	// TimeLimitExceeded the test took too long to process
	TimeLimitExceeded
	// MemoryLimitExceeded the test run used too much RAM
	MemoryLimitExceeded
	// WrongAnswer test run successfully but test outputs differ
	WrongAnswer
	// Accepted all OK
	Accepted
)

func (e *Status) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	for i := 0; i <= len(_Status_index); i++ {
		if Status(i).String() == s {
			*e = Status(i)
			return nil
		}
	}
	return errors.New("invalid testacase status value")
}

func (e Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}
