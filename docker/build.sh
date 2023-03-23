#!/bin/bash

PUBLIC_PREFIX=craftmine/pwp
platform=linux/amd64

NPM_DIR=vue-project

cd $(dirname $0)

function build(){
	tag=$1
	fulltag="${PUBLIC_PREFIX}:${tag}"
	echo
	echo "==> building $fulltag from Dockerfile.$tag"
	echo
	docker build --platform ${platform} \
	 --tag "$fulltag" \
	 --file "Dockerfile.$tag" \
	 .. || return $?
	echo
	echo "==> pushing $fulltag"
	echo
	docker push "$fulltag" || return $?
	return 0
}

echo
cur="${PWD}"
cd "../$NPM_DIR"
npm run build || exit $?
cd "$cur"
build web || exit $?
build v1 || exit $?
