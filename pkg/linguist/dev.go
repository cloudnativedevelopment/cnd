package linguist

import (
	"strings"

	"github.com/cloudnativedevelopment/cnd/pkg/log"
	"github.com/cloudnativedevelopment/cnd/pkg/model"
)

type languageDefault struct {
	image   string
	command []string
	path    string
	scripts map[string]string
	forward model.Forward
}

const (
	javascript       = "javascript"
	golang           = "go"
	python           = "python"
	java             = "java"
	ruby             = "ruby"
	unrecognized     = "unrecognized"
	helloCommandName = "hello"
)

var (
	languageDefaults map[string]languageDefault
	tailCommand      = []string{"tail", "-f", "/dev/null"}
)

func init() {
	languageDefaults = make(map[string]languageDefault)
	languageDefaults[javascript] = languageDefault{
		image:   "okteto/node:11",
		command: tailCommand,
		path:    "/usr/src/app",
		scripts: map[string]string{
			"test":  "yarn run test",
			"start": "yarn install && yarn start",
		},
		forward: model.Forward{Local: 3000, Remote: 3000},
	}

	languageDefaults[golang] = languageDefault{
		image:   "golang:1",
		command: tailCommand,
		path:    "/go/src/app",
		scripts: map[string]string{
			"start": "go run main.go",
		},
		forward: model.Forward{Local: 8080, Remote: 8080},
	}

	languageDefaults[python] = languageDefault{
		image:   "python:3",
		command: tailCommand,
		path:    "/usr/src/app",
		scripts: map[string]string{
			"start": "pip install -r requirements.txt && python app.py",
		},
		forward: model.Forward{Local: 8080, Remote: 8080},
	}

	languageDefaults[java] = languageDefault{
		image:   "gradle:5.1-jdk11",
		command: tailCommand,
		path:    "/home/gradle",
		scripts: map[string]string{
			"boot":  "gradle bootRun",
			"start": "gradle build -continuous --scan",
		},
		forward: model.Forward{Local: 8080, Remote: 8080},
	}

	languageDefaults[ruby] = languageDefault{
		image:   "ruby:2",
		command: tailCommand,
		path:    "/usr/src/app",
		scripts: map[string]string{
			"migrate": "rails db:migrate",
			"start":   "rails s -e development",
		},
		forward: model.Forward{Local: 3000, Remote: 3000},
	}

	languageDefaults[unrecognized] = languageDefault{
		image:   "ubuntu:bionic",
		command: tailCommand,
		path:    "/usr/src/app",
		forward: model.Forward{Local: 8080, Remote: 8080},
	}
}

// GetDevConfig returns the default dev for the specified language
func GetDevConfig(language string) *model.Dev {
	vals := languageDefaults[normalizeLanguage(language)]
	dev := model.NewDev()
	dev.Swap.Deployment.Image = vals.image
	dev.Swap.Deployment.Command = vals.command
	dev.Mount.Source = "."
	dev.Mount.Target = vals.path
	dev.Scripts = vals.scripts

	if dev.Scripts == nil {
		dev.Scripts = make(map[string]string)
	}

	dev.Scripts[helloCommandName] = "echo Your cluster ♥s you"

	dev.Forward = []model.Forward{vals.forward}
	return dev
}

func normalizeLanguage(language string) string {
	lower := strings.ToLower(language)
	switch lower {
	case "typescript":
		return javascript
	case "javascript":
		return javascript
	case "jsx":
		return javascript
	case "python":
		return python
	case "java":
		return java
	case "ruby":
		return ruby
	case "go":
		return golang
	default:
		log.Debugf("unrecognized language: %s", lower)
		return unrecognized
	}
}
