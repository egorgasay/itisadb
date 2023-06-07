./lb -a=':800' -d="$MONGO" &
sleep 2
./gs -a=':1000' -connect=':800' -d="$MONGO" &
./gs -a=':1001' -connect=':800' -d="$MONGO" &
sleep 3
go test .\...