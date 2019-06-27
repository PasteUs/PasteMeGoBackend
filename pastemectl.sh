#!/usr/bin/env bash

if [[ ${#} != 1 ]]; then
    echo "Usage: pastemectl.sh <install|uninstall|upgrade|restart>"
else
    if [[ ${1} == "install" ]]; then
        echo "install" # TODO
    elif [[ ${1} == "uninstall" ]]; then
        echo "uninstall" # TODO
    else
        # TODO
        echo "[ERROR] unsupported operation: ${1}"
    fi
fi
