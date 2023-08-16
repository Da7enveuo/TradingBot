package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
func checkStringArrayChanges(old []string, new []string) (string, string, bool) {
	var typ string
	var check bool = false
	if len(old) != len(new) {
		check = true // array length changed, so elements were added or removed

		if len(old) > len(new) {
			typ = "Del"
		} else if len(old) < len(new) {
			typ = "Add"
		} else {
			typ = "none"
		}
	}
	var item string
	for i := 0; i < len(old); i++ {
		if old[i] != new[i] {
			item = new[i]
			check = true
		}
	}
	return item, typ, check // no changes detected
}
*/
func main() {

	// send data of symbols to aws invoke url for lambda function
	// I want to be able to request somewhere in aws for a list of symbols I want real time, then go through the process with that

	// this is to start this off, we go ahead and grab first data sets

	// before this we need to do a for loop and populate a new "EmaCache" map of string (symbol) and struct with 4 ema's/close candle as []float64
	// we track through that by using a yield on each ema compared to price to determine any type of trend
	// we have to update this every time to get rid of last element in ordered hash and insert the new one in first position
	// when we finish on that analysis we will have array of ema yc's and we have to write something to find a pattern, or soemething that is holding as a trend

	// create the cache for ema's for all the symbols before starting up
	// this will also need to be flagged in the for loop, so when a new symbol gets added we pause the functions making requests and collect data for the new symbols previous ema data
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/DailyData")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SetMaxOpenConns(130) // Set maximum open connections
	db.SetMaxIdleConns(130)
	rows, err := db.Query("SELECT Ticker, MAX(Time) FROM DD GROUP BY Ticker")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var symbols []string

	if rows.Next() {
		for rows.Next() {

			var ticker string
			var timea int64
			err := rows.Scan(&ticker, &timea)
			if err != nil {
				panic(err.Error())
			}
			symbols = append(symbols, ticker)
		}
	}

	var wg *sync.WaitGroup
	wg.Add(len(symbols))
	for _, symbol := range symbols {

		// do analysis here
		d, err := GrabIndicators(symbol)
		if err != nil {
			fmt.Println(err)
		}
		analyzeData(d, symbol, wg)
	}
	time.Sleep(3600 * time.Second)

}
func GrabIndicators(symbol string) ([]ReturnJson, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/DailyData")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SetMaxOpenConns(130) // Set maximum open connections
	db.SetMaxIdleConns(130)
	bytes, err := Ta(symbol, db)
	return bytes, err
}

type Re struct {
	Symbol  string
	EmaOnly bool
}
type Req struct {
	Symbols []string
}

// Have this return something other than bytes
// return a array of ReturnJson
func Ta(symbol string, db *sql.DB) ([]ReturnJson, error) {

	// okay so lets do 3 weeks back for daily, 5 months back for weekly, and
	// we do need to specify an appropriate time here.
	// we also need to make multiple queries to get daily, weekly, monthly data.
	// maybe go like 4 weeks back for daily, 5 months back for weekly, and 2-3 years back for Monthly Analysis?
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM D_TA WHERE Ticker='%s' ORDER BY Time DESC LIMIT 20", symbol))
	if err != nil {
		fmt.Println(err)
	}
	var rr []ReturnJson
	if rows.Next() {
		for rows.Next() {
			var r ReturnJson
			err := rows.Scan(&r)
			if err != nil {
				fmt.Println(err)
			}
			rr = append(rr, r)
		}
	}

	return rr, err
}

/*
type CandleDa struct {
	Id     string     `json:"id"`
	Result CandleData `json:"result"`
	Errors []string   `json:"errors"`
}
type CandleResponse struct {
	Data []CandleDa `json:"data"`
}

type EmaResponse struct {
	Data []EmaData `json:"data"`
}
type EmaData struct {
	Id     string   `json:"id"`
	Result E        `json:"result"`
	Errors []string `json:"errors"`
}
type E struct {
	Value     float64 `json:"value"`
	Backtrack int     `json:"backtrack"`
}

type SB struct {
	Symbol  string `json:"symbol"`
	EmaOnly bool   `json:"emaonly"`
}

type Requeststr struct {
	Secret    string `json:"secret"`
	Construct []Cstr `json:"construct"`
}
type Cstr struct {
	Exchange  string      `json:"exchange"`
	Symbol    string      `json:"symbol"`
	Interval  string      `json:"interval"`
	Indicator []Indicator `json:"indicators"`
}

type Indicator struct {
	Indicator  string `json:"indicator"`
	Period     int    `json:"period,omitempty"`
	BackTracks int    `json:"backtracks,omitempty"`
}
*/
