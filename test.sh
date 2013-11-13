go test ./multiverse
go build -o server.bin ./server
go build -o diorite.bin

./server.bin -travis &
./diorite.bin -local
