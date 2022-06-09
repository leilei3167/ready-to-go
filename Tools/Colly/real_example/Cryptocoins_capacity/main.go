package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fName := "cap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() //将缓冲区内的全部写入底层Writer
	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	c := colly.NewCollector()

	c.OnHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".currency-name-container"),
			e.ChildText(".col-symbol"),
			e.ChildAttr("a.price", "data-usd"),
			e.ChildAttr("a.volume", "data-usd"),
			e.ChildAttr(".market-cap", "data-usd"),
			e.ChildText(".percent-1h"),
			e.ChildText(".percent-24h"),
			e.ChildText(".percent-7d"),
		})

	})
	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
