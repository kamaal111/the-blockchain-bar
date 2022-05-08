#!/bin/sh

# {"from":"andrej","to":"andrej","value":3,"data":""}
# {"from":"andrej","to":"andrej","value":700,"data":"reward"}
# {"from":"andrej","to":"babayaga","value":2000,"data":""}
# {"from":"andrej","to":"andrej","value":100,"data":"reward"}
# {"from":"babayaga","to":"andrej","value":1,"data":""}

FILE="database/transactions.db"

if [ -f "$FILE" ]; then
    rm $FILE
fi

touch $FILE

go run . tx add --from=andrej --to=andrej --value=3
go run . tx add --from=andrej --to=andrej --value=700 --data=reward
go run . tx add --from=andrej --to=babayaga --value=2000
go run . tx add --from=andrej --to=andrej --value=100 --data=reward
go run . tx add --from=babayaga --to=andrej --value=1
