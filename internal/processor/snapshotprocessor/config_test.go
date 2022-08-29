package snapshotprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	assert.NoError(t, Config{}.Validate())
}
