#!/bin/sh

# Example usage:
#
# docker/release-docker --hub docker.io/istio --tags $(git rev-parse --short # HEAD),$(date +%Y%m%d%H%M%S),latest"

function usage() {
    echo "$0 --hub <docker image repository> --tags <comma seperated list of docker image tags>"
    exit 1
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --hub) hub="$2"; shift ;;
        --tags) tags="$2"; shift ;;
        *) usage ;;
    esac
    shift
done

[[ -z $hub ]] && usage
[[ -z $tags ]] && usage

tags=$(echo $tags | sed -e 's/,/ /g')

# TODO expose list of images as command line flag?
images="init init_debug app app_debug runtime runtime_debug"

if [[ "$hub" =~ ^gcr\.io ]]; then
    gcloud docker --authorize-only
fi

set -ex

for image in $images; do
    bazel $BAZEL_ARGS run //docker:$image
    for tag in $tags; do
        docker tag istio/docker:$image $hub/$image:$tag
        docker push $hub/$image:$tag
    done
done
