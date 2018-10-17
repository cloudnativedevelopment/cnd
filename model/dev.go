package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Dev represents a cloud native development environment
type Dev struct {
	Name  string `yaml:"name"`
	Swap  swap   `yaml:"swap"`
	Mount mount  `yaml:"mount"`
}

type swap struct {
	Deployment deployment `yaml:"deployment"`
	Service    service    `yaml:"service"`
}

type mount struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

func (dev *Dev) validate() error {
	file, err := os.Stat(dev.Mount.Source)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Source mount folder does not exists")
	}
	if !file.Mode().IsDir() {
		return fmt.Errorf("Source mount folder is not a directory")
	}
	if dev.Swap.Deployment.File == "" {
		return fmt.Errorf("Swap deployment file cannot be empty")
	}
	if dev.Swap.Deployment.Image == "" {
		return fmt.Errorf("Swap deployment image cannot be empty")
	}
	return nil
}

//ReadDev returns a Dev object from a given file
func ReadDev(devPath string) (*Dev, error) {
	b, err := ioutil.ReadFile(devPath)
	if err != nil {
		return nil, err
	}

	d, err := loadDev(b)
	d.fixPath(devPath)
	return d, nil
}

func loadDev(b []byte) (*Dev, error) {
	dev := Dev{
		Mount: mount{
			Source: ".",
			Target: "/src",
		},
		Swap: swap{
			Deployment: deployment{
				Command: []string{"tail"},
				Args:    []string{"-f", "/dev/null"},
			},
		},
	}

	err := yaml.Unmarshal(b, &dev)
	if err != nil {
		return nil, err
	}

	if err := dev.validate(); err != nil {
		return nil, err
	}

	return &dev, nil
}

func (dev *Dev) fixPath(originalPath string) {
	if !filepath.IsAbs(dev.Mount.Source) {
		if filepath.IsAbs(originalPath) {
			dev.Mount.Source = path.Join(path.Dir(originalPath), dev.Mount.Source)
		} else {
			wd, _ := os.Getwd()
			dev.Mount.Source = path.Join(wd, path.Dir(originalPath), dev.Mount.Source)
		}
	}
}
