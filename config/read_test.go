package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jaswdr/faker"
)

func TestRead(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	tmpConfigPath := f.Name()
	defer func() {
		f.Close()
		os.Remove(tmpConfigPath)
	}()

	testCases := []struct {
		name       string
		configPath string
		yaml       string
		expected   Config
	}{
		{
			name: "normal",
			expected: Config{
				Global: BuildFlags{Ldflags: "-s -w"},
				Packages: []struct {
					Name       string
					BuildFlags `yaml:",inline"`
				}{
					{
						Name:       "hugo",
						BuildFlags: BuildFlags{Tags: "extended"},
					},
				},
			},
			yaml: `
global:
  ldflags: -s -w
packages:
  - name: hugo
    tags: extended`,
		},
		{
			name:       "no config",
			configPath: faker.New().File().FilenameWithExtension(),
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			configPath := v.configPath
			if configPath == "" {
				configPath = tmpConfigPath
				f.WriteString(v.yaml)
			}

			actual, err := Read(configPath)
			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(v.expected, *actual); diff != "" {
				t.Errorf("unexpected result:\n%s", diff)
			}
		})
	}
}
