#!/usr/bin/env bash

fnd-fdfind() {
    set -e
    echo $(fdfind "$@" | fnd)
}
