#!/usr/bin/env bash

if [[ ${#} != 1 ]]; then
    echo "Usage: ${0} <install|uninstall|upgrade>"
else
    if [[ ${1} == "install" ]]; then
        set -x
        mkdir -p /usr/local/pastemed /etc/pastemed && \
        cp pastemed db_transfer pastemectl.sh installer.sh /usr/local/pastemed/ && \
        cp pastemed.service config /etc/pastemed/ && \
        chmod +x /usr/local/pastemed/pastemed && \
        cd /usr/local/pastemed/ && \
        ln -s /usr/local/pastemed/pastemectl.sh /usr/local/bin/pastemectl && \
        chmod +x /usr/local/bin/pastemectl && \
        ln -s /etc/pastemed/pastemed.service /lib/systemd/system/ && \
        systemctl daemon-reload
        set +x
        if [[ ${?} != 0 ]]; then
            echo "Installation finished"
            echo "Config file: /etc/pastemed/config"
        else
            echo "Installation failed"
        fi
    elif [[ ${1} == "uninstall" ]]; then
        exit 0 # TODO
    elif [[ ${1} == "upgrade" ]]; then
        exit 0 # TODO
    else
        echo "[ERROR] unsupported operation: ${1}"
    fi
fi
