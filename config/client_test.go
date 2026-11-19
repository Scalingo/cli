package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScalingoClientBaseConfigTimeout(t *testing.T) {
	t.Run("it should keep the SDK default timeout when none is requested", func(t *testing.T) {
		cfg := Config{}.scalingoClientBaseConfig(ClientConfigOpts{})

		assert.Zero(t, cfg.Timeout)
	})

	t.Run("it should apply the requested timeout", func(t *testing.T) {
		cfg := Config{}.scalingoClientBaseConfig(ClientConfigOpts{Timeout: 2 * time.Second})

		assert.Equal(t, 2*time.Second, cfg.Timeout)
	})
}
