#!/bin/bash

set -ex

# Test harness will checkout code to directory $GOPATH/src/github.com/istio/pilot
# but we depend on being at path $GOPATH/src/istio.io/pilot for imports
if [ "$CI" == "bootstrap" ]; then
    mkdir -p $GOPATH/src/istio.io
    mv $GOPATH/src/github.com/istio/pilot $GOPATH/src/istio.io
    cd $GOPATH/src/istio.io/pilot/

    # use the provided pull head sha
    GIT_SHA=$PULL_PULL_SHA
else
    # use the current commit
    GIT_SHA=$(git rev-parse --verify HEAD)
fi

# Get configuration for the separate test cluster, it must be at
# ~/.kube and platform/kube because different aspects of testing
# & building require it in each place.
gcloud config set container/use_client_certificate True
gcloud container clusters get-credentials e2e-testing --zone us-west1-a --project istio-testing
if [ -e platform/kube/config ]; then
    rm platform/kube/config
fi
ln -s ~/.kube/config platform/kube/

echo "=== Bazel Build ==="
./bin/install-prereqs.sh
bazel fetch -k //...
bazel build //...

echo "=== Go Build ==="
./bin/init.sh

echo "=== Code Check ==="
./bin/check.sh

echo "=== Bazel Tests ==="
bazel test //...

echo "=== Integration Tests ==="
./bin/e2e.sh -tag $GIT_SHA -hub "gcr.io/istio-testing"
