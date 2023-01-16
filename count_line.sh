#!/usr/bin/env sh

find .  -name  "*.go" -exec  cat {} + | grep -v ^$ -c #count
find .  -name  "*.go" -exec  cat {} + | grep -v ^$ #echo

