#!/usr/bin/env bash

DEP="
github.com/gin-gonic/gin
github.com/jinzhu/gorm
github.com/go-sql-driver/mysql
github.com/wonderivan/logger
"

for each in ${DEP}; do
    /usr/bin/env go get -u ${each}
done
