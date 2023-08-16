package main

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
)

type IndicatorSqlEntry struct {
	Time       int64
	Ticker     string
	DatabaseID int
	// macd
	Macd       float64
	MacdHist   float64
	MacdSignal float64

	// Regular RSI
	Rsi float64
	// stochastic rsi
	FastRsi float64
	SlowRsi float64
	// Exponential Moving Averages
	ExponentialMovingAverage200 float64
	ExponentialMovingAverage100 float64
	ExponentialMovingAverage50  float64
	ExponentialMovingAverage20  float64

	// Smoothed Moving Averages
	SmoothedMovingAverage200 float64
	SmoothedMovingAverage100 float64
	SmoothedMovingAverage50  float64
	SmoothedMovingAverage20  float64
	// ichimoku cloud
	Conversion   float64
	Base         float64
	SpanA        float64
	SpanB        float64
	Lagging      float64
	CurrentSpanA float64
	CurrentSpanB float64

	// CCI, 20 based
	CommodityChannelIndex float64

	// MFI
	MoneyFlowIndex float64
	//
}

// want this to be 14 periods by default then?
func CalculateMoneyFlowIndex(closePrices []float64, highPrices []float64, lowPrices []float64, volume []int64, period int) float64 {

	if len(closePrices) != len(highPrices) || len(closePrices) != len(lowPrices) || len(closePrices) != len(volume) {
		panic("Input slices must have the same length")
	}
	var positive_money_flows float64
	var negative_money_flows float64
	var previous_flow float64
	fmt.Printf("Close Length: %v\n", len(closePrices))
	fmt.Printf("Close Length: %v\n", len(highPrices))
	fmt.Printf("Close Length: %v\n", len(lowPrices))
	for i := 0; i <= len(closePrices)-1; i++ {
		if i == 0 {
			previous_flow = (closePrices[0] + highPrices[0] + lowPrices[0]) / 3
		}
		high := highPrices[i]
		low := lowPrices[i]
		close := closePrices[i]
		volume := volume[i]
		typical_price := (high + low + close) / 3
		money_flow := typical_price * float64(volume)

		if money_flow > previous_flow {
			positive_money_flows += money_flow
		} else if money_flow < previous_flow {
			negative_money_flows += math.Abs(money_flow)
		}

		previous_flow = money_flow
	}
	if negative_money_flows == 0 {
		// Avoid division by zero
		return 100.0
	}
	// separate to 2 groups, positive/negative, money flow is the sum of meony flow values for all periods where typical price increased, negative is sum of all where typical price decreased

	money_ratio := positive_money_flows / negative_money_flows
	mfi := 100 - (100 / (1 + money_ratio))
	return mfi
}

// hard code period to 20?

// typical price =  high, low, and close  / 3
// calc tpical price for all periods
// calc moving average (just ma) of typical prices
// calculate the mean deviation by subtractign each typical price form moving average and summing absolute values and dividing by 20
// sum absolute values of
func CalculateCommodityChannelIndex(high []float64, low []float64, close []float64) float64 {

	typicalPrices := make([]float64, len(high))
	for i := 0; i <= len(high)-1; i++ {
		typicalPrices[i] = (high[i] + low[i] + close[i]) / 3.0
	}

	// Step 2: Calculate the Simple Moving Average (SMA) of TP-
	sma := CalculateMA(typicalPrices)

	// Step 3: Calculate the Mean Deviation (MD)
	md := CalculateMD(typicalPrices, sma)

	// Step 4: Calculate the Commodity Channel Index (CCI)
	var multiplier float64 = 0.015
	cci := (typicalPrices[len(typicalPrices)-1] - sma) / (multiplier * md)

	return cci
}

// 14 day period here
// Okay so first we calculate the rsi of the periods
// then we identify the highest and lowest rsi reading for each period.
//  stoch rsi formula: RSI - min[RSI] / max[RSI] - min[RSI] where RSI is the current reading and min/max[RSI] is the lowest/highest rsi reading in the period.
// so it is 100 - (100/1+(average gain%/14)/(average loss%/14))
// we should be passing in 15 days

