package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var port = flag.String("port", ":5050", "The port to run the server on.")
var test = flag.Bool("travis", false, "Exit after a single download for testing.")

var multiverseDL []byte
var multiverseModified time.Time
var multiverseReady sync.RWMutex

func main() {
	flag.Parse()
	multiverseReady.Lock()
	go updateMultiverse(&multiverseReady)

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
			go provideDownload(conn, *test)
		}
	}
}

func updateMultiverse(lock *sync.RWMutex) {
	var err error

	log.Println("Downloading multiverse.")

	multiverseDL, multiverseModified, err = getMultiverseData()

	lock.Unlock()

	if err != nil {
		log.Println("Error!", err.Error())
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
			log.Println("Error getting multiverse data:", err.Error())
			continue
		}

		log.Println("Update applied.")

		lock.Lock()
		multiverseDL, multiverseModified = newData, newMod
		lock.Unlock()
	}
}

func provideDownload(conn net.Conn, done bool) {
	s := bufio.NewScanner(conn)

	multiverseReady.RLock()
	defer multiverseReady.RUnlock()

	defer func() {
		if done {
			log.Println("Exiting after single connection.")
			os.Exit(0)
		}
	}()

	for s.Scan() {
		switch text := s.Text(); text {
		// We want to know the modification time of the multiverse.
		case "multiverseMod":
			conn.Write([]byte(multiverseModified.Format(lastModifiedFormat) + "\n"))
			log.Println("Timestamp accessed.")

		// We want to download the multiverse.
		case "multiverseDL":
			conn.Write(multiverseDL)
			log.Println("Multiverse downloaded.")

		// We're done, close the connection.
		case "close":
			conn.Close()
			log.Println("Client disconnected.")
			return

		default:
			answer := fmt.Sprintf("Unrecognized request '%s'.\n")
			log.Printf(answer)
			conn.Write([]byte(answer))
		}
	}
}
