package main

import (
	"bytes"
	"encoding/gob"
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

	t, _ := time.Parse(lastModifiedFormat, s.readLine())

	return t
}

func connectToLocalServer() (serverConnection, error) {
	return connectToGivenServer("localhost")
}

func connectToDefaultServer() (serverConnection, error) {
	return connectToGivenServer(remoteDBServer)
}

func connectToGivenServer(addr string) (serverConnection, error) {
	conn, err := net.Dial("tcp", addr+serverPort)
	return serverConnection{conn}, err
}

func (s serverConnection) RawMultiverse() []byte {
	var l int32

	s.Write([]byte("multiverseLen\n"))

	dec := gob.NewDecoder(s)
	dec.Decode(&l)
	data := make([]byte, int(l))

	s.Write([]byte("multiverseDL\n"))
	num, err := s.Read(data)
	total := num

	for total < int(l) {
		dataSegment := data[total:]
		num, err := s.Read(dataSegment)
		total += num
		if err != nil {
			println("Error:", err.Error())
			return nil
		}
	}

	if total != int(l) {
		println("Size mismatch")
		println("Target size:", l)
		println("Actual size:", num)
	}

	if err != nil {
		println("Error:", err.Error())
		return nil
	}

	return data
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

func (s serverConnection) readLine() string {
	var char [1]byte
	var buf bytes.Buffer

	defer buf.Reset()

	s.Read(char[:])
	for char[0] != '\n' {
		buf.WriteByte(char[0])
		s.Read(char[:])
	}

	return buf.String()
}
