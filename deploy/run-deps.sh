#!/bin/bash

docker compose stop alloy || true
yes | docker compose rm alloy -v || true

rm -rf .run
mkdir -p .run
touch .run/shorty.log
chmod 777 .run/shorty.log

docker compose up -d