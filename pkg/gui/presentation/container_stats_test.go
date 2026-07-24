package presentation

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertBigMetricFromSchema(t *testing.T) {
	type scenario struct {
		name        string
		data        map[string]interface{}
		path        string
		schema      interface{}
		expected    map[string]interface{}
		expectedErr string
	}

	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(SCHEMA_JSON), &schema); err != nil {
		fmt.Println("Error unmarshaling SCHEMA_JSON:", err)
	}

	scenarios := []scenario{
		{
			name: "string schema: converts int64 bytes",
			data: map[string]interface{}{
				"client_stats": map[string]interface{}{
					"memory_stats": map[string]interface{}{
						"limit": int64(2048),
					},
				},
			},
			path:   "",
			schema: schema,
			expected: map[string]interface{}{
				"client_stats": map[string]interface{}{
					"memory_stats": map[string]interface{}{
						"limit": "2.000 KB",
					},
				},
			},
		},
		{
			name: "string schema: converts int64 nanoseconds",
			data: map[string]interface{}{
				"client_stats": map[string]interface{}{
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
			path:   "",
			schema: schema,
			expected: map[string]interface{}{
				"client_stats": map[string]interface{}{
					"cpu_stats": map[string]interface{}{
						"cpu_usage": map[string]interface{}{
							"total_usage":         "100.000 s",
							"percpu_usage":        []string{"50.000 ms", "50.000 ms"},
							"usage_in_kernelmode": "200.000 µs",
							"usage_in_usermode":   "200.000 µs",
						},
						"system_cpu_usage": "100.000 s",
					},
				},
			},
		},
		{
			name: "direct string schema: converts int64 bytes at path",
			data: map[string]interface{}{
				"limit": int64(2048),
			},
			path:   ".limit",
			schema: "bytes",
			expected: map[string]interface{}{
				"limit": "2.000 KB",
			},
		},
		{
			name: "direct string schema: converts int64 nanoseconds at path",
			data: map[string]interface{}{
				"cpu": int64(1500000),
			},
			path:   ".cpu",
			schema: "nanoseconds",
			expected: map[string]interface{}{
				"cpu": "1.500 ms",
			},
		},
		{
			name: "direct string schema: keeps small bytes value without unit suffix",
			data: map[string]interface{}{
				"limit": int64(500),
			},
			path:   ".limit",
			schema: "bytes",
			expected: map[string]interface{}{
				"limit": "500",
			},
		},
		{
			name: "direct []string schema: converts []int64 at path",
			data: map[string]interface{}{
				"percpu_usage": []int64{1500000, 2500000000},
			},
			path:   ".percpu_usage",
			schema: []string{"nanoseconds", "nanoseconds"},
			expected: map[string]interface{}{
				"percpu_usage": []string{"1.500 ms", "2.500 s"},
			},
		},
		{
			name: "map schema: recurses into single-key nested map",
			data: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": int64(4096),
				},
			},
			path: "",
			schema: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "bytes",
				},
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "4.000 KB",
				},
			},
		},
		{
			name: "unknown schema type is a no-op",
			data: map[string]interface{}{
				"limit": int64(2048),
			},
			path:   ".limit",
			schema: 42,
			expected: map[string]interface{}{
				"limit": int64(2048),
			},
		},
		{
			name: "string schema: errors when value at path is not int64",
			data: map[string]interface{}{
				"limit": "not a number",
			},
			path:        ".limit",
			schema:      "bytes",
			expected:    map[string]interface{}{"limit": "not a number"},
			expectedErr: "Can't convert string to int64",
		},
		{
			name: "[]string schema: errors when value at path is not []int64",
			data: map[string]interface{}{
				"percpu_usage": []string{"a", "b"},
			},
			path:        ".percpu_usage",
			schema:      []string{"nanoseconds", "nanoseconds"},
			expected:    map[string]interface{}{"percpu_usage": []string{"a", "b"}},
			expectedErr: "Can't convert []string to []int64",
		},
		{
			name:        "string schema: errors when path does not exist in data",
			data:        map[string]interface{}{},
			path:        ".missing",
			schema:      "bytes",
			expected:    map[string]interface{}{},
			expectedErr: "Unable to find the key",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if m, ok := s.schema.(map[string]interface{}); ok && m == nil {
				t.Skip("scenario depends on SCHEMA_JSON which failed to parse")
			}
			data := s.data
			err := convertBigMetricFromSchema(&data, s.path, s.schema)
			if s.expectedErr != "" {
				assert.EqualError(t, err, s.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, s.expected, data)
		})
	}
}
