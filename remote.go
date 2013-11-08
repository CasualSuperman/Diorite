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

type serverConnection struct {
	net.Conn
}

func (s serverConnection) Modified() time.Time {
	s.Write([]byte("multiverseMod\n"))

	scan := bufio.NewScanner(s)
	scan.Scan()
	t := scan.Text()

	tim, _ := time.Parse(lastModifiedFormat, t)

	return tim
}

func connectToServer() (serverConnection, error) {
	conn, err := net.Dial("tcp", remoteDBLocation)
	return serverConnection{conn}, err
}

func (s serverConnection) DownloadMultiverse(saveTo io.Writer) (mv m.Multiverse, err error) {
	var conn io.Reader = s

	if saveTo != nil {
		conn = io.TeeReader(s, saveTo)
	}

	s.Write([]byte("multiverseDL\n"))

	return m.Read(conn)
}

func (s serverConnection) Close() {
	s.Write([]byte("close\n"))
	s.Conn.Close()
}