// 1. get price change from prev period to now
// 2. get average loss/gain from the closes within the period we calculating, which is the current day - previous day / 14?
// 3. calc RS by dividing average gain / average loss
// 4. calc RSI by adding stochastic oscillator 100 - (100 / (1 + RS))

// 15th close price is the previous close, which means we should check that the array is 15 in length and make sure to pass in 15 floats
// so maybe we leave the rsi calculations, and just do stochastic off that?

func CalcRsi(close_prices []float64) float64 {
	var average_gain float64
	var average_loss float64
	var rsi float64
	// 2. get average gain/loss
	var pr_close float64 = 0
	var gains float64 = 0
	var losses float64 = 0
	for i, val := range close_prices {
		if i != 0 {
			if val > pr_close {
				gains += val - pr_close
			} else if val < pr_close {
				losses += pr_close - val
			}
		}
		pr_close = val
	}
	if gains != 0 {
		average_gain = gains / 14
	} else {
		average_gain = 0
	}
	if losses != 0 {
		average_loss = losses / 14
	} else {
		average_loss = 0
	}
	rs := average_gain / average_loss
	// now we go through the rsi shit and
	rsi = 100 - float64(100/1+rs)
	// gotta return this bitch after stochasticizing

	// return stoch_rsi, average_gain, average_loss
	return rsi
}

// here I think I will need to start this once we have 14 RSI periods available to work with.
func CalcStochasticRsi(time int64, symbol string, db *sql.DB) {
	// okay I want to have this setup so that we are pulling generate stoch rsi
	// we need 14 completed rsi periods, where we determine the rsi high and lows, then use the formula:
	// currentRsi - min14PeriodRsi / max14PeriodRsi - min14PeriodRsi
	rd, err := db.Query(fmt.Sprintf("SELECT Rsi FROM DD WHERE Ticker = '%v' AND Time < '%v' LIMIT 14 DESC", strings.ToLower(symbol), time))

	if err != nil {
		fmt.Println(err)
	}
	var da []IndicatorSqlEntry
	if rd.Next() {
		for rd.Next() {
			var d IndicatorSqlEntry
			err = rd.Scan(&d.Rsi)
			if err != nil {
				fmt.Println(err)
			}
			da = append(da, d)
		}
	}
	rd.Close()
	// now we have da, which sould have all the sql entries we need, may need to adjust the query
	var low_rsi float64 = 0
	var high_rsi float64 = 0
	for _, sql_entry := range da {
		if sql_entry.Rsi > high_rsi {
			high_rsi = sql_entry.Rsi
		}
		if sql_entry.Rsi < low_rsi {
			low_rsi = sql_entry.Rsi
		}
	}
	c_rsi := da[len(da)+1].Rsi

	stoch_rsi := (c_rsi - low_rsi) / (high_rsi - low_rsi)
	_, err = db.Query(fmt.Sprintf("INSERT INTO D_TA WHERE Dbid = %v (FastRsi) VALUES (%v)", stoch_rsi, da[len(da)+1].DatabaseID))
	if err != nil {
		panic(err)
	}
}

// implement this appropriately
func CalcSlowRSI(time int64, ticker string, db *sql.DB) {
	rd, err := db.Query(fmt.Sprintf("SELECT FastRsi FROM D_TA WHERE Ticker = '%v' AND Time < '%v' LIMIT 9 DESC", strings.ToLower(ticker), time))

	if err != nil {
		fmt.Println(err)
	}
	var da []IndicatorSqlEntry
	if rd.Next() {
		for rd.Next() {
			var d IndicatorSqlEntry
			err = rd.Scan(&d.FastRsi)
			if err != nil {
				fmt.Println(err)
			}
			da = append(da, d)
		}
	}
	var stoch_rsis []float64
	for _, val := range da {
		stoch_rsis = append(stoch_rsis, val.FastRsi)
	}
	slow_rsi := CalculateMA(stoch_rsis)
	_, err = db.Query(fmt.Sprintf("INSERT %v INTO SlowRSI WHERE Dbid = %v", slow_rsi, da[len(da)+1].DatabaseID))
	if err != nil {
		panic(err)
	}
}

