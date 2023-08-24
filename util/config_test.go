package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	_, errInvalidPath := LoadConfig(".", "app")
	c, err := LoadConfig("../", "app")
	assert.Error(t, errInvalidPath)
	assert.Nil(t, err)
	assert.NotEmpty(t, c)
}
