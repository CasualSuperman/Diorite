package web

import (
	"net/http"

	m "github.com/CasualSuperman/Diorite/multiverse"
)

var multiverse m.Multiverse

func Serve(mult m.Multiverse) {
	multiverse = mult
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
