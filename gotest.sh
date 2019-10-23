#!/usr/bin/env bash

BASE=github.com/PasteUs/PasteMeGoBackend/

PACKAGE_LISTS="
server
util
"

cp config.example.json config.json
rm -f server/pasteme.db

if [[ ${#} == 1 ]]; then
    go test -v -count=1 "${BASE}${1}"
    exit ${?}
fi

for PACKAGE in ${PACKAGE_LISTS}; do
    if ! go test -v -count=1 "${BASE}${PACKAGE}"; then
        exit 1
    fi
done

rm -f server/pasteme.db

echo "All test done"
