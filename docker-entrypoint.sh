#!/usr/bin/env sh

CONFIG_PATH="/etc/pastemed/config.json"
EXAMPLE_CONFIG_PATH="/usr/local/pastemed/config.example.json"
DATA_PATH="/data"

# 如果配置文件不存在，就使用默认配置
if [ ! -f "${CONFIG_PATH}" ]; then
  cp ${EXAMPLE_CONFIG_PATH} ${CONFIG_PATH}
fi

/usr/local/pastemed/pastemed -c ${CONFIG_PATH} -d ${DATA_PATH}
