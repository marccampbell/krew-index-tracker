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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/corneliusweig/krew-index-tracker/pkg/github"
	"github.com/corneliusweig/krew-index-tracker/pkg/globals"
	"github.com/corneliusweig/krew-index-tracker/pkg/homebrew"
	"github.com/sirupsen/logrus"
)

type requestHandler struct{}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, &requestHandler{}))
}

func (h *requestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Println("Transfer triggered")
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		logrus.Fatal("GitHub token was not set")
	}

	var failed bool
	retry := backoff.WithContext(
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), globals.DefaultRetries),
		r.Context(),
	)
	if err := backoff.RetryNotify(func() error {
		return github.SaveDownloadCountsToBigQuery(r.Context(), token, true)
	}, retry, func(err error, duration time.Duration) {
		logrus.Warnf("Failed after %s with %s", duration, err)
	}); err != nil {
		logrus.Errorf("GitHub scraping failed: %s", err)
		failed = true
	}

	if err := backoff.RetryNotify(func() error {
		return homebrew.SaveAnalyticsToBigQuery(r.Context())
	}, retry, func(err error, duration time.Duration) {
		logrus.Warnf("Failed after %s with %s", duration, err)
	}); err != nil {
		logrus.Errorf("Homebrew scraping failed: %s", err)
		failed = true
	}

	if failed {
		logrus.Infof("Scraping failed")
		w.WriteHeader(http.StatusFailedDependency)
		return
	}

	logrus.Infof("All good")
	w.WriteHeader(http.StatusOK)
}
