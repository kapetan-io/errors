package errors_test

import (
	"fmt"
	"testing"

	"github.com/kapetan-io/errors"
	"github.com/stretchr/testify/assert"
)

func TestLast(t *testing.T) {
	err := errors.New("bottom")
	err = errors.With("sonic", "boom").Errorf("last: %w", err)
	err = fmt.Errorf("second: %w", err)
	err = errors.With("key", "value").Errorf("first: %w", err)
	err = fmt.Errorf("top: %w", err)

	// errors.As() returns the "first" error in the chain with attributes
	var first errors.HasAttrs
	assert.True(t, errors.As(err, &first))
	assert.Equal(t, "first: second: last: bottom", first.(error).Error())

	// errors.Last() returns the last error in the chain with attributes
	var last errors.HasAttrs
	assert.True(t, errors.Last(err, &last))
	assert.Equal(t, "last: bottom", last.(error).Error())

	// If no attributes are found, then should not set target and should return false
	assert.False(t, errors.Last(errors.New("no attributes"), &last))
	assert.Equal(t, "last: bottom", last.(error).Error())
}

func TestErrorf(t *testing.T) {
	err := errors.Errorf("wrap: %w", errors.New("error"))
	assert.EqualError(t, err, "wrap: error")

	var a errors.HasAttrs
	assert.True(t, errors.As(err, &a))
	as, pc := a.Attrs()
	assert.True(t, pc != 0)
	assert.Equal(t, 0, len(as))
}

func TestError(t *testing.T) {
	err := errors.Error("error")
	assert.EqualError(t, err, "error")

	var a errors.HasAttrs
	assert.True(t, errors.As(err, &a))
	as, pc := a.Attrs()
	assert.True(t, pc != 0)
	assert.Equal(t, 0, len(as))
}

func TestWrap(t *testing.T) {
	err := errors.Wrap(errors.New("error"))

	var a errors.HasAttrs
	assert.True(t, errors.As(err, &a))
	as, pc := a.Attrs()
	assert.True(t, pc != 0)
	assert.Equal(t, 0, len(as))
}
