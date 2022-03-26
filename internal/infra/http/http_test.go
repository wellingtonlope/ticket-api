package http

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapError(t *testing.T) {
	t.Run("should return an error in json format", func(t *testing.T) {
		got := wrapError(errors.New("i'm an error"))
		assert.Equal(t, "{\"message\":\"i'm an error\"}", got)
	})

	t.Run("should return an json without message", func(t *testing.T) {
		got := wrapError(nil)
		assert.Equal(t, "{\"message\":\"\"}", got)
	})
}

func TestWrapBody(t *testing.T) {
	t.Run("should return an json with the body", func(t *testing.T) {
		type test struct {
			Foo string `json:"foo"`
		}
		got := wrapBody(test{Foo: "bar"})
		assert.Equal(t, "{\"foo\":\"bar\"}", got)
	})

	t.Run("should return a empty json", func(t *testing.T) {
		got := wrapBody(nil)
		assert.Equal(t, "{}", got)
	})
}
