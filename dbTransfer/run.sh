#!/usr/bin/env bash

# Old database config
export PASTEMED_OLD_DB_USERNAME=pasteme_cn
export PASTEMED_OLD_DB_PASSWORD=password
export PASTEMED_OLD_DB_SERVER=web
export PASTEMED_OLD_DB_PORT=3306
export PASTEMED_OLD_DB_DATABASE=pasteme_cn

# New database config
export PASTEMED_DB_USERNAME=username
export PASTEMED_DB_PASSWORD=password
export PASTEMED_DB_SERVER=mysql
export PASTEMED_DB_PORT=3306
export PASTEMED_DB_DATABASE=pasteme

./db_transfer
