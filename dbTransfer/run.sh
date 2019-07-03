#!/usr/bin/env bash

# Database config
PASTEMED_DB_ENVS="
PASTEMED_DB_USERNAME=username
PASTEMED_DB_PASSWORD=password
PASTEMED_DB_SERVER=mysql
PASTEMED_DB_PORT=3306
PASTEMED_DB_DATABASE=pasteme
"

for ENV in ${PASTEMED_DB_ENVS}; do
    export ${ENV}
done

./db_transfer
