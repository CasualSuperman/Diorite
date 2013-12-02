package web

import (
	"net/http"

	m "github.com/CasualSuperman/Diorite/multiverse"
	gs "github.com/CasualSuperman/gosocket"
)

var multiverse m.Multiverse

func Serve(mult m.Multiverse) {
	multiverse = mult

	socketServer := gs.NewServer()
	http.Handle("/gs/", socketServer)

	socketServer.Handle("nameSearch", fuzzyNameSearch)
	socketServer.Handle("card", sendCard)

	http.Handle("/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/name/", http.StripPrefix("/name/", http.HandlerFunc(nameSearch)))
	http.ListenAndServe(":6060", nil)
}

func nameSearch(rq http.ResponseWriter, req *http.Request) {
	cards := multiverse.FuzzyNameSearch(req.URL.Path, 1)
	if len(cards) < 1 {
		rq.WriteHeader(http.StatusNotFound)
		rq.Write([]byte("Unable to locate card."))
	} else {
		rq.Write([]byte(cards[0].String()))
	}
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

func fuzzyNameSearch(msg gs.Msg) {
	var searchTerm string
	msg.Receive(&searchTerm)

	results := multiverse.FuzzyNameSearch(searchTerm, 15)
	cards := make([]webCard, len(results))
	for i, card := range results {
		cards[i].Name = card.Name
		cards[i].MultiverseID = int(card.Printings[len(card.Printings)-1].ID)
	}

	msg.Respond(cards)
}

func sendCard(msg gs.Msg) {
	var id m.MultiverseID
	msg.Receive(&id)
	cards, _ := multiverse.Search(id)
	if len(cards) > 0 {
		msg.Respond(cards[0])
	} else {
		msg.Respond(nil)
	}
}
