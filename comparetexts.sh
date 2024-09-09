#!/bin/bash

# expected use:
# ./comparetexts.sh result.txt resultTest.txt

# Check if two arguments (files) are passed
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 file1.txt file2.txt"
    exit 1
fi

# Check if both files exist
if [ ! -f "$1" ]; then
    echo "File $1 does not exist."
    exit 1
fi

if [ ! -f "$2" ]; then
    echo "File $2 does not exist."
    exit 1
fi

# Compare the two files
if cmp -s "$1" "$2"; then
    echo "The files are the same."
else
    echo "The files are different."
fi
