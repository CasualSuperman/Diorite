package multiverse

import (
	"encoding/gob"
	"io"
	"sync"
	"time"
)

type jsonRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

type jsonCard struct {
	Name string `json:"name"`

	ManaCost string   `json:"manaCost"`
	Cmc      float32  `json:"cmc"`
	Colors   []string `json:"colors"`

	Type       string   `json:"type"`
	Supertypes []string `json:"supertypes"`
	Types      []string `json:"types"`
	Subtypes   []string `json:"subtypes"`

	Rarity string `json:"rarity"`

	Text string `json:"text"`

	Flavor string `json:"flavor"`

	Artist string `json:"artist"`
	Number string `json:"number"`

	Power     string `json:"power"`
	Toughness string `json:"toughness"`

	Layout       string `json:"layout"`
	MultiverseId int    `json:"multiverseid"`

	ImageName string `json:"imageName"`

	Rulings []jsonRuling `json:"rulings"`
}

type jsonSet struct {
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	ReleaseDate string     `json:"releaseDate"`
	Border      string     `json:"border"`
	Type        string     `json:"type"`
	Block       string     `json:"block"`
	Cards       []jsonCard `json:"cards"`
}

type OnlineMultiverse struct {
	Sets     map[string]jsonSet
	Modified time.Time
}

func (o OnlineMultiverse) WriteTo(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(o)
}

func inflateOnlineMultiverse(r io.Reader) OnlineMultiverse {
	d := gob.NewDecoder(r)
	var om OnlineMultiverse
	d.Decode(&om)
	return om
}

var processedCards = struct {
	sync.RWMutex
	cards map[string]*Card
}{cards: make(map[string]*Card)}

func (jc *jsonCard) convert() *Card {
	return nil
}

func CardFromJson(jc *jsonCard) *Card {
	processedCards.RLock()
	if c, ok := processedCards.cards[jc.Name]; ok {
		processedCards.RUnlock()
		return c
	}
	processedCards.RUnlock()
	processedCards.Lock()
	c := new(Card)
	processedCards.cards[jc.Name] = c
	processedCards.Unlock()

	copyCardFields(jc, c)

	return c
}
