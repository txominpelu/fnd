#!/usr/bin/env bash

# A copy of __fzf_history found at:
# https://github.com/junegunn/fzf/blob/master/README.md

bind '"\C-r": "\C-x1\e^\er"'
bind -x '"\C-x1": __fnd_history';

__fnd_history ()
{
    __ehc "$({ echo "index command"; history; } | fnd --sorter index --line_format tabular --output_column "command")"
}

__ehc()
{
if
        [[ -n $1 ]]
then
        bind '"\er": redraw-current-line'
        bind '"\e^": magic-space'
        READLINE_LINE=${READLINE_LINE:+${READLINE_LINE:0:READLINE_POINT}}${1}${READLINE_LINE:+${READLINE_LINE:READLINE_POINT}}
        READLINE_POINT=$(( READLINE_POINT + ${#1} ))
else
        bind '"\er":'
        bind '"\e^":'
fi
}
