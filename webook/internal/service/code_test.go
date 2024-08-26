package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeGenerate(t *testing.T) {
	output := fmt.Sprintf("%06d", 1)
	expectedOutput := "000001"

	assert.Equal(t, expectedOutput, output)
}
