#!/usr/bin/env bash

fnd-kill() {
    CHOICE=$(ps aux | fnd --line_format tabular --output_column PID)
    kill -9 $CHOICE
}
