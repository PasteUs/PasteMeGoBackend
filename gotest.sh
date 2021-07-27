#!/usr/bin/env bash

rm -f server/pasteme.db

UNITTEST=1 go test -count=1 -cover ./...

rm -f server/pasteme.db
rm -f server/pasteme.log

echo "All test done"
