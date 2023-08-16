package main

// the way all of these will be structured will be a map with the minutes in the structure "MINUTEm" (ex. 1m, 2m, 3m) and the interface will either be a float64 or array of float64
// i also have decided to include a single back track
func StochRSIBackTrackAnalysis(rsiData Stochrsi, previousRsi Stochrsi) int {
	// within this backtrack analysis, we will score each individual timeframe for as far back as it goes.
	// we will then just take an average of all the scores
	var fkscore int
	// total score is just total score to get average

	//assuming that 0 is the previous event and 1 is the current event
	//also fuck fastd, honestly fastk better and more sensitive imo
	if previousRsi.Fastk < 20 && rsiData.Fastk > 20 {
		// a cross from below 20 to above 20 should result in a buy from rsi perspective
		fkscore = 15
	} else if previousRsi.Fastk < 80 && rsiData.Fastk > 80 {
		fkscore = 15
	} else if previousRsi.Fastk > 80 && rsiData.Fastk < 80 {
		// crossing downard from the 80 should result in a sell from rsi perspective
		fkscore = 0
	} else {
		// if its slightly in the buy zone, then its attractive still
		if rsiData.Fastk >= 20 && rsiData.Fastk < 30 {
			fkscore = 6
		} else if rsiData.Fastk >= 30 && rsiData.Fastk < 40 {
			fkscore = 10
		} else if rsiData.Fastk >= 40 && rsiData.Fastk < 50 {
			fkscore = 14
		} else if rsiData.Fastk >= 50 && rsiData.Fastk < 60 {
			fkscore = 18
		} else if rsiData.Fastk >= 60 && rsiData.Fastk < 80 {
			fkscore = 23
		} else if rsiData.Fastk >= 80 {
			fkscore = 27
		}
		// if its not below 20 or above 80, it must be somewhere in the middle, so will just say fuck it its 15 for now, may adjust to return adjusted score

	}

	if previousRsi.Fastk > previousRsi.Fastd && rsiData.Fastk < rsiData.Fastd {
		fkscore += 4
	} else if previousRsi.Fastk < previousRsi.Fastd && rsiData.Fastk > rsiData.Fastd {
		fkscore -= 4
	}
	return fkscore

}
