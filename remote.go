package main

import (
	"bufio"
	"io"
	"net"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

const serverPort = ":5050"
const remoteDBServer = "diorite.casualsuperman.com"
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

func connectToLocalServer() (serverConnection, error) {
	return connectToServer("localhost")
}

func connectToDefaultServer() (serverConnection, error) {
	return connectToServer(remoteDBServer)
}

func connectToServer(addr string) (serverConnection, error) {
	conn, err := net.Dial("tcp", addr+serverPort)
	return serverConnection{conn}, err
}

func (s serverConnection) DownloadMultiverse(saveTo io.Writer) (mv m.Multiverse, err error) {
	var conn io.Reader

	if saveTo != nil {
		conn = io.TeeReader(s, saveTo)
	} else {
		conn = s
	}

	s.Write([]byte("multiverseDL\n"))

	return m.Read(conn)
}

func (s serverConnection) Close() {
	s.Write([]byte("close\n"))
	s.Conn.Close()
}
