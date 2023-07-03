#!/usr/bin/env bash

yq -i '.fileignoreconfig |= sort_by(.filename)' .talismanrc