package main

import (
	"bufio"
	"io"
	"net"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

const remoteDBLocation = "diorite.casualsuperman.com:5050"
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

func downloadMultiverse(saveTo io.Writer) (mv m.Multiverse, err error) {
	var conn io.Reader
	netConn, err := net.Dial("tcp", remoteDBLocation)
	conn = netConn

	if err != nil {
		return
	}

	defer netConn.Close()

	netConn.Write([]byte("multiverseDL\n"))

	if saveTo != nil {
		conn = io.TeeReader(conn, saveTo)
	}

	return m.Read(conn)
}
