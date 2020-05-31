# Copyright 2019 Cornelius Weig.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SRC:=$(shell find . -name '*.go')

krew-index-tracker: $(SRC)
	go build -trimpath -tags netgo -ldflags "-s -w -extldflags '-static'" -o $@ ./app/krew-index-tracker

krew-index-tracker-http: $(SRC)
	go build -trimpath -tags netgo -ldflags "-s -w -extldflags '-static'" -o $@ ./app/http

.PHONY: lint
lint: $(SRC)
	hack/run-lint.sh

.PHONY: test
test: $(SRC)
	hack/verify-boilerplate.sh && \
	go test ./...

.PHONY: build
build:
	gcloud builds submit --config cloudbuild.yaml --substitutions=TAG_NAME="$$(git describe --tags --always)" .
