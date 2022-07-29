package file

import "testing"

func TestTrimExt(t *testing.T) {
	testCases := []testCase{
		{
			name:     "with extension",
			input:    withExt,
			expected: withoutExt,
		},
		{
			name:     "without extension",
			input:    withoutExt,
			expected: withoutExt,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			result := TrimExt(v.input)

			if isWindows() {
				if result != v.expected {
					t.Errorf("expected %s, got %s", v.expected, result)
				}
			} else if result != v.input {
				t.Errorf("expected %s, got %s", v.input, result)
			}
		})
	}
}
