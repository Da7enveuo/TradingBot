package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

		})
	}
	return symbols
}

func Req(symbol string, timea string, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	url := "https://query1.finance.yahoo.com/v7/finance/download/%v?period1=%v&period2=%v&interval=%v&events=history&includeAdjustedClose=true"
	day_req_url := fmt.Sprintf(url, symbol, timea, time.Now().Unix(), "1d")
	d_req, err := http.Get(day_req_url)
	if err != nil {
		fmt.Println(err)
	}
	err = CSV_SQL(symbol, db, d_req)
	if err != nil {
		fmt.Println(err)
	}

}

// here we adjust the high and low  witht he ratio of the adjusted close vs close and do that for high/low
func CSV_SQL(ticker string, db *sql.DB, bd *http.Response) error {

	defer bd.Body.Close()
	reader := csv.NewReader(bd.Body)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	ct := 0

	for _, row := range records {
		if ct == 0 {
			ct += 1
		} else {
			stmt, err := db.Prepare("INSERT INTO DD (Time, Ticker, Open, High, Low, Close, Volume) VALUES (?, ?, ?, ?, ?, ?, ?);")
			if err != nil {
				fmt.Println(err)
			}
			var day_sql_entry SqlEntry

			adj, err := strconv.ParseFloat(row[5], 64)
			if err != nil {
				fmt.Println(err)
			}
			close, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				fmt.Println(err)
			}
			adj_multiplier := adj / close

			for c_index, column_value := range row {
				// at c_index 0 is date, 1 is open, 2 is high, 3 is low, 4 is close, 5 is adjusted close, 6 is volume.
				// want to append these to gether and insert as a record to the db table. Also want to convert these dates. Comes in format 3/10/2023
				switch c_index {
				case 0:
					date, err := time.Parse("2006-01-02", column_value)
					if err != nil {
						fmt.Println(err)
					}
					day_sql_entry.Time = date.Unix()
				case 1:
					f, err := strconv.ParseFloat(column_value, 64)
					if err != nil {
						fmt.Println(err)
					}
					day_sql_entry.Open = f * adj_multiplier
				case 2:
					f, err := strconv.ParseFloat(column_value, 64)
					if err != nil {
						fmt.Println(err)
					}
					day_sql_entry.High = f * adj_multiplier
				case 3:
					f, err := strconv.ParseFloat(column_value, 64)
					if err != nil {
						fmt.Println(err)
					}
					day_sql_entry.Low = f * adj_multiplier
				case 5:
					f, err := strconv.ParseFloat(column_value, 64)
					if err != nil {
						fmt.Println(err)
					}
					day_sql_entry.Close = f
				case 6:
					f, err := strconv.ParseInt(column_value, 10, 64)
					if err != nil {
						fmt.Println(err)
					}
					day_sql_entry.Volume = f
				default:
					continue
				}
			}
			_, err = stmt.Exec(day_sql_entry.Time, string(ticker), day_sql_entry.Open, day_sql_entry.High, day_sql_entry.Low, day_sql_entry.Close, day_sql_entry.Volume)
			if err != nil {
				fmt.Println(err)
			}
			stmt.Close()
			// Insert a new record into the database
		}
	}

	return nil
}

// func GrabNewSymbolData (symbols []string) (bool){
func PopulateNewSymbols(symbols []string) {
	url := "https://query1.finance.yahoo.com/v7/finance/download/%v?period1=%v&period2=%v&interval=%v&events=history&includeAdjustedClose=true"
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/DailyData")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	for _, symbol := range symbols {
		day_req_url := fmt.Sprintf(url, symbol, "728265600", time.Now().Unix(), "1d")

		d_req, err := http.Get(day_req_url)

		if err != nil {
			fmt.Println(err)
		}
		defer d_req.Body.Close()
		err = CSV_SQL(symbol, db, d_req)
		if err != nil {
			fmt.Println(err)
		}
	}
	time.Sleep(time.Second * 60)
}
