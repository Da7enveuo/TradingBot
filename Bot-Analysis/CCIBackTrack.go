package main

//
func CCIAnalysis(cci TAapi_Response, prevcci TAapi_Response) int {
	// we will go through and take average? just trend comparison? trend comparision would be the line of best fit?of
	// go through and compare which points are higher and which points are lower, and how many are consecutively higher/lower

	var fkscore int
	// total score is just total score to get average

	//assuming that 0 is the prevcciious event and 1 is the current event
	//also fuck fastd, honestly fastk better and more sensitive imo
	if prevcci.Value < -100 && cci.Value > -100 {
		// a cross from below 20 to above 20 should result in a buy from rsi perspective
		fkscore = 30
	} else if prevcci.Value > 100 && cci.Value < 100 {
		// crossing downard from the 80 should result in a sell from rsi perspective
		fkscore = 0
	} else {
		if cci.Value < prevcci.Value {
			fkscore = -5
		}

		// if its slightly in the buy zone, then its attractive still
		if cci.Value >= -100 && cci.Value < -60 {
			fkscore = 8
		} else if cci.Value >= -60 && cci.Value < -20 {
			fkscore = 14
		} else if cci.Value >= -20 && cci.Value < 20 {
			fkscore = 18
		} else if cci.Value >= 20 && cci.Value < 60 {
			fkscore = 23
		} else if cci.Value >= 100 {
			fkscore = 27
		}
		// if its not below 20 or above 80, it must be somewhere in the middle, so will just say fuck it its 15 for now, may adjust to return adjusted score

	}
	return fkscore
}
