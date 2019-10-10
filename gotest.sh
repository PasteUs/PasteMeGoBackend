#!/usr/bin/env bash

BASE=github.com/PasteUs/PasteMeGoBackend/

PACKAGE_LISTS="
server
util
util/generator
"

PASTEMED_TEST_ENVS="
PASTEMED_DB_USERNAME=username
PASTEMED_DB_PASSWORD=password
PASTEMED_DB_SERVER=mysql
PASTEMED_DB_PORT=3306
PASTEMED_DB_DATABASE=pasteme
PASTEMED_RUNTIME=debug
"

for PASTEMED_TEST_ENV in ${PASTEMED_TEST_ENVS}; do
    export ${PASTEMED_TEST_ENV}
done

if [[ ${#} == 1 ]]; then
    go test -v -count=1 ${BASE}${1}
    exit ${?}
fi

for PACKAGE in ${PACKAGE_LISTS}; do
    go test -v -count=1 ${BASE}${PACKAGE}
    if [[ ${?} != 0 ]]; then
        exit 1
    fi
done

echo "All test done"
