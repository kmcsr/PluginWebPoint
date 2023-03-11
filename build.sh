#!/bin/bash

WEB_ONLY=
RUN=
DEBUG="${DEBUG}"
DEV="${DEV}"
NPM_DIR=vue-project

function build_web(){
	cd "$NPM_DIR"
	if [[ "$DEBUG" == true ]]; then
		NODE_ENV=development npm run build_dev || return $?
	else
		npm run build || return $?
	fi

	cd ..
	rm -rf "./dist"
	cp -a "${NPM_DIR}/dist" ./dist
	return $?
}

function build_watch(){
	cd "$NPM_DIR"
	NODE_ENV=development exec npm run build_watch
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
		-D | --dev)
			DEV=true
			;;
	esac
	shift
done

export DEBUG

cd $(dirname $0)

if [[ "$DEV" == true ]]; then
	build_watch &
	pid=$!
	go run ./dev/
	_e=$?
	kill -s SIGINT "$pid"
	exit $_e
fi

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
