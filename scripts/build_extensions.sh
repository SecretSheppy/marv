#!/bin/bash

for dir in */; do
    dirname=$(basename "$dir")
    go_file="$dir$dirname.go"

    if [[ -f "$go_file" ]]; then
        go build -buildmode=plugin -o "$dir$dirname.so" "$go_file"

        if [[ $? -eq 0 ]]; then
            echo "Compiled $go_file to $dirname.so"
        else
            echo "Failed to compile $go_file"
        fi
    else
        echo "Go file $go_file does not exist"
    fi
done