#!/usr/bin/env bash

fnd-ps-aux() {
    ps aux | fnd --line_format tabular --output_column PID
}
