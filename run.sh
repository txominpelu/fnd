#!/usr/bin/env bash

guid=$(python run.py list  | fzf --preview 'python run.py description {1}' --preview-window hidden --bind ctrl-k:preview-up,ctrl-j:preview-down,?:toggle-preview | cut -d '.' -f 1)
xdg-open "https://remoteok.io/l/$guid"
