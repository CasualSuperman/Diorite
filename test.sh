go test ./multiverse
go build ./server
go build

./server/server -travis &
./Diorite -local
