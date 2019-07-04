#!/usr/bin/env bash

if [[ ${#} != 1 ]]; then
    echo "Usage: ${0} <install|uninstall|upgrade>"
else
    if [[ ${1} == "install" ]]; then
        set -x
        /usr/bin/env bash ${0} uninstall && \
        git clone --depth=1 https://github.com/LucienShui/PasteMeBackend.git -b build /usr/local/pastemed && \
        mkdir -p /etc/pastemed && \
        cd /usr/local/pastemed && \
        mv pastemed.service config.sh /etc/pastemed/ && \
        ln -s ${PWD}/pastemectl.sh /usr/local/bin/pastemectl && \
        chmod +x /usr/local/bin/pastemectl && \
        ln -s /etc/pastemed/pastemed.service /lib/systemd/system/ && \
        systemctl daemon-reload
        if [[ ${?} == 0 ]]; then
            echo "Installation finished"
            echo "Config file: /etc/pastemed/config.sh"
        else
            echo "Installation failed"
        fi
    elif [[ ${1} == "uninstall" ]]; then
        set -x
        rm /lib/systemd/system/pastemed.service
        rm /usr/local/bin/pastemectl
        rm -rf /usr/local/pastemed
        rm -rf /etc/pastemed
        set +x
        echo "uninstall finished"
    elif [[ ${1} == "upgrade" ]]; then
        cd /usr/local/pastemed
        git pull
    else
        echo "[ERROR] unsupported operation: ${1}"
    fi
fi
