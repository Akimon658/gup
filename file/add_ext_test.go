package file

import (
	"os"
	"testing"

	"github.com/jaswdr/faker"
)

type testCase struct {
	name     string
	input    string
	expected string
}

var (
	withExt    string
	withoutExt string
)

func TestMain(m *testing.M) {
	withoutExt = faker.New().Lorem().Word()
	withExt = withoutExt + extWin

	os.Exit(m.Run())
}

func TestAddExt(t *testing.T) {
	testCases := []testCase{
		{
			name:     "without extension",
			input:    withoutExt,
			expected: withExt,
		},
		{
			name:     "with extension",
			input:    withExt,
			expected: withExt,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			result := AddExt(v.input)

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
