#!/usr/bin/env bash

name=$2
digest=$4
export HOME=/root
ratifyOutput=$(~/.ratify/bin/ratify verify -c ~/.ratify/config.json -s $name --digest $digest)
isSuccess=$(echo $ratifyOutput | jq .isSuccess)
echo $ratifyOutput | jq . 
if [ "$isSuccess" = "true" ]; then
    echo "ratify verification succeeded"
    exit 0
fi
exit 1