#!/bin/bash

# Function for encryption: (x + y) mod 3
encrypt() {
    local x=$1
    local y=$2
    echo $(( (x + y) % 3 ))
}

# Function for decryption: (z - y) mod 3
decrypt() {
    local z=$1
    local y=$2
    echo $(( (z - y + 3) % 3 ))  # +3 to avoid negative results
}

# Check if key filename was provided
if [ -z "$1" ]; then
    echo "Usage: $0 key.txt < infile > outfile"
    exit 1
fi

# Read key file
key=$(cat "$1")

# Read message from stdin
message=$(cat)

# Check if message and key have the same length
if [ ${#message} -ne ${#key} ]; then
    echo "Error: Message and key must have the same length"
    exit 1
fi

# Decrypt
decrypted_result=""
for (( i=0; i<${#message}; i++ )); do
    z=${message:$i:1}
    y=${key:$i:1}

    # Only process valid characters (0, 1, 2)
    if [[ $z =~ ^[0-2]$ && $y =~ ^[0-2]$ ]]; then
        decrypted_result+=$(decrypt $z $y)
    else
        echo "Error: Invalid character in message or key"
        exit 1
    fi
done

# Output result
echo "$decrypted_result"
