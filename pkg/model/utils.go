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

package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/okteto/okteto/pkg/log"
)

// FileExists return true if the file exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		log.Infof("failed to check if %s exists: %s", name, err)
	}

	return true
}

// CopyFile copies a binary between from and to
func CopyFile(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}

	// skipcq GSC-G302 syncthing is a binary so it needs exec permissions
	toFile, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		return err
	}

	defer toFile.Close()

	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return err
	}

	return nil
}

// GetValidNameFromFolder returns a valid kubernetes name for a folder
func GetValidNameFromFolder(folder string) (string, error) {
	dir, err := filepath.Abs(folder)
	if err != nil {
		return "", fmt.Errorf("error inferring name: %s", err)
	}
	name := filepath.Base(dir)
	name = strings.ToLower(name)
	name = ValidKubeNameRegex.ReplaceAllString(name, "-")
	log.Infof("autogenerated name: %s", name)
	return name, nil
}

//GetValidNameFromFolder returns a valid kubernetes name for a folder
func GetValidNameFromGitRepo(folder string) (string, error) {
	repo, err := GetRepositoryURL(folder)
	if err != nil {
		return "", err
	}
	name := translateURLToName(repo)
	return name, nil
}

func translateURLToName(repo string) string {
	repoName := findRepoName(repo)

	if strings.HasSuffix(repoName, ".git") {
		repoName = repoName[:strings.LastIndex(repoName, ".git")]
	}
	name := ValidKubeNameRegex.ReplaceAllString(repoName, "-")
	return name
}
func findRepoName(repo string) string {
	possibleName := strings.ToLower(repo[strings.LastIndex(repo, "/")+1:])
	if possibleName == "" {
		possibleName = repo
		nthTrim := strings.Count(repo, "/")
		for i := 0; i < nthTrim-1; i++ {
			possibleName = strings.ToLower(possibleName[strings.Index(possibleName, "/")+1:])
		}
		possibleName = possibleName[:len(possibleName)-1]
	}
	return possibleName
}
func GetRepositoryURL(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", fmt.Errorf("failed to analyze git repo: %w", err)
	}

	origin, err := repo.Remote("origin")
	if err != nil {
		if err != git.ErrRemoteNotFound {
			return "", fmt.Errorf("failed to get the git repo's remote configuration: %w", err)
		}
	}

	if origin != nil {
		return origin.Config().URLs[0], nil
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("failed to get git repo's remote information: %w", err)
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("git repo doesn't have any remote")
	}

	return remotes[0].Config().URLs[0], nil
}

func getDependentCyclic(s *Stack) []string {
	visited := make(map[string]bool)
	stack := make(map[string]bool)
	cycle := make([]string, 0)
	for svcName := range s.Services {
		if dfs(s, svcName, visited, stack) {
			for svc, isInStack := range stack {
				if isInStack {
					cycle = append(cycle, svc)
				}
			}
			return cycle
		}
	}
	return cycle
}

func dfs(s *Stack, svcName string, visited, stack map[string]bool) bool {
	isVisited := visited[svcName]
	if !isVisited {
		visited[svcName] = true
		stack[svcName] = true

		svc := s.Services[svcName]
		for dependentSvc := range svc.DependsOn {
			if !visited[dependentSvc] && dfs(s, dependentSvc, visited, stack) {
				return true
			} else if value, ok := stack[dependentSvc]; ok && value {
				return true
			}
		}
	}
	stack[svcName] = false
	return false
}
