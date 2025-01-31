#!/bin/bash

encrypt() {
    local x=$1
    local y=$2
    echo $(( (x + y) % 3 ))
}

decrypt() {
    local z=$1
    local y=$2
    echo $(( (z - y + 3) % 3 ))
}

if [ "$1" = "-d" ]; then
    DECRYPT=true
    KEYFILE=$2
else
    DECRYPT=false
    KEYFILE=$1
fi

if [ -z "$KEYFILE" ]; then
    echo "Usage: $0 [-d] key.txt < infile > outfile"
    exit 1
fi

# Read and trim key file
key=$(tr -d '\r\n' < "$KEYFILE")

# Read and trim message from stdin
message=$(tr -d '\r\n')

if [ ${#message} -ne ${#key} ]; then
    echo "Error: Message and key must have the same length"
    exit 1
fi

result=""
for (( i=0; i<${#message}; i++ )); do
    x=${message:$i:1}
    y=${key:$i:1}

    if [[ $x =~ ^[0-2]$ && $y =~ ^[0-2]$ ]]; then
        if [ "$DECRYPT" = true ]; then
            result+=$(decrypt $x $y)
        else
            result+=$(encrypt $x $y)
        fi
    else
        echo "Error: Invalid character in message or key"
        exit 1
    fi
done

echo "$result"
