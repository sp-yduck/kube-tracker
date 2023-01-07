#!/usr/local/bin/bash

git clone $KTRACKERREPO
pushd kube-tracker
export DATETIME=$(date '+%Y%m%d-%H%M')
git checkout -b $DATETIME
kube-tracker --config $KTRACKERCONFIG --dir $KTRACKERDIR
git config --global user.email $GITEMAIL
git config --global user.name $GITUSER
git add $KTRACKERDIR
git commit -m "$DATETIME; periodical commit"
git push origin $DATETIME
