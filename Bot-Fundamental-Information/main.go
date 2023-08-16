package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Cik struct {
	Ticker string
	Cik    string
	// padded CIK
	AdjustedCik string
}

// use https://www.sec.gov/files/company_tickers.json to get a json of all comapnies and their corresponding cik's (or https://www.sec.gov/include/ticker.txt, might be better with this one, just text based.)
// then use https://data.sec.gov/submissions/CIK##########.json and replace with cik number.This gives us some useless shit, we need to get the 10k data from the xbrt (busines financial xml reporting format)
// then we can store the information in our database.
// please note that CIK's are 10 digits, if we get one that is not 10 digits we just have to pad the beginning with 0's util we reach 10 digits.
// 0000789019
func removeDuplicates(arr []string) []string {
	unique := []string{}
	visited := make(map[string]bool)

	for _, num := range arr {
		if !visited[num] {
			unique = append(unique, num)
			visited[num] = true
		}
	}

	return unique
}
func NdaqSymbols() []string {

	var symbols []string
	resp, err := http.Get("https://en.wikipedia.org/wiki/Nasdaq-100")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the table with id "constituents"
	table := doc.Find("#constituents")

	// Check if the table was found
	if table.Length() == 0 {
		fmt.Println("Table not found")

	} else {

		// Find all the rows in the table, excluding the header row
		rows := table.Find("tbody tr")

		// Iterate through the rows and print the values
		rows.Each(func(i int, row *goquery.Selection) {
			// Find the cells in the row
			cells := row.Find("td")
			symbol := cells.Eq(1).Text()
			symbols = append(symbols, symbol)

		})
	}
	return symbols
}
func GetDJIASymbols() []string {
	// Make an HTTP GET request to fetch the HTML content
	resp, err := http.Get("https://en.wikipedia.org/wiki/Dow_Jones_Industrial_Average")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	table := doc.Find("#constituents")
	if table.Length() == 0 {
		fmt.Println("Table not found")

	}
	// Find the <tbody> element
	tbody := table.Find("tbody")
	var symbols []string
	// Iterate over each <tr> element in the <tbody>
	tbody.Find("tr").Each(func(i int, tr *goquery.Selection) {
		// Find the <th> and <td> elements within the <tr>

		tds := tr.Find("td")

		symbol := tds.Eq(1).Find("a").Text()

		symbols = append(symbols, symbol)
		// Print the extracted data

	})
	return symbols
}
func GetSPYSymbols() []string {
	// Make an HTTP GET request to the URL
	resp, err := http.Get("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the table with id "constituents"
	table := doc.Find("#constituents")
	var symbols []string
	// Check if the table was found
	if table.Length() == 0 {
		fmt.Println("Table not found")

	} else {

		// Find all the rows in the table, excluding the header row
		rows := table.Find("tbody tr")

		// Iterate through the rows and print the values
		rows.Each(func(i int, row *goquery.Selection) {
			// Find the cells in the row
			cells := row.Find("td")

			// Extract the values from the cells
			symbol := cells.Eq(0).Find("a").Text()
			symbols = append(symbols, symbol)

			// Print the values

		})
	}
	return symbols
}
func main() {
	symbols := GetSPYSymbols()

	dj := GetDJIASymbols()

	nd := NdaqSymbols()

	symbols = append(symbols, nd...)
	symbols = append(symbols, dj...)
	symbols = removeDuplicates(symbols)
	fmt.Println(len(symbols))
	resp, err := http.Get("https://www.sec.gov/include/ticker.txt")
	if err != nil {
		fmt.Println("Error fetching ticker data:", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	//var ciks []Cik
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		symbol := parts[0]
		cik := parts[1]
		for _, sym := range symbols {
			if symbol == strings.ToLower(sym) {
				fmt.Printf("Symbol: %s, CIK: %s\n", sym, cik)
			}
		}

	}

	// honestly not sure what happens after this
	/*
		sec_url := "https://data.sec.gov/submissions/CIK%s.json"
		for _, symbol := range ciks {
			ur := fmt.Sprintf(sec_url, symbol.AdjustedCik)

			resp, err := http.Get(ur)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning ticker data:", err)
		}
	*/
}
