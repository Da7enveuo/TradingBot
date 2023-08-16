package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*

Program Structure:

2 tables in mysql database:
	- Daily Data (keeps data on ohlcv)
	- TA Data (keeps TA data based on the timeframes in daily data, sorted on epoch keys as primary key)
Main:
	loops through a list of symbols
		if we have data in the database then continue from there and populate to the most recent daily data
		if not then we pull all historic data and populate in the data base

	have a secondary service that will uniq out the symbols, and retrieve any missing data from the most recent date fetched. Will run once a day.
*/

type SqlEntry struct {
	DBid   int   `json:"Dbid,omitempty"`
	Time   int64 // epoch time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

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
func isEpochToday(epochTime int64) bool {
	// Convert the epoch time to a time.Time value
	t := time.Unix(epochTime, 0)

	// Get the current date
	currentDate := time.Now().Local().Format("2006-01-02")

	// Format the epoch time as a date string
	epochDate := t.Local().Format("2006-01-02")

	// Compare the current date with the epoch date
	return currentDate == epochDate
}
func DailyCheck() {
	// probably should storee these creds better haha
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
	var wg sync.WaitGroup
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
			// put in the latest ticker time we have and now time to get the newest data.
			if !isEpochToday(timea) {
				//notWg = true
				wg.Add(1)
				go Req(ticker, strconv.FormatInt(timea, 10), db, &wg)

			} else {
				fmt.Printf("Not today %s\n", ticker)
			}

		}
		wg.Wait()
	} else {
		symbols = GetSPYSymbols()
		symbols = append(symbols, NdaqSymbols()...)
		symbols = append(symbols, GetDJIASymbols()...)
		symbols = removeDuplicates(symbols)
		// run function
		for _, symbol := range symbols {
			wg.Add(1)
			go Req(symbol, "1", db, &wg)
		}
		wg.Wait()
		fmt.Println("Past the waiting")
	}

	TechnicalsCheck(db, symbols)
	// this marks the end of the checking and injesting of the DD table data.

	// after this we have our data. Now we need to populate our data based on it .do this in a loop for each thing we need to do?
	//SELECT Ticker, FROM_UNIXTIME(MAX(Time)) as Latest_Time, FROM_UNIXTIME(MIN(Time)) as First_Time FROM DD GROUP BY Ticker;

}

