/*
Copyright 2019 Cornelius Weig.

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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/corneliusweig/krew-index-tracker/pkg/github/util"
	"github.com/corneliusweig/krew-index-tracker/pkg/homebrew"
)

var (
	token         string
	isUpdateIndex bool
)

var rootCmd = &cobra.Command{
	Use:     "krew-index-tracker",
	Example: "krew-index-tracker",
	Short:   "Upload krew repo statistics to BigQuery",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := util.ContextWithCtrlCHandler(context.Background())

		if err := github.SaveDownloadCountsToBigQuery(ctx, token, isUpdateIndex); err != nil {
			return err
		}

		return homebrew.SaveAnalyticsToBigQuery(ctx)
	},
}

func main() {
	rootCmd.Flags().StringVar(&token, "token", "", "Specify personal Github Token if you are hitting a rate limit anonymously. https://github.com/settings/tokens")
	rootCmd.Flags().BoolVar(&isUpdateIndex, "update-index", false, "Call git to ensure that the index is up to date")
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
