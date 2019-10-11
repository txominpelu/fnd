#!/usr/bin/env bash

fnd-apt-install() {
    QUERY=$1
    if [ -z "$QUERY" ]
    then
        echo >&2 "Should pass a text to search for a package"
        return 1
    else
        CHOICE=$({ echo "pkg - description"; apt-cache search $QUERY; } | fnd --line_format tabular --output_column pkg)
        sudo apt install $CHOICE
    fi
}
