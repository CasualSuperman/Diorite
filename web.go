package main

import (
	"net/http"

	m "github.com/CasualSuperman/Diorite/multiverse"
	gs "github.com/CasualSuperman/gosocket"
)

type server struct {
	multiverse *m.Multiverse
}

func (s *server) UpdateMultiverse(mlt *m.Multiverse) {
	s.multiverse = mlt
}

func NewServer(mlt *m.Multiverse) *server {
	return &server{mlt}
}

func (s *server) Serve(port string, exit chan exitSignal) {
	socketServer := gs.NewServer()
	http.Handle("/gs/", socketServer)

	socketServer.Handle("nameSearch", s.fuzzyNameSearch)
	socketServer.Handle("card", s.sendCard)
	socketServer.Errored(func(err error) {
		println(err.Error())
	})

	if !*keepserver {
		socketServer.On(gs.Disconnect, func(c *gs.Conn) {
			exit <- exitSignal{0, "Web interface closed."}
		})
	}

	http.Handle("/", http.FileServer(http.Dir("./content")))
	http.ListenAndServe(port, nil)
}

type wsRequest struct {
	Type    string
	Request string
}

type wsResponse struct {
	Type     string
	Response interface{}
}

type webCard struct {
	Name         string
	MultiverseID int
}

func (s *server) fuzzyNameSearch(msg gs.Msg) {
	var searchTerm string
	msg.Receive(&searchTerm)

	results := s.multiverse.FuzzyNameSearch(searchTerm, 15)
	cards := make([]webCard, len(results))
	for i, card := range results {
		cards[i].Name = card.Name
		cards[i].MultiverseID = int(card.Printings[len(card.Printings)-1].ID)
	}

	msg.Respond(cards)
}

func (s *server) sendCard(msg gs.Msg) {
	var id m.MultiverseID
	msg.Receive(&id)
	cards, _ := s.multiverse.Search(id)
	if len(cards) > 0 {
		err := msg.Respond(cards[0])
		if err != nil {
			println(err.Error())
		}
	} else {
		println("no cards")
		msg.Respond(nil)
	}
}
