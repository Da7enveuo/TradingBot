package main

import (
	"sync"
)

// func analyzeData(indicators []ReturnJson, symbol string, Flags [3]PersistenceFlags, ws *sync.WaitGroup) error {
func analyzeData(indicators []ReturnJson, symbol string, ws *sync.WaitGroup) error {
	defer ws.Done()

	var wg *sync.WaitGroup
	// change this up a bit for sure

	// okay here is where we need to pay attention to the amount of data and how much we get from our queries

	// oka
	var d_score chan float64
	var w_score chan float64
	var m_score chan float64
	wg.Add(3)
	// 3 chans, struct {index, result} for each.
	// maybe no chans, just do the analysis, return the results, and weigh in the flags that we have registered.
	DayAnalysis(indicators, wg, d_score)
	WeekAnalysis(indicators, wg, w_score)
	MonthAnalysis(indicators, wg, m_score)

	wg.Wait()
	day_score := <-d_score
	week_score := <-w_score
	month_score := <-m_score

	buy_day := false
	sell_day := false
	buy_week := false
	sell_week := false
	buy_month := false
	sell_month := false

	if day_score > 80.0 {
		// indicating a buy on day, whcih shoudl mean within a week or two it pops
		buy_day = true
	} else if day_score < 20.0 {
		// indicating a sell on day, which should mean within a week ro two it drops
		sell_day = true
	}
	if week_score > 80.0 {
		buy_week = true
	} else if week_score < 20.0 {
		sell_week = true
	}
	if month_score > 80.0 {
		buy_month = true
	} else if month_score < 20.0 {
		sell_month = true
	}

	if sell_month || sell_week || sell_day {
		// any of these are set than fuck it
		return nil
	} else {
		if buy_month && buy_week && buy_day {
			// we want to execute a nice little buy here

		} else if buy_month && buy_week && !buy_day && !sell_day {
			// showing buy signals on week and month but not day, need to add a little bit of nuance here or something

		}
		// other than that I don't really wanna fuck with anything else, should be good with what we have with lots of adjustments
	}
	//then send data to do further analysis from here (by further analysis I mean calculating profit loss with the inclusion of the fees from vendor)
	// this is gonna be slightly different

	return nil
}

func DayAnalysis(indicators []ReturnJson, wg *sync.WaitGroup, score chan float64) {
	defer wg.Done()
	// so how many rows do I have here? How many do I want? just 2?
	mfiscore := MFIFunctions(indicators[0].MFIData, indicators[1].MFIData)
	rsiscore := RSIFunctions(indicators[0].StochData, indicators[1].StochData)
	macdscore := MacdFunctions(indicators[0].MacdData, indicators[1].MacdData)
	cciscore := CCIFunctions(indicators[0].CCIData, indicators[1].CCIData)
	/*for i := 0; i <= len(indicators); i++ {

	}*/

	// we get this from TA api, will write code above
	// put the previous data in here along with current price to compare a touch or cross of emas
	cdatascore := CandleFunctions(indicators[0].CandleData, indicators[1].CandleData)
	// loop to collect emas/smas in arrays
	var closes []float64
	var highs []float64
	var lows []float64
	var volume []int64
	var sma20 []float64
	var sma50 []float64
	var sma100 []float64
	var sma200 []float64
	for i := 0; i <= len(indicators); i++ {
		closes = append(closes, indicators[i].CandleData.Close)
		highs = append(highs, indicators[i].CandleData.High)
		lows = append(lows, indicators[i].CandleData.Low)
		volume = append(volume, indicators[i].CandleData.Volume)

		sma20 = append(sma20, indicators[i].SMA21.Value)
		sma50 = append(sma50, indicators[i].SMA50.Value)
		sma100 = append(sma100, indicators[i].SMA100.Value)
		sma200 = append(sma200, indicators[i].SMA200.Value)
	}
	hscore, lscore volscore, cscore := StdDev(highs, lows, volume, closes)
	emascore := SMAFunctions(sma20, sma50, sma100, sma200, closes)

	var total_score float64 = float64(rsiscore+cciscore+emascore+mfiscore+cdatascore+macdscore) / 1.06
	score <- total_score
}

func WeekAnalysis(indicators []ReturnJson, wg *sync.WaitGroup, score chan float64) {
	defer wg.Done()

	mfiscore := MFIFunctions(indicators[0].MFIData, indicators[1].MFIData)
	rsiscore := RSIFunctions(indicators[0].StochData, indicators[1].StochData)
	macdscore := MacdFunctions(indicators[0].MacdData, indicators[1].MacdData)
	cciscore := CCIFunctions(indicators[0].CCIData, indicators[1].CCIData)
	var closes []float64
	var highs []float64
	var lows []float64
	var volume []int64
	var sma20 []float64
	var sma50 []float64
	var sma100 []float64
	var sma200 []float64
	for i := 0; i <= len(indicators); i++ {
		closes = append(closes, indicators[i].CandleData.Close)
		highs = append(highs, indicators[i].CandleData.High)
		lows = append(lows, indicators[i].CandleData.Low)
		volume = append(volume, indicators[i].CandleData.Volume)

		sma20 = append(sma20, indicators[i].SMA21.Value)
		sma50 = append(sma50, indicators[i].SMA50.Value)
		sma100 = append(sma100, indicators[i].SMA100.Value)
		sma200 = append(sma200, indicators[i].SMA200.Value)
	}
	// we get this from TA api, will write code above
	// put the previous data in here along with current price to compare a touch or cross of emas
	cdatascore := CandleFunctions(indicators[0].CandleData, indicators[1].CandleData)
	emascore := SMAFunctions(sma20, sma50, sma100, sma200, closes)

	var total_score float64 = float64(rsiscore+cciscore+emascore+mfiscore+cdatascore+macdscore) / 1.06
	score <- total_score
}

func MonthAnalysis(indicators []ReturnJson, wg *sync.WaitGroup, score chan float64) {
	defer wg.Done()
	mfiscore := MFIFunctions(indicators[0].MFIData, indicators[1].MFIData)
	rsiscore := RSIFunctions(indicators[0].StochData, indicators[1].StochData)
	macdscore := MacdFunctions(indicators[0].MacdData, indicators[1].MacdData)
	cciscore := CCIFunctions(indicators[0].CCIData, indicators[1].CCIData)
	var closes []float64
	var highs []float64
	var lows []float64
	var volume []int64
	var sma20 []float64
	var sma50 []float64
	var sma100 []float64
	var sma200 []float64
	for i := 0; i <= len(indicators); i++ {
		closes = append(closes, indicators[i].CandleData.Close)
		highs = append(highs, indicators[i].CandleData.High)
		lows = append(lows, indicators[i].CandleData.Low)
		volume = append(volume, indicators[i].CandleData.Volume)

		sma20 = append(sma20, indicators[i].SMA21.Value)
		sma50 = append(sma50, indicators[i].SMA50.Value)
		sma100 = append(sma100, indicators[i].SMA100.Value)
		sma200 = append(sma200, indicators[i].SMA200.Value)
	}
	// we get this from TA api, will write code above
	// put the previous data in here along with current price to compare a touch or cross of emas
	cdatascore := CandleFunctions(indicators[0].CandleData, indicators[1].CandleData)
	emascore := SMAFunctions(sma20, sma50, sma100, sma200, closes)

	var total_score float64 = float64(rsiscore+cciscore+emascore+mfiscore+cdatascore+macdscore) / 1.06
	score <- total_score
}
