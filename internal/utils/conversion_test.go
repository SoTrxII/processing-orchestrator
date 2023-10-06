package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatSize_Zero(t *testing.T) {
	assert.Equal(t, "0B", FormatSize(0))
}

func TestFormatSize_BytesOnly(t *testing.T) {
	assert.Equal(t, "1.00B", FormatSize(1))
	assert.Equal(t, "1023.00B", FormatSize(1023))
}

func TestFormatSize_Large(t *testing.T) {
	assert.Equal(t, "1.00KB", FormatSize(1024))
	assert.Equal(t, "1.00MB", FormatSize(1024*1024))
	assert.Equal(t, "3.00TB", FormatSize(3*1024*1024*1024*1024))
}
