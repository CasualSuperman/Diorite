package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var port = flag.String("port", ":5050", "The port to run the server on.")

var multiverseDL []byte
var multiverseModified time.Time

func main() {
	go updateMultiverse()

	log.Println("Starting query server.")

	ln, err := net.Listen("tcp", *port)

	if err != nil {
		log.Fatalln("Unable to bind socket:", err.Error())
	}

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
	var err error

	log.Println("Downloading multiverse.")

	multiverseDL, multiverseModified, err = getMultiverseData()

	if err != nil {
		log.Fatalln("Unable to download multiverse.")
	}

	log.Println("Multiverse downloaded.")

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

	// We want to know the modification time of the multiverse.
	case "multiverseMod":
		conn.Write([]byte(multiverseModified.Format(lastModifiedFormat) + "\n"))
		log.Println("Timestamp accessed.")
		conn.Close()

	// We want to download the multiverse.
	case "multiverseDL":
		conn.Write(multiverseDL)
		log.Println("Multiverse downloaded.")
		conn.Close()

	default:
		answer := fmt.Sprintf("Unrecognized request '%s'.\n")
		log.Printf(answer)
		conn.Write([]byte(answer))
	}
}
