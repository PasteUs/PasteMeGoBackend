#!/usr/bin/env bash

if [[ ${#} != 1 ]]; then
    echo "Usage: ${0} <start|stop|restart|status|log>"
else
    if [[ ${1} == "start" || ${1} == "stop" || ${1} == "restart" ]]; then
        systemctl ${1} pastemed
    elif [[ ${1} == "status" ]]; then
        systemctl status pastemed | grep Active
    elif [[ ${1} == "log" ]]; then
        journalctl -e -u pastemed -o cat | cat
    else
        echo "[ERROR] unsupported operation: ${1}"
        exit 1
    fi
fi
