package introspection

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGetStructColumns(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name: "normal case",
			input: struct {
				A string `gorm:"column:a"`
				B int    `gorm:"column:b"`
				C bool   `gorm:"column:c"`
				D float64
				E bool `gorm:"-"`
			}{},
			expected: []string{"a", "b", "c", "d"},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			var cache sync.Map

			actual := getStructColumnsWithCache(c.input, &cache)
			assert.Equal(t, c.expected, actual)
		})
	}
}
