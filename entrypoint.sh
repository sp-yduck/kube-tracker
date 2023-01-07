#!/usr/local/bin/bash

git clone github.com/sp-yduck/kube-tracker
pushd kube-tracker
export DATETIME=$(date '+%Y%m%d-%H:%M')
git checkout -b $DATETIME
kube-tracker --config KTRACKERCONFIG --dir KTRACKERDIR
git add .
git commit -m "$DATETIME; periodical commit"
git push origin $DATETIME
