#!/bin/bash

# Run golangci-lint and store the output
output=$(golangci-lint run)

# Convert relative paths to absolute paths for lines matching "FILENAME:LINE:COLUMN:"
# The regular expression looks for patterns like "filename.go:123:45:"
echo "$output" | sed -E "s|^([a-zA-Z0-9_./-]+\.go:[0-9]+:[0-9]+:)|$(pwd)/\1|"

