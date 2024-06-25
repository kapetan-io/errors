package errors_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/kapetan-io/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFields(t *testing.T) {
	err := &ErrTest{Msg: "query error"}
	wrap := errors.Fields{"key1", "value1"}.Errorf("message: %w", err)
	assert.NotNil(t, wrap)

	t.Run("Unwrap should return ErrTest", func(t *testing.T) {
		u := errors.Unwrap(wrap)
		require.NotNil(t, u)
		assert.Equal(t, "query error", u.Error())
	})

	t.Run("Extract fields as a normal map", func(t *testing.T) {
		m := errors.ToMap(wrap)
		require.NotNil(t, m)

		assert.Equal(t, "value1", m["key1"])
		assert.Equal(t, "message: query error", m["err"])
		assert.Len(t, m, 2)
	})

	t.Run("Can use errors.Is() from std `errors` package", func(t *testing.T) {
		assert.True(t, errors.Is(err, &ErrTest{}))
		assert.True(t, errors.Is(wrap, &ErrTest{}))
	})

	t.Run("Can use errors.As() from std `errors` package", func(t *testing.T) {
		myErr := &ErrTest{}
		assert.True(t, errors.As(wrap, &myErr))
		assert.Equal(t, myErr.Msg, "query error")
	})

	t.Run("errors.As with HasFields", func(t *testing.T) {
		var f errors.HasFields
		require.True(t, errors.As(wrap, &f))

		require.NotNil(t, f)
		b := bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(&b, nil))
		log.Error("test log fields", f.Fields()...)
		assert.Contains(t, b.String(), "test log fields")
		assert.Contains(t, b.String(), "key1=value1")

		assert.Equal(t, "message: query error", wrap.Error())
		out := fmt.Sprintf("%+v", wrap)
		assert.Contains(t, out, `message: query error (key1=value1)`)

		if errors.As(wrap, &f) {
			slog.Error("this error has log fields", f.Fields()...)
		}
	})

	t.Run("errors.ToAttr() all fields", func(t *testing.T) {
		b := bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(&b, nil))
		log.Error("test log fields", errors.ToAttr(wrap)...)
		assert.Contains(t, b.String(), "test log fields")
		assert.Contains(t, b.String(), `err="message: query error"`)
		assert.Contains(t, b.String(), "key1=value1")

		assert.Equal(t, "message: query error", wrap.Error())
		out := fmt.Sprintf("%+v", wrap)
		assert.True(t, strings.Contains(out, `message: query error (key1=value1)`))
	})

	t.Run("Wrap() should return nil, if error is nil", func(t *testing.T) {
		got := errors.Fields{"some", "context"}.Wrap(nil)
		assert.Nil(t, got)
	})
}

func TestSlogAttributes(t *testing.T) {
	err := errors.Fields{
		"key1", "value1",
		slog.String("key2", "value2"),
		slog.Duration("expired", time.Second),
	}.Errorf("message")
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("%+v", err), "message (key1=value1, key2=value2, expired=1s)")
	slog.Info("an error occurred!", errors.ToAttr(err)...)
	err = errors.New("blah")
	slog.Info("an error occurred!", errors.ToAttr(err)...)

}

func TestErrorf(t *testing.T) {
	err := errors.New("this is an error")
	wrap := errors.Fields{"key1", "value1", "key2", "value2"}.Errorf("message: %w", err)
	err = fmt.Errorf("wrapped: %w", wrap)
	assert.Equal(t, fmt.Sprintf("%s", err), "wrapped: message: this is an error")
}

func TestNestedFields(t *testing.T) {
	err := errors.New("this is an error")
	err = errors.Fields{"key1", "value1"}.Errorf("message: %w", err)
	err = fmt.Errorf("second: %w", err)
	err = errors.Fields{"key2", "value2"}.Errorf("message: %w", err)
	err = fmt.Errorf("first: %w", err)

	t.Run("ToMap() collects all values from nested errors.Fields", func(t *testing.T) {
		m := errors.ToMap(err)
		assert.NotNil(t, m)
		assert.Equal(t, "value1", m["key1"])
		assert.Equal(t, "value2", m["key2"])
	})

	t.Run("ToAttr() collects all values from nested errors.Fields", func(t *testing.T) {
		f := errors.ToAttr(err)
		require.NotNil(t, f)
		b := bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(&b, nil))
		log.Error("test log fields", f...)
		assert.Contains(t, b.String(), "test log fields")
		assert.Contains(t, b.String(), "key1=value1")
		assert.Contains(t, b.String(), "key2=value2")
	})
}

func TestFieldsFmtDirectives(t *testing.T) {
	t.Run("Wrap() with a message", func(t *testing.T) {
		err := errors.Fields{"key1", "value1"}.Errorf("shit happened: %w", errors.New("error"))
		assert.Equal(t, "shit happened: error", fmt.Sprintf("%s", err))
		assert.Equal(t, "shit happened: error", fmt.Sprintf("%v", err))
		assert.Equal(t, "shit happened: error (key1=value1)", fmt.Sprintf("%+v", err))
		assert.Equal(t, "*errors.fields", fmt.Sprintf("%T", err))
	})

	t.Run("Wrap() without a message", func(t *testing.T) {
		err := errors.Fields{"key1", "value1"}.Wrap(errors.New("error"))
		assert.Equal(t, "error", fmt.Sprintf("%s", err))
		assert.Equal(t, "error", fmt.Sprintf("%v", err))
		assert.Equal(t, "error (key1=value1)", fmt.Sprintf("%+v", err))
		assert.Equal(t, "*errors.fields", fmt.Sprintf("%T", err))
	})
}

func TestFieldsErrorValue(t *testing.T) {
	err := io.EOF
	wrap := errors.Fields{"key1", "value1"}.Errorf("message: %w", err)
	assert.True(t, errors.Is(wrap, io.EOF))
}

func TestHasFields(t *testing.T) {
	hf := &ErrHasFields{M: "error", F: []any{"file", "errors.go"}}
	err := errors.Fields{"key1", "value1"}.Wrap(hf)
	m := errors.ToMap(err)
	require.NotNil(t, m)
	assert.Equal(t, "value1", m["key1"])
	assert.Equal(t, "errors.go", m["file"])
}

func TestFieldsError(t *testing.T) {
	t.Run("Fields.Error() should create a new error", func(t *testing.T) {
		err := errors.Fields{"key1", "value1"}.Error("error")
		m := errors.ToMap(err)
		require.NotNil(t, m)
		assert.Equal(t, "value1", m["key1"])
		assert.Equal(t, "error", err.Error())
	})

	t.Run("Fields.Errorf() should create a new error", func(t *testing.T) {
		err := errors.Fields{"key1", "value1"}.Errorf("error '%d'", 1)
		m := errors.ToMap(err)
		require.NotNil(t, m)
		assert.Equal(t, "value1", m["key1"])
		assert.Equal(t, "error '1'", err.Error())
	})
}
