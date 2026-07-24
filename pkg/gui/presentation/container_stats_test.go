package presentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertBigMetricFromSchema(t *testing.T) {
	type scenario struct {
		name        string
		data        map[string]interface{}
		expected    map[string]interface{}
		expectedErr string
	}

	scenarios := []scenario{
		{
			name: "string schema: converts int64 bytes",
			data: map[string]interface{}{
				"ClientStats": map[string]interface{}{
					"memory_stats": map[string]interface{}{
						"limit": int64(2048),
					},
				},
			},
			expected: map[string]interface{}{
				"ClientStats": map[string]interface{}{
					"memory_stats": map[string]interface{}{
						"limit": "2.000KB",
					},
				},
			},
		},
		{
			name: "string schema: converts int64 nanoseconds",
			data: map[string]interface{}{
				"ClientStats": map[string]interface{}{
					"cpu_stats": map[string]interface{}{
						"cpu_usage": map[string]interface{}{
							"total_usage":         int64(100000000000),
							"percpu_usage":        []int64{50000000000, 50000000000},
							"usage_in_kernelmode": int64(200000),
							"usage_in_usermode":   int64(200000),
						},
						"system_cpu_usage": int64(100000000000),
					},
				},
			},
			expected: map[string]interface{}{
				"ClientStats": map[string]interface{}{
					"cpu_stats": map[string]interface{}{
						"cpu_usage": map[string]interface{}{
							"total_usage":         "100.000s",
							"percpu_usage":        []string{"50.000s", "50.000s"},
							"usage_in_kernelmode": "200.000µs",
							"usage_in_usermode":   "200.000µs",
						},
						"system_cpu_usage": "100.000s",
					},
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			data := s.data
			err := convertBigMetricFromSchema(&data)
			if s.expectedErr != "" {
				assert.EqualError(t, err, s.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, s.expected, data)
		})
	}
}
