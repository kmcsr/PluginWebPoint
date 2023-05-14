#!/bin/sh

program=$1
interval=$2

[ -n "$program" ] || { echo "ERR: must give a program path"; exit 1; }
[ -n "$interval" ] || interval=3600

if ! git --version; then
	echo 'WARN: Git not found'
	echo 'Installing git...'
	apk add git || {
		echo "ERR: Cannot install git"
		exit 1
	}
fi

while true; do
	echo
	echo "Running: $program"
	echo
	"$program" >> "${program}.log" 2>&1
	echo 'Exit status:' $?
	echo '==> Done'
	echo
	sleep "$interval"
done
