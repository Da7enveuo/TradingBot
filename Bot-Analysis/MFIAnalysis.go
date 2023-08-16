package main

func MFIAnalysis(mfiData TAapi_Response, prevdata TAapi_Response) int {

	// total score is just total score to get average

	var fkscore int
	// total score is just total score to get average

	//assuming that 0 is the previous event and 1 is the current events
	//also fuck fastd, honestly fastk better and more sensitive imo
	if prevdata.Value < 20 && mfiData.Value > 20 {
		// a cross from below 20 to above 20 should result in a buy from rsi perspective
		fkscore = 30
	} else if prevdata.Value > 80 && mfiData.Value < 80 {
		// crossing downard from the 80 should result in a sell from mfi perspective
		fkscore = 0
	} else {
		if mfiData.Value < prevdata.Value {
			fkscore = -5
		}

		// if its slightly in the buy zone, then its attractive still
		if mfiData.Value >= 20 && mfiData.Value < 30 {
			fkscore = 6
		} else if mfiData.Value >= 30 && mfiData.Value < 40 {
			fkscore = 10
		} else if mfiData.Value >= 40 && mfiData.Value < 50 {
			fkscore = 14
		} else if mfiData.Value >= 50 && mfiData.Value < 60 {
			fkscore = 18
		} else if mfiData.Value >= 60 && mfiData.Value < 80 {
			fkscore = 23
		} else if mfiData.Value >= 80 {
			fkscore = 27
		}
		// if its not below 20 or above 80, it must be somewhere in the middle, so will just say fuck it its 15 for now, may adjust to return adjusted score

	}
	return fkscore

}
