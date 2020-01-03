package model

import (
	"reflect"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestForward_MarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		data     Forward
	}{
		{
			name:     "basic",
			expected: "8080:9090",
			data:     Forward{Local: 8080, Remote: 9090},
		},
		{
			name:     "service",
			expected: "8080:svc",
			data:     Forward{Local: 8080, Remote: 0, Service: true, ServiceName: "svc"},
		},
		{
			name:     "service-with-port",
			expected: "8080:svc:5214",
			data:     Forward{Local: 8080, Remote: 5214, Service: true, ServiceName: "svc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := yaml.Marshal(tt.data)
			if err != nil {
				t.Error(err)
			}

			outStr := strings.Trim(string(b), "\n")
			if outStr != tt.expected {
				t.Errorf("didn't marshal correctly. Actual '%+v', Expected '%+v'", outStr, tt.expected)
			}

		})
	}
}

func TestForward_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expected  Forward
		expectErr bool
	}{
		{
			name:     "basic",
			data:     "8080:9090",
			expected: Forward{Local: 8080, Remote: 9090},
		},
		{
			name:     "equal",
			data:     "8080:8080",
			expected: Forward{Local: 8080, Remote: 8080},
		},
		{
			name:      "service",
			data:      "8080:svc",
			expectErr: false,
			expected:  Forward{Local: 8080, Remote: 0, Service: true, ServiceName: "svc"},
		},
		{
			name:      "service-with-port",
			data:      "8080:svc:5214",
			expectErr: false,
			expected:  Forward{Local: 8080, Remote: 5214, Service: true, ServiceName: "svc"},
		},
		{
			name:      "service-with-bad-port",
			data:      "8080:svc:bar",
			expectErr: true,
		},
		{
			name:      "too-many-parts",
			data:      "8080:8081:8082",
			expectErr: true,
		},
		{
			name:      "service-at-end",
			data:      "8080:8081:svc",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Forward
			if err := yaml.Unmarshal([]byte(tt.data), &result); err != nil {
				if tt.expectErr {
					return
				}

				t.Fatal(err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("didn't unmarshal correctly. Actual '%+v', Expected '%+v'", result, tt.expected)
			}

			out, err := yaml.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			outStr := string(out)
			outStr = strings.TrimSuffix(outStr, "\n")

			if !reflect.DeepEqual(outStr, tt.data) {
				t.Errorf("didn't unmarshal correctly. Actual '%+v', Expected '%+v'", outStr, tt.data)
			}
		})
	}
}
