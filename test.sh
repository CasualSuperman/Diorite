set -e

# Run our tests
go test ./multiverse

# Compile the binaries
go build -o server.bin ./server
go build -o diorite.bin

# Start the server
./server.bin -travis &

# Wait for it to bind
sleep 0.1

# Start our client
./diorite.bin -local -noserver

# Clean up
rm diorite.bin server.bin
