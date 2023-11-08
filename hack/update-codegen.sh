#!/usr/bin/env bash

# Copyright 2018 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
pushd "${SCRIPT_ROOT}"
SCRIPT_ROOT=$(pwd)
popd

# Note that we use code-generator from `${GOPATH}/pkg/mod/` because we cannot vendor it
# via `go mod vendor` to the project's /vendor directory.
# Reference: https://github.com/kubernetes/code-generator/issues/57
CODEGEN_VERSION=$(grep 'k8s.io/code-generator' go.sum | awk '{print $2}' | sed 's/\/go.mod//g' | head -1)
CODEGEN_PKG=$(echo `go env GOPATH`"/pkg/mod/k8s.io/code-generator@${CODEGEN_VERSION}")
chmod +x "${CODEGEN_PKG}/generate-groups.sh"

echo "${CODEGEN_PKG}/generate-groups.sh"

APIS_PKG=api
GROUP=kubeflow
VERSION=v1beta1
DOMAIN=org
REPO=github.com/kuizhiqing/trainingjob-operator

rm -rf client ${DOMAIN}
rm -rf "${APIS_PKG}/${GROUP}" && mkdir -p "${APIS_PKG}/${GROUP}" && cp -r "${APIS_PKG}/${VERSION}/" "${APIS_PKG}/${GROUP}"

"${CODEGEN_PKG}/generate-groups.sh" "client,informer,lister" \
    ${REPO}/client ${REPO}/${APIS_PKG} \
    ${GROUP}:${VERSION} --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

rm -rf api/${GROUP}

mv ${REPO}/client . && rm -rf ${DOMAIN}

find client/ -name "*.go" | xargs sed -i "s/${APIS_PKG}\/${GROUP}\/${VERSION}/${APIS_PKG}\/${VERSION}/"
