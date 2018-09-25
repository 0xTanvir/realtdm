package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fName := "realtdm.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Just for 1st page
	payload := map[string]string{
		"filterPageNumber": "1",
		"filterCasesPerPage":   "100",
		"filterFiltered": "1",
	}

	// Write CSV header
	writer.Write([]string{"Status", "Case Number", "Date Created", "Application Number", "Parcel Number"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("#county-setup tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(".text-left"),
			e.ChildText("td:nth-child(3)"),
			e.ChildText("td:nth-child(4)"),
			e.ChildText("td:nth-child(5)"),
			e.ChildText("td:nth-child(6)"),
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	c.Post("https://miamidade.realtdm.com/public/cases/list",payload)

	log.Printf("Scraping finished, check file %q for results\n", fName)
}