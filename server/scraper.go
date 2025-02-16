package main

import (
	"net/http"
	"net/url"
	"regexp"
	"sort"

	"golang.org/x/net/html"
)

var resultsLink = regexp.MustCompile("^/query.*&p=\\d+$")
var pageLinkText = regexp.MustCompile("^\\[.*\\]$")

const baseUrl = "http://magiccards.info/query?q="

func getGenericList(baseUrl string) (cardList []string, err error) {
	urls := make(map[string]string)
	cards := make(map[string]string)

	pageUrl, err := url.Parse(baseUrl)

	if err != nil {
		return
	}

	page, err := http.Get(pageUrl.String())

	if err != nil {
		return
	}

	root, err := html.Parse(page.Body)
	page.Body.Close()

	if err != nil {
		return
	}

	scrape(root, urls, cards)

	for _, url := range urls {
		pageUrl, err = pageUrl.Parse(url)

		if err != nil {
			return
		}

		if pageUrl.Query()["p"][0] == "1" {
			continue
		}

		page, err = http.Get(pageUrl.String())

		if err != nil {
			return
		}

		root, err = html.Parse(page.Body)
		page.Body.Close()

		if err != nil {
			return
		}

		scrape(root, urls, cards)
	}

	for card := range cards {
		cardList = append(cardList, card)
	}

	sort.Strings(cardList)
	return
}

func getRestrictList(format string) ([]string, error) {
	return getGenericList(baseUrl + "restricted%3A" + format)
}

func getBanList(format string) ([]string, error) {
	return getGenericList(baseUrl + "banned%3A" + format)
}

func isPageLink(n *html.Node) bool {
	return n.FirstChild != nil && n.FirstChild.Type == html.TextNode && pageLinkText.MatchString(n.FirstChild.Data)
}

func scrape(n *html.Node, urls, cards map[string]string) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" {
					if resultsLink.MatchString(a.Val) && isPageLink(n) {
						urls[n.FirstChild.Data] = a.Val
					}
					break
				}
			}
		case "img":
			cardName := ""
			isCard := false
			for _, img := range n.Attr {
				switch img.Key {
				case "width":
					if img.Val == "312" {
						isCard = true
					}
				case "alt":
					cardName = img.Val
				}
			}
			if isCard {
				cards[cardName] = ""
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrape(c, urls, cards)
	}
}
