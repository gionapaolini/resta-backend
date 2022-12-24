#!/bin/bash
SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")
ROOT=$SCRIPTPATH/..
cd $ROOT

find -name go.mod -printf "%h\n" | sort -nr | while read -r i; do cd "$i"; go test -count=1 -cover ./...; cd $ROOT; done