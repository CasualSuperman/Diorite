package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

var local = flag.Bool("local", false, "Connect to a server running on localhost.")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	canSaveMultiverse := true
	multiverseLoaded := false

	err := os.MkdirAll(StorageDir, os.ModePerm|os.ModeDir)

	if err != nil {
		log.Println("Unable to create application storage directory. Multiverse will not be saved.")
		canSaveMultiverse = false
	}

	var multiverse m.Multiverse

	if canSaveMultiverse {
		multiverseFile, err := os.Open(MultiverseFileName)

		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				log.Println("No local database available. A local copy will be downloaded.")
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Println("Loading local multiverse.")
			multiverse, err = m.Read(multiverseFile)
			multiverseFile.Close()

			if err != nil {
				log.Printf("Unable to load multiverse: %s\n", err)
			} else {
				log.Println("Multiverse loaded.")
				multiverseLoaded = true
			}
		}

	}

	if multiverseLoaded {
		log.Println("Checking for multiverse updates.")
	} else {
		log.Println("Downloading online multiverse.")
	}

	var server serverConnection

	if *local {
		server, err = connectToLocalServer()
	} else {
		server, err = connectToDefaultServer()
	}

	if err != nil {
		log.Println(err)
		if multiverseLoaded == false {
			log.Fatalln("No local multiverse available, and unable to download copy. Unable to continue.")
		}
		log.Println("Warning: Online database unavailable. Your card index may be out of date.")
	} else if server.Modified().After(multiverse.Modified) {
		var saveTo io.Writer

		if canSaveMultiverse {
			saveTo, err = os.Create(MultiverseFileName)
			if err != nil {
				log.Println("Unable to save update to multiverse. Continuing, but it will be redownloaded on next startup.")
				saveTo = nil
			}
		}

		log.Println("Multiverse update available! Downloading now.")
		newM, err := server.DownloadMultiverse(saveTo)

		if err != nil {
			log.Printf("Error downloading: %s\n", err)
			if !multiverseLoaded {
				log.Fatalln("Downloading multiverse failed and no local database available. Unable to continue.")
			}
			log.Println("Unable to download most recent multiverse. Continuing with an out-of-date version.")
		} else {
			log.Println("Multiverse downloaded!")
			log.Println("Cards in multiverse:", newM.Cards.Len())
			multiverse = newM
		}
		server.Close()
	} else {
		log.Println("No updates available.")
		server.Close()
	}

	searchbench, _ := os.Create("searchbench")
	pprof.StartCPUProfile(searchbench)
	for i := 0; i < 1000; i++ {
		multiverse.Search(m.Formats.Standard)
		multiverse.Search(
			m.And{
				m.ManaColors.Blue,
				m.Or{
					m.ManaColors.Green,
					m.ManaColors.Black,
				},
				m.Not{
					m.Or{
						m.ManaColors.White,
						m.ManaColors.Red,
					},
				},
			})
	}
	pprof.StopCPUProfile()

	printedIn := func(setName string) m.Filter {
		return m.Func(func(c *m.Card) bool {
			for _, printing := range c.Printings {
				if printing.Set.Name == setName {
					return true
				}
			}
			return false
		})
	}

	_ = printedIn

	hasText := func(text string) m.Filter {
		text = strings.ToLower(text)
		return m.Func(func(c *m.Card) bool {
			return strings.Contains(strings.ToLower(c.Text), text)
		})
	}

	cards, err := multiverse.Search(
		m.And{
			m.Formats.Standard,
			hasText("life"),
			printedIn("Theros"),
			m.Or{
				m.ManaColors.White,
				m.ManaColors.Green,
				m.ManaColors.Black,
			},
		})

	if err != nil {
		fmt.Printf("Error! %s\n", err.Error())
	} else {
		results := cards.Sort(m.Sorts.Cmc)
		fmt.Printf("%d Results\n", len(results))
		for _, card := range results {
			fmt.Printf("%s\t%s\n=====\n%s\n\n\n\n", card.Name, card.Cost, card.Text)
		}
	}
}
