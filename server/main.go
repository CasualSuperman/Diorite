package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

var port = flag.String("port", ":5050", "The port to run the server on.")
var test = flag.Bool("travis", false, "Exit after a single download for testing.")

var multiverse m.Multiverse

var banlist []formatList
var banlistHash uint64
var banlistModified time.Time

var downloadData []byte
var downloadModified time.Time

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)
	flag.Parse()

	var initialDownloadComplete sync.WaitGroup

	initialDownloadComplete.Add(1)

	go doDownload(&initialDownloadComplete)

	go watchForUpdates()

	handleConnections(&initialDownloadComplete)
}

func doDownload(done *sync.WaitGroup) {
	defer done.Done()
	multiverseChan := make(chan m.Multiverse)
	banlistChan := make(chan formatList, len(m.Formats.List))
	errChan := make(chan error)
	updates := false

	go getFormatLists(banlistChan, errChan)

	if mod, err := onlineModifiedAt(); err == nil && mod.After(multiverse.Modified) {
		go getMultiverse(multiverseChan, errChan)

		select {
		case newMultiverse := <-multiverseChan:
			if multiverse.Modified.IsZero() || newMultiverse.Modified.After(multiverse.Modified) {
				multiverse = newMultiverse
				updates = true
			}
		case err := <-errChan:
			log.Fatal(err.Error())
		}
	} else if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Multiverse up to date.")
	}

	formats := make([]formatList, len(m.Formats.List))
	formatsLeft := len(formats)

	for formatsLeft > 0 {
		select {
		case format := <-banlistChan:
			for i, f := range m.Formats.List {
				if f == format.Format {
					formats[i] = format
					formatsLeft--
				}
			}
		case err := <-errChan:
			log.Fatal(err.Error())
		}
	}

	hash := generateFormatsHash(formats)

	if hash != banlistHash {
		updates = true
		banlistHash = hash
		banlistModified = time.Now()
	} else {
		log.Println("Banned/restricted card list up to date.")
	}

	if updates {
		log.Println("Clearing banlists.")
		clearBanlists()

		log.Println("Marking banned/restricted cards.")
		for i, card := range multiverse.Cards {
			for _, formatList := range formats {
				if formatList.Banned[card.Name] {
					multiverse.Cards[i].Banned = append(multiverse.Cards[i].Banned, formatList.Format)
				}
				if formatList.Restricted[card.Name] {
					multiverse.Cards[i].Restricted = append(multiverse.Cards[i].Restricted, formatList.Format)
				}
			}
		}

		downloadModified = time.Now()
		multiverseUpdate := multiverse.Modified
		multiverse.Modified = downloadModified

		log.Println("Finalizing multiverse payload.")
		var b bytes.Buffer
		err := multiverse.Write(&b)

		if err != nil {
			log.Fatalln(err.Error())
		}

		downloadData = b.Bytes()
		multiverse.Modified = multiverseUpdate

		b.Reset()
	}

	log.Println("Multiverse ready.")
}

func watchForUpdates() {
	updateTick := time.Tick(time.Hour * 12)
	var unused sync.WaitGroup

	for {
		unused.Add(1)
		<-updateTick
		doDownload(&unused)
	}
}

func handleConnections(ready *sync.WaitGroup) {
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
			go provideDownload(conn, *test, ready)
		}
	}
}

func provideDownload(conn net.Conn, done bool, ready *sync.WaitGroup) {
	s := bufio.NewScanner(conn)

	defer func() {
		if done {
			log.Println("Exiting after single connection.")
			os.Exit(0)
		}
	}()

	ready.Wait()

	for s.Scan() {
		switch text := s.Text(); text {
		// We want to know the modification time of the multiverse.
		case "multiverseMod":
			conn.Write([]byte(downloadModified.Format(lastModifiedFormat) + "\n"))
			log.Println("Timestamp accessed.")

		// We want to download the multiverse.
		case "multiverseDL":
			conn.Write(downloadData)
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
