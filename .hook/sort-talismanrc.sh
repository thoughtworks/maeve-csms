#!/usr/bin/env bash

yq -y -i '.fileignoreconfig |= sort_by(.filename)' .talismanrc
