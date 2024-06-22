package errors_test

import (
	"fmt"
	"testing"

	"github.com/kapetan-io/errors"
	"github.com/stretchr/testify/assert"
)

type ErrTest struct {
	Msg string
}

func (e *ErrTest) Error() string {
	return e.Msg
}

func (e *ErrTest) Is(target error) bool {
	_, ok := target.(*ErrTest)
	return ok
}

type ErrHasFields struct {
	M string
	F []any
}

func (e *ErrHasFields) Error() string {
	return e.M
}

func (e *ErrHasFields) Is(target error) bool {
	_, ok := target.(*ErrHasFields)
	return ok
}

func (e *ErrHasFields) Fields() []any {
	return e.F
}

func TestLast(t *testing.T) {
	err := errors.New("bottom")
	err = errors.Fields{"sonic", "boom"}.Errorf("last: %w", err)
	err = fmt.Errorf("second: %w", err)
	err = errors.Fields{"key", "value"}.Errorf("first: %w", err)
	err = fmt.Errorf("top: %w", err)

	// errors.As() returns the "first" error in the chain with fields
	var first errors.HasFields
	assert.True(t, errors.As(err, &first))
	assert.Equal(t, "first: second: last: bottom", first.(error).Error())

	// errors.Last() returns the last error in the chain with fields
	var last errors.HasFields
	assert.True(t, errors.Last(err, &last))
	assert.Equal(t, "last: bottom", last.(error).Error())

	// If no fields are found, then should not set target and should return false
	assert.False(t, errors.Last(errors.New("no fields"), &last))
	assert.Equal(t, "last: bottom", last.(error).Error())
}
