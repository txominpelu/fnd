#!/usr/bin/env bash

python run.py list  | fzf --preview 'python run.py description {1}' --preview-window hidden --bind ctrl-k:preview-up,ctrl-j:preview-down,?:toggle-preview
