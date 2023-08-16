package main

// data from TA api is passed in

// instead lets pass the cache and these emas, add them to the cache, and do the analysis
func SMAFunctions(SMA20 []float64, SMA50 []float64, SMA100 []float64, SMA200 []float64, close []float64) int {

	// so backtrack is gonna be way different than how it is now
	/* we want to:
	- find any crosses with EMA's
	- find when price crosses an ema
	- determine a support or resistance level relative to price
	- determine trend by measuring ratios between EMA's

	then we aggregate out.

	*/
	score := SMAAnalysis(SMA20, SMA50, SMA100, SMA200, close)

	return score
}

// 22 points
func CCIFunctions(cciData TAapi_Response, prevData TAapi_Response) int {
	score := CCIAnalysis(cciData, prevData)
	return score
}

// 28 points
func RSIFunctions(StochBackTrack Stochrsi, prevStochData Stochrsi) int {
	// just get a list of all the rsi and get average
	// so what if we just do backtrac
	score := StochRSIBackTrackAnalysis(StochBackTrack, prevStochData)

	return score
}

// 20 points
func MFIFunctions(mfiData TAapi_Response, prevdata TAapi_Response) int {

	score := MFIAnalysis(mfiData, prevdata)
	return score
}

// total of 5-10?_ points
func CandleFunctions(cdata CandleData, pcd CandleData) int {
	var score int
	/*
		do analysis on the open/close of candle, close of candle would be current price
		do analysis on volume to add a "severity", aka high volume then high priority and add to score
		thats it
	*/
	return score
}

// 13 points
func MacdFunctions(mdata MacdData, pmd MacdData) int {
	var score int = MacdValAnalysis(mdata.Val, mdata.Val9DEma, mdata.ValHist, pmd.Val, pmd.Val9DEma, pmd.ValHist)
	return score
}
