package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
)

// EXAMPLES
// https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1313&_nkw=water+floss&_sacat=0
// https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1313&_nkw=chicken&_sacat=0
var ebaySearchURLFormat = "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1313&_nkw=%s&_sacat=0"
var csvFileName = "scraped_ebay.csv"

func main() {
	if len(os.Args) != 2 {
		log.Println("Missing search item argument")
		os.Exit(1)
	}
	searchItemArg := os.Args[1]
	searchItem := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(searchItemArg)), " ", "+")

	var previousPageURL string
	nextPageURL := fmt.Sprintf(ebaySearchURLFormat, searchItem)

	// CSV cleanup
	if _, err := os.Stat(csvFileName); err == nil {
		os.Remove(csvFileName)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("unexpected os.Stat error: %v", err)
	}
	// write CSV headers
	writeCSV([]string{"Name","Price","URL"})

	c := colly.NewCollector(
		colly.AllowedDomains("www.ebay.com", "ebay.com"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
		)

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Scraping %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Printf("Scraping error: %s\n", e.Error())
	})

	c.OnHTML("ul.srp-results>li.s-item", func(h *colly.HTMLElement) {
		listSelection := h.DOM
		itemSelection := listSelection.Find("a.s-item__link")
		priceSelection := listSelection.Find("span.s-item__price")

		itemName := strings.TrimSpace(itemSelection.Text())
		itemURL, _ := itemSelection.Attr("href")
		itemPrice := strings.TrimSpace(priceSelection.Text())

		writeCSV([]string{itemName, itemPrice, itemURL})
	})

	c.OnHTML("nav.pagination>a.pagination__next", func(h *colly.HTMLElement) {
		previousPageURL = nextPageURL
		nextPageURL = h.Attr("href")
	})

	var pagesScraped int
	c.OnScraped(func(r *colly.Response){
		pagesScraped += 1
	})

	for {
		if pagesScraped >= 5 || previousPageURL == nextPageURL {
			break
		}
		c.Visit(nextPageURL)
	}
}

func writeCSV(scrapedData []string) {
	csvFile, err := os.OpenFile(csvFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0777)
	if err != nil {
		log.Fatalf("writeCSV OpenFile error: %v", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	if err := csvWriter.Write(scrapedData); err != nil {
		log.Fatalf("writeCSV Write error: %v", err)
	}
}