#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Search for text with rg and open a file where a match happened
fnd-rg-edit() {
    set -e
    QUERY=$1
    if [ -z "$QUERY" ]
    then
        echo >&2 "Should pass a text to search for with rg"
        return 1
    else
        CHOICE=$(rg $QUERY --line-number | python3 $DIR/rg-to-fnd.py | fnd --line_format json \
            --search_type fuzzy \
            --output_template="+{{.line}} {{.file}}" \
            --display_columns="file,line,content" \
            file)
        if [ ! -z "$CHOICE" ]
        then
            $EDITOR $(echo $CHOICE)
        fi
    fi
}
