#!/usr/bin/env bash

go test -count=1 -cover ./... -args -c "${PWD}/config.json"

echo "All test done"
