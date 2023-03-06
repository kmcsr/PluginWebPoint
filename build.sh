#!/bin/bash

WEB_ONLY=
RUN=
DEBUG="${DEBUG}"

function build_web(){
	_subdir=vue-project
	cd "$_subdir"
	if [[ "$DEBUG" == true ]]; then
		npm run build || return $?
	else
		npm run build_dev || return $?
	fi

	cd ..
	rm -rf "./dist"
	cp -a "${_subdir}/dist" ./dist
	return $?
}

function build_app(){
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./output/plugin_web_point .
	return $?
}

function run_app(){
	go run . "$@"
	return $?
}

while [ -n "$1" ]; do
	case $1 in
		-w | --web-only)
			WEB_ONLY=true
			;;
		-r | --run)
			RUN=true
			;;
		-d | --debug)
			DEBUG=true
			;;
	esac
	shift
done

cd $(dirname $0)

build_web || exit $?

if [[ "$WEB_ONLY" != true ]]; then
	if [[ "$RUN" == true ]]; then
		if [[ "$DEBUG" == true ]]; then
			run_app -debug || exit $?
		else
			run_app || exit $?
		fi
	else
		build_app || exit $?
	fi
fi
