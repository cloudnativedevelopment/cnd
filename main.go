package main

import (
	"os"

	"github.com/cloudnativedevelopment/cnd/cmd"
	"github.com/cloudnativedevelopment/cnd/pkg/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.Init(logrus.WarnLevel, cmd.GetActionID())
	os.Exit(cmd.Execute())
}
