#!/bin/bash

# Check if the directory is provided as an argument
if [ $# -eq 0 ]; then
    echo "Usage: $0 <directory>"
    exit 1
fi

# Get the directory from the arguments
directory=$1

# Check if the provided argument is a directory
if [ ! -d "$directory" ]; then
    echo "Error: $directory is not a directory."
    exit 1
fi

# Initialize an array to hold the file paths
file_array=()

# Loop through all files in the directory and add them to the array
for file in "$directory"/*; 
do
    # Check if it's a file (and not a directory or other type)
    if [ -f "$file" ]; then
        file_array+=("$file")
    fi
done

# Run the edge_detector executable with the files as arguments
time go run . "${file_array[@]}"
