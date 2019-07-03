#!/usr/bin/env bash

# Database config
export PASTEMED_DB_USERNAME=username
export PASTEMED_DB_PASSWORD=password
export PASTEMED_DB_SERVER=mysql
export PASTEMED_DB_PORT=3306
export PASTEMED_DB_DATABASE=pasteme

./db_transfer
