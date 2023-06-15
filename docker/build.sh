#!/bin/bash

PUBLIC_PREFIX=craftmine/pwp
BUILD_PLATFORMS=(linux/arm64 linux/amd64) #

NPM_DIR=vue-project

cd $(dirname $0)

function build(){
	tag=$1
	platform=$2
	fulltag="${PUBLIC_PREFIX}:${tag}"
	echo
	echo "==> building $fulltag from Dockerfile.$tag"
	echo
	DOCKER_BUILDKIT=1 docker build --platform ${platform} \
	 --tag "$fulltag" \
	 --file "Dockerfile.$tag" \
	 .. || return $?
	echo
	if [ -n "$TAG" ]; then
		docker tag "$fulltag" "${fulltag}-${TAG}" || return $?
		echo "==> pushing $fulltag ${fulltag}-${TAG}"
		echo
		(docker push "$fulltag" && docker push "${fulltag}-${TAG}") || return $?
	fi
	return 0
}

echo
# cur="${PWD}"
# cd "../$NPM_DIR"
# npm install || exit $?
# npm run build || exit $?
# cd "$cur"

for platform in "${BUILD_PLATFORMS[@]}"; do
	build web $platform || exit $?
	build ghupdater $platform || exit $?
	build v1 $platform || exit $?
	build reverse_proxy $platform || exit $?
done
