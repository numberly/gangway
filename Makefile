# Copyright © 2017 Heptio
# Copyright © 2017 Craig Tracey
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

PROJECT := gangway
# Where to push the docker image.
REGISTRY ?= numberly
IMAGE := $(REGISTRY)/$(PROJECT)
SRCDIRS := ./cmd/gangway

VERSION ?= master

all: build

build: deps
	go build ./...

install:
	go install -v ./cmd/gangway/...

setup:
	curl -s -o assets/prism-core.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-core.min.js
	curl -s -o assets/prism-bash.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-bash.min.js
	curl -s -o assets/prism-yaml.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-yaml.min.js
	curl -s -o assets/prism-powershell.min.js https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/components/prism-powershell.min.js
	curl -s -o assets/prism-tomorrow.min.css https://raw.githubusercontent.com/PrismJS/prism/v1.28.0/themes/prism-tomorrow.min.css
	curl -s -o assets/materialize.min.css https://raw.githubusercontent.com/Dogfalo/materialize/v1-dev/dist/css/materialize.min.css
	curl -s -o assets/materialize.min.js https://raw.githubusercontent.com/Dogfalo/materialize/v1-dev/dist/js/materialize.min.js

check: test lint gofmt misspell

deps:
	go mod tidy && go mod vendor && go mod verify

test:
	go test -v ./...

misspell:
	@go get github.com/client9/misspell/cmd/misspell
	misspell \
		-i clas \
		-locale US \
		-error \
		cmd/* docs/* *.md

gofmt:
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"

image:
	docker build . -t $(IMAGE):$(VERSION)

push:
	docker push $(IMAGE):$(VERSION)

.PHONY: all deps test image
