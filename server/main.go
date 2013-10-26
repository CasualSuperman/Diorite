package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"time"
)

var port = flag.String("port", ":5050", "The port to run the server on.")

var multiverseDL []byte
var multiverseModified time.Time

func main() {
	var err error
	log.Println("Downloading multiverse.")

	multiverseDL, multiverseModified, err = getMultiverseData()

	if err != nil {
		log.Fatalln("Unable to download multiverse.")
	}

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

		newData, newMod, err := getMultiverseData()

		if err != nil {
			continue
		}

		log.Println("Update applied.")

		multiverseDL, multiverseModified = newData, newMod
	}
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
