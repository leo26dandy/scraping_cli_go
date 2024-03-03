package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// initiation variables
var (
	urls         string
	cssSelector  string
	autoPaginate bool
	linkSelector string
)

// initiate flags to provided variables
func init() {
	flag.StringVar(&urls, "urls", "", "Comma-separated list of URLs to scrape.")
	flag.StringVar(&cssSelector, "selector", "", "CSS selector for the elements to extract.")
	flag.BoolVar(&autoPaginate, "auto-paginate", false, "Enable auto pagination based on a pattern in the next page links.")
	flag.StringVar(&linkSelector, "link-selector", "", "CSS selector for links to follow and scrape.")
	flag.Parse()

	if urls == "" || cssSelector == "" {
		log.Fatal("Missing required flags: -urls and -selector.")
	}
}

// main runs here
func main() {
	c := colly.NewCollector()
	count := 0

	// Define callback for handling scraped data
	c.OnHTML(cssSelector, func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		count++
	})

	// Implement auto pagination (optional)
	if autoPaginate {
		nextPagePattern := regexp.MustCompile(`page=(\d+)`) // Replace with your specific pattern
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			nextLink := e.Attr("href")
			if match := nextPagePattern.FindStringSubmatch(nextLink); match != nil {
				nextPage := fmt.Sprintf("https://www.example.com/page/%s", match[1]) // Replace with base URL
				c.Visit(nextPage)
			}
		})
	}

	// Implement link following (optional)
	if linkSelector != "" {
		c.OnHTML(linkSelector, func(e *colly.HTMLElement) {
			c.Visit(e.Attr("href"))
		})
	}

	// Visit the specified URLs
	for _, url := range strings.Split(urls, ",") {
		err := c.Visit(url)
		if err != nil {
			log.Println(err)
		}
	}

	// Print Sum of scraped data
	fmt.Println("\nTotal Scraped Data:", count)
}
