package main

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type CryptoData struct {
	Name string
	Symbol string
	Price string
}

func main() {
	scrapeURL := "https://coinmarketcap.com/"
	var cryptoDataSlice []CryptoData

	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("www.coinmarketcap.com", "coinmarketcap.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Scraping %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Printf("Error: %s\n", e.Error())
	})

	c.OnHTML("tbody tr", func(h *colly.HTMLElement) {
		cSelection := h.DOM
		nameAndSymbol := cSelection.Find("div.sc-16r8icm-0.sc-1teo54s-1.dNOTPP")
		// [ INFO ]
		// nameAndSymbol.Text() returns Bitcoin1BTC

		nsChildNodes := nameAndSymbol.Children().Nodes
		var cryptoName string
		if len(nsChildNodes) > 0 {
			cryptoName = nameAndSymbol.FindNodes(nsChildNodes[0]).Text()
		}
		// [ ALTERNATIVE - cryptoName only ]
		// ??? - p.hKkaxT does not work
		// cryptoName := nameAndSymbol.ChildrenFiltered("p.sc-1eb5slv-0").Text()

		cryptoSymbol := nameAndSymbol.Find("p.coin-item-symbol").Text()
		cryptoPrice := cSelection.Find("div.cLgOOr a.cmc-link").Text() // don't need to target the span

		if cryptoName != "" && cryptoSymbol != "" && cryptoPrice != "" {
			cryptoDataSlice = append(cryptoDataSlice, CryptoData{Name: cryptoName, Symbol: cryptoSymbol, Price: cryptoPrice})
		}
	})

	c.OnScraped(func(r *colly.Response){
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")
		enc.Encode(cryptoDataSlice)
	})

	c.Visit(scrapeURL)
}

