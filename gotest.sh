#!/usr/bin/env bash

BASE=github.com/PasteUs/PasteMeGoBackend/

PACKAGE_LISTS="
server
util
"

rm -f server/pasteme.db

if [[ ${#} == 1 ]]; then
    go test -count=1 -cover "${BASE}${1}"
    exit ${?}
fi

for PACKAGE in ${PACKAGE_LISTS}; do
    if ! go test -count=1 -cover "${BASE}${PACKAGE}"; then
        exit 1
    fi
done

rm -f server/pasteme.db
rm -f server/pasteme.log

echo "All test done"
