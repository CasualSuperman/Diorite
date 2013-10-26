package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"net"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

var port = flag.String("port", ":5050", "The port to run the server on.")
var multiverseModified time.Time
var multiverseDL []byte

func main() {
	log.Println("Downloading multiverse.")

	multiverse, err := downloadMultiverse()

	if err != nil {
		log.Fatalln("Unable to download multiverse.")
	}

	getDlData(multiverse)
	multiverse = nil

	log.Println("Multiverse downloaded.")

	go updateMultiverse()

	log.Println("Starting server.")

	ln, err := net.Listen("tcp", *port)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
		} else {
			go provideDownload(conn)
		}
	}
}

func updateMultiverse() {
	dailyUpdate := time.Tick(time.Hour * 12)

	for {
		<-dailyUpdate
		log.Println("Checking for update.")
		mod, err := onlineModifiedAt()

		if mod == multiverseModified {
			log.Println("Multiverse up-to-date.")
			continue
		}

		if err != nil {
			continue
		}

		log.Println("Update found! Downloading...")
		newM, err := downloadMultiverse()

		if err != nil {
			continue
		}

		log.Println("Update applied.")

		getDlData(newM)
		multiverseModified = newM.Modified
	}
}

func getDlData(multiverse *m.Multiverse) {
	var b bytes.Buffer
	multiverse.Write(&b)
	multiverseDL = b.Bytes()
}

func provideDownload(conn net.Conn) {
	s := bufio.NewScanner(conn)
	s.Scan()

	switch text := s.Text(); text {
	case "multiverseMod":
		log.Println("Timestamp accessed.")
		conn.Write([]byte(multiverseModified.Format(lastModifiedFormat) + "\n"))
		conn.Close()
	case "multiverseDL":
		log.Println("Multiverse downloaded.")
		conn.Write(multiverseDL)
		conn.Close()
	default:
		log.Printf("Unrecognized request '%s'.\n", text)
	}
}
