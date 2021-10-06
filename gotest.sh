#!/usr/bin/env bash

BASE=github.com/PasteUs/PasteMeGoBackend/

PACKAGE_LISTS="
model/paste
handler/paste
router
"

export UNITTEST=1

clear() {
    rm -f "${1}/pasteme.db"
    rm -f "${1}/pasteme.log"
}

if [[ ${#} == 1 ]]; then
    if [[ ${1} == "clear" ]]; then
        find "${PWD}" -name "*.log" -exec rm -f {} \; && find "${PWD}" -name "*.db" -exec rm -f {} \;
        exit ${?}
    fi
    clear "${1}"
    go test -count=1 -cover "${BASE}${1}"
    exit ${?}
fi

for PACKAGE in ${PACKAGE_LISTS}; do
    clear "${PACKAGE}"

    if [[ ${PACKAGE} == "util" ]]; then
        if ! go test -count=1 -cover "${BASE}${PACKAGE}"; then
            echo "test ${PACKAGE} failed"
            exit 1
        fi
    else
        if ! go test -count=1 -cover "${BASE}${PACKAGE}" -args -c "${PWD}/config.json" --debug; then
            echo "test ${PACKAGE} failed"
            exit 1
        fi
    fi

    echo "test ${PACKAGE} finished"
done

echo "All test done"
