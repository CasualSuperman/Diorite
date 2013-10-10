package main

import (
	_ "encoding/json"
	_ "io/ioutil"
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
