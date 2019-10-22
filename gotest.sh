#!/usr/bin/env bash

BASE=github.com/PasteUs/PasteMeGoBackend/

PACKAGE_LISTS="
server
util
util/generator
"

cp config.example.json config.json
rm -r server/pasteme.db

if [[ ${#} == 1 ]]; then
    go test -v -count=1 ${BASE}${1}
    exit ${?}
fi

for PACKAGE in ${PACKAGE_LISTS}; do
    go test -v -count=1 ${BASE}${PACKAGE}
    if [[ $? != 0 ]]; then
        exit 1
    fi
done

echo "All test done"
