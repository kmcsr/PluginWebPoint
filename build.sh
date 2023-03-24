#!/bin/bash

WEB_ONLY=
RUN=
DEBUG="${DEBUG}"
DEV="${DEV}"
NPM_DIR=vue-project

function build_web(){
	_old_pwd="${PWD}"
	cd "$NPM_DIR"
	if [[ "$DEBUG" == true ]]; then
		NODE_ENV=development npm run build_dev || return $?
	else
		npm run build || return $?
	fi
	cd "${_old_pwd}"
	return $?
}

function build_watch(){
	_old_pwd="${PWD}"
	cd "$NPM_DIR"
	NODE_ENV=development exec npm run build_watch || return $?
	cd "${_old_pwd}"
	return $?
}

function build_app(){
	rm -rf "./cmds/plugin_web_point/dist"
	cp -a "${NPM_DIR}/dist" "./cmds/plugin_web_point/dist"

	CGO_ENABLED=0 go build -o ./output/plugin_web_point ./cmds/plugin_web_point
	return $?
}

function build_handle(){
	_handle=$1
	echo "==> Building handle ${_handle}"
	CGO_ENABLED=0 go build -o "./output/pwp_${_handle}" "./handlers/${_handle}"
	return $?
}

function run_app(){
	go run ./cmds/dev "$@"
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
	go run ./cmds/dev/
	_e=$?
	kill -s SIGINT "$pid"
	exit $_e
fi

build_web || exit $?

if [[ "$WEB_ONLY" != true ]]; then
	if [[ "$RUN" == true ]]; then
		if [[ "$DEBUG" == true ]]; then
			run_app -debug
		else
			run_app
		fi
		exit $?
	fi
	echo '==> Building app'
	cp ./robots.txt ./output/robots.txt
	GOARCH=amd64 GOOS=linux build_app || exit $?
	echo '==> Building ghupdater'
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./output/ghupdater ./cmds/ghupdater || exit $?
	GOARCH=amd64 GOOS=linux build_handle dev || exit $?
	GOARCH=amd64 GOOS=linux build_handle v1 || exit $?
	echo '==> Done'
fi
