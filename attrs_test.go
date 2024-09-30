package errors_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kapetan-io/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

// NOTE: Tests are sensitive to line changes, only add new tests to the end of this file

func TestAttrsWithCodeLoc(t *testing.T) {
	var w bytes.Buffer
	log := slog.New(slog.NewTextHandler(&w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		}}))

	err := errors.With("foo", "bar").Error("this is an error")
	require.Error(t, err)

	log.LogAttrs(context.Background(), slog.LevelInfo, err.Error(), errors.AttrsWithCodeLoc(err)...)

	assert.Contains(t, w.String(), "msg=\"this is an error\"")
	assert.Contains(t, w.String(), "foo=bar")
	assert.Contains(t, w.String(), "code.filepath=")
	assert.Contains(t, w.String(), "errors/attrs_test.go")
	assert.Contains(t, w.String(), "code.function=github.com/kapetan-io/errors_test.TestAttrsWithCodeLoc")
	assert.Contains(t, w.String(), "code.lineno=26")

	w.Reset()

	err = errors.With(slog.String("friendship", "magic")).Errorf("wrapping the previous: %w", err)
	require.Error(t, err)

	log.LogAttrs(context.Background(), slog.LevelInfo, err.Error(), errors.AttrsWithCodeLoc(err)...)

	assert.Contains(t, w.String(), "msg=\"wrapping the previous: this is an error\"")
	assert.Contains(t, w.String(), "foo=bar")
	assert.Contains(t, w.String(), "friendship=magic")
	assert.Contains(t, w.String(), "code.filepath=")
	assert.Contains(t, w.String(), "errors/attrs_test.go")
	assert.Contains(t, w.String(), "code.function=github.com/kapetan-io/errors_test.TestAttrsWithCodeLoc")
	assert.Contains(t, w.String(), "code.lineno=26")
	//t.Log(buf.String())
	// level=INFO msg="wrapping the previous: this is an error"
	//  friendship=magic
	//  foo=bar
	//  code.filepath=/Users/thrawn/Development/errors/attrs_test.go
	//  code.function=github.com/kapetan-io/errors_test.TestAttrsWithCodeLoc
	//  code.lineno=26
}

func TestAttrs(t *testing.T) {
	err := errors.New("query error")
	wrap := errors.With("key1", "value1").Errorf("message: %w", err)
	assert.NotNil(t, wrap)

	t.Run("UnwrapReturnsWrappedError", func(t *testing.T) {
		u := errors.Unwrap(wrap)
		require.NotNil(t, u)
		assert.Equal(t, "query error", u.Error())
	})

	t.Run("ExtractAttrs", func(t *testing.T) {
		as := errors.AttrsFrom(wrap)
		require.NotNil(t, as)

		assert.True(t, as[0].Equal(slog.Any("key1", "value1")))
		assert.Len(t, as, 1)
	})

	t.Run("IsFromStdErrorsPackage", func(t *testing.T) {
		again := fmt.Errorf("wrap again: %w", wrap)
		assert.True(t, errors.Is(again, &errors.ErrAttrs{}))
	})

	t.Run("AsFromStdErrorsPackage", func(t *testing.T) {
		myErr := &errors.ErrAttrs{}
		again := fmt.Errorf("wrap again: %w", wrap)
		assert.True(t, errors.As(again, &myErr))
		assert.Equal(t, "message: query error", myErr.Error())
	})

	t.Run("AsWithHasAttrs", func(t *testing.T) {
		var a errors.HasAttrs
		require.True(t, errors.As(wrap, &a))

		require.NotNil(t, a)
		b := bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(&b, nil))
		log.LogAttrs(context.Background(), slog.LevelError, "test log attrs", errors.AttrsFrom(a)...)
		assert.Contains(t, b.String(), "test log attrs")
		assert.Contains(t, b.String(), "key1=value1")
	})

	t.Run("WrapShouldReturnNilIfErrorIsNil", func(t *testing.T) {
		got := errors.With("some", "context").Wrap(nil)
		assert.Nil(t, got)
	})

	t.Run("FormatTest", func(t *testing.T) {
		assert.Equal(t, "message: query error", wrap.Error())

		// %s
		assert.Equal(t, `message: query error`, fmt.Sprintf("%s", wrap))
		// %v
		assert.Equal(t, `message: query error`, fmt.Sprintf("%v", wrap))
		// %+v
		assert.Equal(t, `message: query error (key1=value1)`, fmt.Sprintf("%+v", wrap))
	})

	t.Run("With", func(t *testing.T) {
		err := errors.New("query error")
		wrap := errors.WithAttr(slog.String("key", "value")).Errorf("message: %w", err)
		assert.NotNil(t, wrap)
		assert.Equal(t, `message: query error (key=value)`, fmt.Sprintf("%+v", wrap))
	})
}
