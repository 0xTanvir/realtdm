package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

func main() {
	var numberOfPage int
	if len(os.Args) > 1 {
		numberOfPage, _ = strconv.Atoi(os.Args[1])
	} else {
		fmt.Print("Total number of page: ")
		fmt.Scanln(&numberOfPage)
	}

	fName := "realtdm.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Status", "Case Number", "Date Created", "Application Number", "Parcel Number", "Sale Date"})

	fmt.Println("Scrapper started......")

	var payload map[string]string
	for i := 0; i < numberOfPage; i++ {
		page := strconv.Itoa(i + 1)
		payload = map[string]string{
			"filterPageNumber":   page,
			"filterCasesPerPage": "100",
			"filterFiltered":     "1",
		}

		// Instantiate default collector
		c := colly.NewCollector()

		c.OnHTML("#county-setup tbody tr", func(e *colly.HTMLElement) {
			writer.Write([]string{
				e.ChildText(".text-left"),
				e.ChildText("td:nth-child(3)"),
				e.ChildText("td:nth-child(4)"),
				e.ChildText("td:nth-child(5)"),
				e.ChildText("td:nth-child(6)"),
				e.ChildText("td:nth-child(7)"),
			})
		})

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		})

		c.OnScraped(func(response *colly.Response) {
			fmt.Println("Scrapped: Page ", page)
		})

		c.Post("https://miamidade.realtdm.com/public/cases/list", payload)

	}

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
