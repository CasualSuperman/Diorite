package main

import (
	"bufio"
	"net"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

const remoteDBLocation = "localhost:5050"
const lastModifiedFormat = time.RFC1123

func onlineModifiedAt() (time.Time, error) {
	conn, err := net.Dial("tcp", remoteDBLocation)

	if err != nil {
		return time.Time{}, err
	}

	conn.Write([]byte("multiverseMod\n"))

	s := bufio.NewScanner(conn)
	s.Scan()
	t := s.Text()
	conn.Close()

	return time.Parse(lastModifiedFormat, t)
}

func downloadMultiverse() (mv m.Multiverse, err error) {
	conn, err := net.Dial("tcp", remoteDBLocation)

	if err != nil {
		return
	}

	conn.Write([]byte("multiverseDL\n"))
	defer conn.Close()

	return m.Read(conn)
}