func TechnicalsCheck(db *sql.DB, symbols []string) {
	var row_ int64
	var Pulled_data []SqlEntry
	// iterate through the symbols
	fmt.Println("Got To Technicals Check")
	for _, symbol := range symbols {
		fmt.Println("working here")
		/*
			basically we:
				for each symbol:
					- query D_TA for most recent time
					- use that time found to query DD for the entry at that time and getting the 200 entries before (for moving averages n shit)
					- calculate all the technicals for that time period
					- Push the data to D_TA
					- increment and get rid of the oldest cached data.
					- keep going. If we have mor than 200 to go (200 timeframe difference between DD and D_TA) then we just continue past the cached portion and continue.

				We will need a different function for when we start vs when we finish.
				Also some of the technical indicators will be available earlier than others.
				And another thing, must also keep track of database id's sice they are primary keyed
				Last thing, the close to adj close have a similar ratio, so lets adjust the high low open to match the adjustment, therefore we have no strange datapoints Caused by stock splits.
		*/

		// get the max time from DD table
		// get the max time from D_TA table
		rows, err := db.Query(fmt.Sprintf("SELECT MAX(Time) FROM D_TA WHERE Ticker='%s'", strings.ToLower(symbol)))

		if err != nil {
			fmt.Println("Not Found, setting to 1")
			row_ = 1
		} else {

			for rows.Next() {
				var maxTime int64
				err := rows.Scan(&maxTime)
				if err != nil {
					fmt.Println("Not Found, setting to 1")
					row_ = 1
				} else {
					// Use the value of maxTime here
					fmt.Println("Maximum Time:", maxTime)
					row_ = maxTime
				}
			}
		}
		rows.Close()
		// use max time we got from D_TA and proceed from there, grabbing all data that is higher than what we put in, and including data 200 before for our ema and sma calculations.
		roww, err := db.Query(fmt.Sprintf("SELECT DBid, Time, Open, High, Low, Close, Volume FROM DD WHERE Ticker='%s' AND Time<%v ORDER BY Time DESC LIMIT 200", strings.ToLower(symbol), row_))
		if err != nil {
			fmt.Println(err)
			row_ = 1
		}
		// this stays persistant so we don't pull multiple times.

		if roww.Next() {
			for roww.Next() {
				var s SqlEntry
				// iterate through the data collected
				// so here is the thing, we want to grab sql data back from this period. SO need to run sql query here and pull enough periods back.
				// then as we iterate through we keep the data but just increment up depending on what we need
				err = roww.Scan(&s.DBid, &s.Time, &s.Open, &s.High, &s.Low, &s.Close, &s.Volume)
				if err != nil {
					fmt.Println(err)
				}
				Pulled_data = append(Pulled_data, s)
			}

		}
		roww.Close()
		// grab remaining data and append to end of the cache.
		rowz, err := db.Query(fmt.Sprintf("SELECT DBid, Time, Open, High, Low, Close, Volume FROM DD WHERE Ticker='%s' AND Time >%v ORDER BY Time DESC", strings.ToLower(symbol), row_))
		if err != nil {
			fmt.Println(err)
			row_ = 1
		}
		// this stays persistant so we don't pull multiple times.
		if rowz.Next() {
			for rowz.Next() {
				var s SqlEntry
				// iterate through the data collected
				// so here is the thing, we want to grab sql data back from this period. SO need to run sql query here and pull enough periods back.
				// then as we iterate through we keep the data but just increment up depending on what we need
				err = rowz.Scan(&s.DBid, &s.Time, &s.Open, &s.High, &s.Low, &s.Close, &s.Volume)
				if err != nil {
					fmt.Println(err)
				}
				Pulled_data = append(Pulled_data, s)
			}

		}
		rowz.Close()

		rw, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM DD WHERE Ticker = '%v';", strings.ToLower(symbol)))
		if err != nil {
			fmt.Println(err)
		}

		var DD_ct int
		if rw.Next() {
			fmt.Println("Got in here")
			err := rw.Scan(&DD_ct)
			if err != nil {
				fmt.Println(err)
			}

		}
		rw.Close()
		rd, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM D_TA WHERE Ticker = '%v';", strings.ToLower(symbol)))

		if err != nil {
			fmt.Println(err)
		}

		var d_ct int
		if rd.Next() {
			err := rd.Scan(&d_ct)
			if err != nil {
				fmt.Println(err)
			}
		}
		rd.Close()
		var da int = DD_ct - d_ct

		var adjclose []float64
		var high []float64
		var low []float64
		var volume []int64
		var time []int64
		var dbid []int
		for _, s := range Pulled_data {
			high = append(high, s.High)
			low = append(low, s.Low)
			volume = append(volume, s.Volume)
			adjclose = append(adjclose, s.Close)
			time = append(time, s.Time)
			dbid = append(dbid, s.DBid)
		}
		fmt.Printf("Ticker: %v, DD_ct: %v, d_ct: %v \n", symbol, DD_ct, d_ct)
		// here we loop through the difference between 2 symbols in each table, The data in DD is in Pulled_data
		for i := 0; i <= da-1; i++ {
			fmt.Println("Got Here")
			// now do a loop to grab an array of the data I want and throw in here

			// something is fucked up with this
			// 2214 - 20 - 1 : 2214 - 1+1 2193:2214
			cci := CalculateCommodityChannelIndex(high[len(high)-20-i:len(high)-i+1], low[len(high)-20-i:len(high)-i+1], adjclose[len(high)-20+i:len(high)-i+1])
			fmt.Printf("CCI: %v\n", cci)
			ema20 := CalculateEMA(adjclose[len(high)-20-i:len(high)-i], 20)[0]
			fmt.Printf("EMA20: %v\n", ema20)
			ema50 := CalculateEMA(adjclose[len(high)-50-i:len(high)-i], 50)[0]
			fmt.Printf("EMA50: %v\n", ema50)
			ema100 := CalculateEMA(adjclose[len(high)-100-i:len(high)-i], 100)[0]
			fmt.Printf("EMA100: %v\n", ema100)
			ema200 := CalculateEMA(adjclose[len(high)-200-i:len(high)-i], 200)[0]
			fmt.Printf("EMA200: %v\n", ema200)
			macd := CalculateMACD(adjclose[len(high)-26-i:len(high)-i], 26)
			fmt.Printf("MACD: %v\n", macd)
			// for macd, need a 9 period moving average of the macd as well, write another fuction I am too tired.
			mfi := CalculateMoneyFlowIndex(adjclose[len(high)-14-i:len(high)-i+1], high[len(high)-14-i:len(high)-i+1], low[len(high)-14-i:len(high)-i+1], volume[len(high)-14-i:len(high)-i+1], 14)
			fmt.Printf("MFI: %v\n", mfi)
			sma20 := CalculateSMA(adjclose[len(high)-20+i:len(high)-i], 20)[0]
			fmt.Printf("SMA20: %v\n", sma20)
			sma50 := CalculateSMA(adjclose[len(high)-50+i:len(high)-i], 50)[0]
			fmt.Printf("SMA50: %v\n", sma50)
			sma100 := CalculateSMA(adjclose[len(high)-100+i:len(high)-i], 100)[0]
			fmt.Printf("SMA100: %v\n", sma100)
			sma200 := CalculateSMA(adjclose[len(high)-200+i:len(high)-i], 200)[0]
			fmt.Printf("SMA200: %v\n", sma200)

			conv, base, spana, spanb, lag := CalculateIchimoku(high, low, adjclose)

			//push to database with shared ID

			var f IndicatorSqlEntry = IndicatorSqlEntry{
				DatabaseID:                  dbid[i],
				Time:                        time[i],
				Ticker:                      symbol,
				CommodityChannelIndex:       cci,
				FastRsi:                     0, // replace thius somehow
				ExponentialMovingAverage200: ema200,
				ExponentialMovingAverage100: ema100,
				ExponentialMovingAverage50:  ema50,
				ExponentialMovingAverage20:  ema20,
				Macd:                        macd,
				MoneyFlowIndex:              mfi,
				SmoothedMovingAverage200:    sma200,
				SmoothedMovingAverage100:    sma100,
				SmoothedMovingAverage50:     sma50,
				SmoothedMovingAverage20:     sma20,
				Conversion:                  conv,
				Base:                        base,
				SpanA:                       spana,
				SpanB:                       spanb,
				Lagging:                     lag,
			}
			stmt, err := db.Prepare("INSERT INTO D_TA (Dbid, Time, Ticker, Macd, FastRsi, SlowRsi, ExponentialMovingAverage200, ExponentialMovingAverage100, ExponentialMovingAverage50, ExponentialMovingAverage20, SmoothedMovingAverage200, SmoothedMovingAverage100, SmoothedMovingAverage50, SmoothedMovingAverage20, Conversion, Base, SpanA, SpanB, CurrentSpanA, CurrentSpanB, CommodityChannelIndex, MoneyFlowIndex) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);")
			if err != nil {
				fmt.Println(err)
			}
			_, err = stmt.Exec(f.DatabaseID, f.Time, string(symbol), f.Macd, f.FastRsi, f.SlowRsi, f.ExponentialMovingAverage200, f.ExponentialMovingAverage100, f.ExponentialMovingAverage50, f.ExponentialMovingAverage20, f.SmoothedMovingAverage200, f.SmoothedMovingAverage100, f.SmoothedMovingAverage50, f.SmoothedMovingAverage20, f.Conversion, f.Base, f.SpanA, f.SpanB, f.CurrentSpanA, f.CurrentSpanB, f.CommodityChannelIndex, f.MoneyFlowIndex)
			if err != nil {
				fmt.Println(err)
			}
			stmt.Close()
		}
	}
}
func databaseExists() bool {
	// Open a connection to the MySQL server
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	// Ping the server to make sure the connection works
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Check if the "DailyData" database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)", "DailyData").Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return exists
}

func main() {
	/* the way that this works at the moment, we have our data already in sql data base,
	we check if we have an sql database setup
		if not send email.
	if we do
		check if we have our tables
			if not send email.
		then do the DailyCheck function to check symbols and times and pull more information.
		then we want to check if there are any entries in the day that don't have an entry in the TA analysis, and fix if they can.
	Also want to add a feature that makes this a service to run daily.
	*/
	c := databaseExists()
	if !c {
		panic("No SQL Database Configured")
		// We do have the sql
	} else {
		//
		DailyCheck()
		ScheduleChecker()
	}
}
