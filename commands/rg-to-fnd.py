#!/usr/bin/env python3
import fileinput
import json
import sys

if __name__ == "__main__":
    for l in fileinput.input():
        splits = l.split(":")
        filename = splits[0]
        linenumb = int(splits[1])
        content = ":".join(splits[2:])
        d = {"file": filename, "line": linenumb, "content": content}
        print(json.dumps(d))


