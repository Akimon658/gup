package file

import (
	"runtime"
	"testing"

	"github.com/jaswdr/faker"
)

func TestAddExt(t *testing.T) {
	fileName := faker.New().Lorem().Word()
	fileNameWithExtension := fileName + extWin
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "without extension",
			input:    fileName,
			expected: fileNameWithExtension,
		},
		{
			name:     "with extension",
			input:    fileNameWithExtension,
			expected: fileNameWithExtension,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			result := AddExt(v.input)
			if runtime.GOOS == "windows" {
				if result != v.expected {
					t.Errorf("expected %s, got %s", v.expected, result)
				}
			} else {
				if result != v.input {
					t.Errorf("expected %s, got %s", v.input, result)
				}
			}
		})
	}
}
