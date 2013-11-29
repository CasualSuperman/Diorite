package web

import (
	"net/http"

	ws "code.google.com/p/go.net/websocket"
	m "github.com/CasualSuperman/Diorite/multiverse"
)

var multiverse m.Multiverse

func Serve(mult m.Multiverse) {
	multiverse = mult
	http.Handle("/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/name/", http.StripPrefix("/name/", http.HandlerFunc(nameSearch)))
	http.Handle("/ws", ws.Handler(websocketServer))
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
	Type string
	Request string
}

type wsResponse struct {
	Type string
	Response interface{}
}

func websocketServer(w *ws.Conn) {
	open := true

	for open {
		var msg wsRequest
		err := ws.JSON.Receive(w, &msg)
		if err != nil {
			println(err.Error())
			break
		}

		switch msg.Type {
		case "close":
			open = false
		case "nameSearch":
			cards := multiverse.FuzzyNameSearch(msg.Request, 15)
			names := make([]string, len(cards))
			for i, card := range cards {
				names[i] = card.Name
			}
			ws.JSON.Send(w, wsResponse{
				"nameSearch",
				names,
			})
		default:
			println("Type:", msg.Type)
		}
	}
}
