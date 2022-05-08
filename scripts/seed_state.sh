#!/bin/sh

FILE="database/transactions.db"

if [ -f "$FILE" ]; then
    rm $FILE
fi

touch $FILE

add_transaction() {
    go run . tx add --from="$1" --to="$2" --value="$3" --data="$4"
}

# day 1
add_transaction andrej andrej 3
add_transaction andrej andrej 700 reward
add_transaction andrej babayaga 2000
add_transaction andrej andrej 100 reward
add_transaction babayaga andrej 1

# day 2
add_transaction babayaga andcaesarrej 1000
add_transaction babayaga andrej 50
add_transaction andrej andrej 100 reward

# day 3
add_transaction andrej andrej 100 reward
