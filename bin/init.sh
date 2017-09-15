#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -x

# Ensure expected GOPATH setup
PDIR=`pwd`
if [ $PDIR != "${GOPATH-$HOME/go}/src/istio.io/pilot" ]; then
       echo "Pilot not found in GOPATH/src/istio.io/"
       exit 1
fi

# Building and testing with Bazel
bazel build //...

source "${PDIR}/bin/use_bazel_go.sh"
go version

# Clean up vendor dir
rm -rf $(pwd)/vendor

# Vendorize bazel dependencies
bin/bazel_to_go.py

# Remove doubly-vendorized k8s dependencies
rm -rf vendor/k8s.io/*/vendor

# Link proto gen files
mkdir -p vendor/istio.io/api/proxy/v1/config
for f in dest_policy.pb.go  http_fault.pb.go  l4_fault.pb.go  proxy_mesh.pb.go  route_rule.pb.go ingress_rule.pb.go egress_rule.pb.go; do
  ln -sf $(pwd)/bazel-genfiles/external/io_istio_api/proxy/v1/config/$f \
    vendor/istio.io/api/proxy/v1/config/
done

mkdir -p vendor/istio.io/pilot/test/grpcecho
for f in $(pwd)/bazel-genfiles/test/grpcecho/*.pb.go; do
  ln -sf $f vendor/istio.io/pilot/test/grpcecho/
done

# Link envoy binary
ln -sf "$(pwd)/bazel-genfiles/proxy/envoy/envoy" proxy/envoy/

# Link CRD generated files
ln -sf "$(pwd)/bazel-genfiles/adapter/config/crd/types.go" \
  adapter/config/crd/

# Some linters expect the code to be installed
go install ./...
