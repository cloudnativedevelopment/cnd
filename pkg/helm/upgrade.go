// Copyright 2020 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"fmt"

	"github.com/okteto/okteto/pkg/model"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

//Upgrade upgrades an okteto stack
func Upgrade(c *action.Upgrade, settings *cli.EnvSettings, s *model.Stack, repoName, chartName, chartVersion string, vals map[string]interface{}, wait bool) error {
	c.Namespace = s.Namespace
	c.Atomic = wait
	c.MaxHistory = 2
	c.Version = chartVersion
	chartPath, err := c.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repoName, chartName), settings)
	if err != nil {
		return err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("error loading stack repositoryt: %s", err)
	}

	_, err = c.Run(s.Name, chart, vals)
	if err != nil {
		return fmt.Errorf("error upgrading stack: %s", err)
	}
	return nil
}
