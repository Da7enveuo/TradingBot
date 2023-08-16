package main

// 13 poitns total
func MacdValAnalysis(macd float64, macdhist float64, macdsignal float64, prevmacd float64, prevmacdsignal float64, prevmacdhist float64) int {
	var score int
	// 4 points
	if macd < prevmacd {
		score -= 4
	} else if macd > prevmacd {
		score += 4
	}
	// 3 points
	if macdhist < prevmacdhist {
		score -= 2
	} else if macdhist > prevmacdhist {
		score += 2
	}
	// 2 points, slowest reacting
	if macdsignal < prevmacdsignal {
		score -= 2
	} else if macdsignal > prevmacdsignal {
		score += 2
	}
	// if
	// 6 points taken already
	// assign this 5 points
	if macd > macdsignal && prevmacd < prevmacdsignal {
		score += 5
	} else if macd < macdsignal && prevmacd > prevmacdsignal {
		score -= 5
	} else {
		if macd > 0 {
			score += 2
		} else if macd < 0 {
			score -= 2
		} else {
			score += 0
		}

	}

	// and now we measure the macd itself, if it is coming close and is negative, bearish, no points

	// okay so now we need analysis to determine shit
	/*
		essentially range the macd and compare with sign, and add in the hist as well, max points when they cross
			basically when macdHist (or the bar chart on Trading View) crosses upwards, thats a buy signal,
				when it crosses downwards thats a sell signal

			when macd (the blue line in trading view) crosses the red line (signal, 9day ema OF MACD)
	*/
	return score
}
