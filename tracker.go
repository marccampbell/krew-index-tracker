/*
Copyright 2019 Cornelius Weig

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/corneliusweig/krew-index-tracker/pkg/git"
	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/pkg/errors"
	"sigs.k8s.io/krew/pkg/index"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/krew/pkg/index/indexscanner"
)

const (
	indexDir   = "index"
	pluginsDir = indexDir + "/plugins"
)

var (
	token         string
	isUpdateIndex bool
)

type pluginHandle struct {
	index.PluginSpec
	owner, repo string
}

var rootCmd = &cobra.Command{
	Use:     "krew-index-tracker",
	Example: "krew-index-tracker",
	Short:   "Generate a markdown changelog of merged pull requests since last release",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := git.UpdateAndCleanUntracked(isUpdateIndex, indexDir); err != nil {
			logrus.Fatal(err)
		}

		repos, err := getRepoList()
		if err != nil {
			logrus.Fatal(err)
		}

		ctx := contextWithCtrlCHandler(context.Background())
		releaseFetcher := github.NewReleaseFetcher(ctx, token)
		summary, err := releaseFetcher.RepoSummary(repos[0].owner, repos[0].repo)
		marshal, _ := json.Marshal(summary)
		fmt.Printf("%s", string(marshal))
	},
}

func main() {
	rootCmd.Flags().StringVar(&token, "token", "", "Specify personal Github Token if you are hitting a rate limit anonymously. https://github.com/settings/tokens")
	rootCmd.Flags().BoolVar(&isUpdateIndex, "update-index", false, "Call git to ensure that the index is up to date")
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func getRepoList() ([]pluginHandle, error) {
	plugins, err := indexscanner.LoadPluginListFromFS(pluginsDir)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read index")
	}
	gitHubRepo := regexp.MustCompile(".*github.com/([^/]+)/([^/]+).*")
	res := make([]pluginHandle, 0, len(plugins))
	for _, plugin := range plugins {
		submatch := gitHubRepo.FindStringSubmatch(plugin.Spec.Homepage)
		if len(submatch) < 3 {
			continue
		}
		res = append(res, pluginHandle{
			PluginSpec: plugins[0].Spec,
			owner:      submatch[1],
			repo:       submatch[2],
		})
	}
	return res, nil
}

func contextWithCtrlCHandler(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGPIPE)

	go func() {
		<-sigs
		signal.Stop(sigs)
		cancel()
		logrus.Infof("Aborted.")
	}()

	return ctx
}
