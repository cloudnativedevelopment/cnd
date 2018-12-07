package model

import (
	"fmt"
	"os"
	"testing"
)

func Test_fixPath(t *testing.T) {
	wd, _ := os.Getwd()

	var tests = []struct {
		name     string
		source   string
		target   string
		devPath  string
		expected string
	}{
		{
			name:     "relative-source",
			source:   ".",
			target:   "/go/src/github.com/okteto/cnd",
			devPath:  "/go/src/github.com/okteto/cnd/cnd.yml",
			expected: "/go/src/github.com/okteto/cnd"},
		{
			name:     "relative-source-abs",
			source:   "/go/src/github.com/okteto/cnd",
			target:   "/src/github.com/okteto/cnd",
			devPath:  "cnd.yml",
			expected: "/go/src/github.com/okteto/cnd"},
		{
			name:     "relative-dev-path",
			source:   "k8/src",
			target:   "/go/src/github.com/okteto/cnd",
			devPath:  "cnd/cnd.yml",
			expected: fmt.Sprintf("%s/cnd/k8/src", wd),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := Dev{
				Mount: mount{
					Source: tt.source,
					Target: tt.target,
				},
			}

			dev.fixPath(tt.devPath)
			if dev.Mount.Source != tt.expected {
				t.Errorf("%s != %s", dev.Mount.Source, tt.expected)
			}
		})
	}
}

func Test_loadDev(t *testing.T) {
	manifest := []byte(`
swap:
  deployment:
    name: deployment
    container: core
    image: codescope/core:0.1.8
    command: ["uwsgi"]
    args: ["--gevent", "100", "--http-socket", "0.0.0.0:8000", "--mount", "/=codescope:app", "--python-autoreload", "1"]
mount:
  source: /Users/fernandomayofernandez/PycharmProjects/codescope-core
  target: /app`)
	d, err := loadDev(manifest)
	if err != nil {
		t.Fatal(err)
	}

	if d.Swap.Deployment.Name != "deployment" {
		t.Errorf("name was not parsed: %+v", d)
	}

	if len(d.Swap.Deployment.Command) != 1 || d.Swap.Deployment.Command[0] != "uwsgi" {
		t.Errorf("command was not parsed: %+v", d)
	}

	if len(d.Swap.Deployment.Args) != 8 || d.Swap.Deployment.Args[4] != "--mount" {
		t.Errorf("args was not parsed: %+v", d)
	}
}

func Test_loadDevDefaults(t *testing.T) {
	manifest := []byte(`
swap:
  deployment:
    name: service
    container: core
mount:
  source: /Users/fernandomayofernandez/PycharmProjects/codescope-core
  target: /app`)
	d, err := loadDev(manifest)
	if err != nil {
		t.Fatal(err)
	}

	if d.Swap.Deployment.Command != nil || len(d.Swap.Deployment.Command) != 0 {
		t.Errorf("command was not parsed: %+v", d)
	}

	if d.Swap.Deployment.Args != nil || len(d.Swap.Deployment.Args) != 0 {
		t.Errorf("args was not parsed: %+v", d)
	}
}
