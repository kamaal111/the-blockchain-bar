#!/bin/sh

FILE="database/block.db"

if [ -f "$FILE" ]; then
    rm $FILE
fi

touch $FILE

go run ./cmd/tbbmigrate
