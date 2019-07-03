#!/usr/bin/env bash

time for (( i = 1; i <= 10; i++ )); do
    curl -s -X POST -H Content-Type:application/json \
         -d '{"content":"Hello","lang":"plain","password":"123"}' \
         localhost:8000
done