/*
CCI=
.015×Mean Deviation
Typical Price−MA
​*/
// subtract typical price from movinga average
func CalculateMD(typical_prices []float64, ma float64) float64 {
	var sum float64
	for i := 0; i <= len(typical_prices)-1; i++ {
		sum += math.Abs(typical_prices[i]) - ma
	}
	return sum / 20
}
func CalculateMA(typical_prices []float64) float64 {
	var sum float64
	for i := 0; i <= len(typical_prices)-1; i++ {
		sum += typical_prices[i]
	}
	return sum / float64(len(typical_prices))
}
func CalculateMACD(closingPrices []float64, timePeriod int) float64 {
	// calculate the 12 day EMA
	// calculate the 26 day EMA
	// calculate the macd by subtracting the 26 day EMA from the 12 day EMA to get macd
	// calculate the 9 day moving average based on the macd
	ema12 := CalculateEMA(closingPrices[len(closingPrices)-12:], 12)
	ema26 := CalculateEMA(closingPrices[len(closingPrices)-26:], 26)
	macd := ema12[len(ema12)-1] - ema26[len(ema26)-1]

	return macd
}
func CalculateSMA(closingPrices []float64, timePeriod int) []float64 {
	// Validate input parameters
	if len(closingPrices) < timePeriod {
		panic("Insufficient data for the given time period.")
	}
	if timePeriod <= 0 {
		panic("Invalid time period. Time period should be greater than 0.")
	}

	// Initialize the SMA array
	sma := make([]float64, len(closingPrices)-timePeriod+1)

	// Calculate the SMAs
	for i := timePeriod; i <= len(closingPrices); i++ {
		sum := 0.0
		for j := i - timePeriod; j < i; j++ {
			sum += closingPrices[j]
		}
		sma[i-timePeriod] = sum / float64(timePeriod)
	}

	return sma
}

func CalculateEMA(closingPrices []float64, timePeriod int) []float64 {
	// Validate input parameters
	if len(closingPrices) < timePeriod {
		panic("Insufficient data for the given time period.")
	}
	if timePeriod <= 0 {
		panic("Invalid time period. Time period should be greater than 0.")
	}

	// Calculate the smoothing constant
	k := 2.0 / float64(timePeriod+1)

	// Initialize the EMA array
	ema := make([]float64, len(closingPrices)-timePeriod+1)

	// Calculate the first EMA using the simple moving average
	sma := CalculateSMA(closingPrices[:timePeriod], timePeriod)
	ema[0] = sma[len(sma)-1]

	// Calculate the subsequent EMAs using the previous EMA and the current closing price
	for i := timePeriod; i < len(closingPrices); i++ {
		ema[i-timePeriod+1] = k*(closingPrices[i]-ema[i-timePeriod]) + ema[i-timePeriod]
	}

	return ema
}

func CalculateIchimoku(highs []float64, lows []float64, closes []float64) (float64, float64, float64, float64, float64) {
	// for the "current" span a/b we should just get span a from 14 time frames ago.

	// tenkan and kijun uses the same function but tenkan uses 9 periods of data, whereas kijun uses 26 periods of dataeww
	tenkan := CalcTenkan(highs[:8], lows[:8])  // Conversion
	kijun := CalcTenkan(highs[:25], lows[:25]) // Base
	spana := tenkan + kijun/2                  // Span A
	spanb := CalcTenkan(highs[:56], lows[:56]) // Span B
	chikou := closes[9]                        // lagging span

	return tenkan, kijun, spana, spanb, chikou
}

func CalcTenkan(highs []float64, lows []float64) float64 {
	var high float64
	var low float64
	for i := 0; i <= len(highs)-1; i++ {
		if highs[i] > high {
			high = highs[i]
		}
		if lows[i] < low {
			low = lows[i]
		}
	}
	/*
		The Tenkan-sen is calculated by adding the highest high and the lowest low over a specific period and dividing the result by two. The formula is as follows:
		Tenkan-sen = (Highest High + Lowest Low) / 2
		Typically, the period used for the calculation is the past 9 periods.
	*/
	return (high + low) / 2
}
