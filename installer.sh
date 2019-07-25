#!/usr/bin/env bash

if [[ ${#} != 1 ]]; then
    echo "Usage: ${0} <install|uninstall|upgrade>"
else
    set -x
    if [[ ${1} == "install" ]]; then
        /usr/bin/env bash ${0} uninstall && \
        git clone --depth=1 https://github.com/LucienShui/PasteMeBackend.git -b build /usr/local/pastemed && \
        cd /usr/local/pastemed && \
        mkdir -p /etc/pastemed && \
        ln -s ${PWD}/config.sh /etc/pastemed/ && \
        ln -s ${PWD}/pastemectl.sh /usr/local/bin/pastemectl && \
        ln -s ${PWD}/pastemed.service /lib/systemd/system/ && \
        systemctl daemon-reload && \
        systemctl enable pastemed
        if [[ ${?} == 0 ]]; then
            echo -e "\033[32mInstall finished\033[0m"
            echo "Config file: /etc/pastemed/config.sh"
        else
            echo -e '\033[31mInstallation failed\033[0m'
        fi
    elif [[ ${1} == "uninstall" ]]; then
        systemctl stop pastemed
        systemctl disable pastemed
        rm -f /lib/systemd/system/pastemed.service
        systemctl daemon-reload
        rm -f /usr/local/bin/pastemectl
        rm -rf /usr/local/pastemed
        rm -rf /etc/pastemed
        echo -e "\033[32mUninstall finished\033[0m"
    elif [[ ${1} == "upgrade" ]]; then
        cd /usr/local/pastemed
        git pull
    else
        echo "[ERROR] unsupported operation: ${1}"
    fi
fi
