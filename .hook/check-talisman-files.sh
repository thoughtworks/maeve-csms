#!/usr/bin/env bash

files=$(for f in $(yq '.fileignoreconfig.[].filename' .talismanrc); do if [ ! -f "$f" ]; then echo "$f"; fi; done)
count=$(echo "$files" | sed '/^$/d' | wc -l)

if [[ "$count" -gt 0 ]];
then
  echo "The following files in .talismanrc do not exist:"
  echo "$files"
  exit 1
fi