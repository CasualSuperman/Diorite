package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type stsJSON struct {
	Cards []stsCard

	Pagination *stsPagination
}

type stsCard struct {
	Name string `json:"name"`
}

type stsPagination struct {
	TotalCards   int `json:"total-cards"`
	TotalPages   int `json:"total-pages"`
	CurrentPage  int `json:"current-page"`
	CardsPerPage int `json:"cards-per-page"`
}

func getBannedList(format string) ([]string, error) {
	return paginateSTS("banned", format)
}

func getRestrictedList(format string) ([]string, error) {
	return paginateSTS("restricted", format)
}

func paginateSTS(list, format string) ([]string, error) {
	baseUrl := "http://searchthecity.me/api?q=" + list + "%3A" + format + "&p="
	var cards []string
	var page = stsJSON{Pagination: &stsPagination{CurrentPage: 0, TotalPages: 1}}
	i := 0

	for page.Pagination.CurrentPage < page.Pagination.TotalPages {
		url := baseUrl + strconv.Itoa(page.Pagination.CurrentPage+1)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&page)

		if err != nil {
			return nil, err
		}

		if page.Pagination == nil {
			return nil, nil
		}

		if i == 0 {
			cards = make([]string, page.Pagination.TotalCards)
		}

		for _, card := range page.Cards {
			cards[i] = card.Name
			i++
		}
	}

	return cards, nil
}
