package main

import "math"

/* we want to:
- find any crosses with sma's
- find when price crosses an sma
- determine a support or resistance level relative to price
- determine trend by measuring ratios between sma's

we need to exaggerate our levels for stop losses and such

2.5x atr for stop losses
then we aggregate out.
*/

// total points:
func SMAAnalysis(sma21 []float64, sma50 []float64, sma100 []float64, sma200 []float64, close []float64) int {
	// get the trend score, 3 points total
	tscore := DetermineTrend(sma21, sma50, sma100, sma200)
	// determine ratios between the sma levels. We will also use this as apart of s/r levels
	// this gets the ratio score and the prediction score, as well as support/resistance that gets wrapped into ratio score
	// 6 for ratio_score, 3 for pcscore, 3 for pscore
	em_score, pre_score := DetermineSMARatios(sma21, sma50, sma100, sma200, close)
	/*
		so I want to get:
			- ratios comparing each sma with their longer/shorter counterpart (gold/death cross)
			- ratios comparing each sma to close price
			- trend of the ratios
	*/
	// return the determined score from sma analysis on previous periods.
	return tscore + em_score + pre_score
}

// this just adds up the last 20 timeframes to see which side it is leaning towards for each sma
// we think of the longer timeframe sma's as more reliable and will score them as such.
func DetermineTrend(sma21 []float64, sma50 []float64, sma100 []float64, sma200 []float64) int {
	var score int

	// this is going to be looking at each sma, determining if it is going up or down, then aggregating out

	var sma100_score int
	var sma200_score int
	var sma50_score int
	var sma21_score int
	// for loop to go through the sma cache
	for i := 0; i <= len(sma200)-1; i++ {
		// we go backwards from the first timeframe and check linearly with respect to time
		if i == 0 {
		} else {
			sma200_score += CheckPrevsma(sma200[i], sma200[i-1])
			sma100_score += CheckPrevsma(sma100[i], sma100[i-1])
			sma50_score += CheckPrevsma(sma50[i], sma50[i-1])
			sma21_score += CheckPrevsma(sma21[i], sma21[i-1])
		}
	}
	// each "score" instance will return a number between -3 - 3
	score += FindScore(sma200_score)
	score += FindScore(sma100_score)
	score += FindScore(sma50_score)
	score += FindScore(sma21_score)
	// which results in a total score ranging from -12 - 12
	// do one last abstraction, if between -12 - -4 then 0, else if -3 - 5 then 1,
	if score >= -12 && score <= -6 {
		score = 0
	} else if score > -6 && score <= 0 {
		score = 1
	} else if score > 0 && score <= 6 {
		score = 2
	} else if score > 6 {
		score = 3
	}
	// essentially a scalar value to assess the trend score of all these moving averages.
	return score
}

func CheckPrevsma(current float64, previous float64) int {
	var rscore int
	if previous > current {
		rscore -= 1

	} else if previous < current {
		rscore += 1
	}
	return rscore
}

func FindScore(score int) int {
	var rscore int
	if score > 0 {
		// means positive overall trend
		if score < 5 {
			rscore += 0
		} else if score > 5 && score < 10 {
			rscore += 1
		} else if score > 10 {
			rscore += 2
		}
		rscore += 1
	} else if score < 0 {
		// means negative overall trend
		rscore -= 3
	} else {
		rscore += 0
	}
	return rscore
}

func FindCross(current_sma float64, current_close float64, previous_sma float64, previous_close float64) (int, int) {
	var crosses int
	var side int

	if current_sma > current_close && previous_sma < previous_close {
		crosses = 1
		side = 1
	} else if current_sma < current_close && previous_sma > previous_close {
		crosses = 1
		side = -1
	} else {
		crosses = 0
		side = 0
	}
	return side, crosses
}
func ChecksmaCross(score int, ccount int) int {
	var score_ int
	if score > 0 {
		score_ = 2
		if ccount > 1 {
			// lots of crosses
			score_ += 2
		} else if ccount == 1 {
			// means one upward cross
			score_ += 2
		}
	} else if score < 0 {
		score_ = -2
	} else {
		score_ = 0
	}

	return score_
}

// determine support/resistance levels
func DetermineSMARatios(sma21 []float64, sma50 []float64, sma100 []float64, sma200 []float64, close []float64) (int, int) {
	var score int

	var sma200score int
	var crosses200 int
	var sma100score int
	var crosses100 int
	var sma50score int
	var crosses50 int
	var sma21score int
	var crosses21 int

	var sum200 float64
	var sum100 float64
	var sum50 float64
	var sum21 float64

	// basically just check to see crosses on smas (intra-sma as well), do a trend score on front and back halves, and do a single gaussian distribution w/unit variance to determine probabilities of price ranges
	for i := 0; i <= len(sma21); i++ {
		// check if price crossed
		if i == 0 {

		} else {
			em200 := sma200[i] / close[i]
			em100 := sma100[i] / close[i]
			em50 := sma50[i] / sma50[i]
			em21 := sma21[i] / sma21[i]
			/*
				prevem200 := smaCache.sma200[i] / smaCache.Close[i]
				prevem100 := smaCache.sma100[i] / smaCache.Close[i]
				prevem50 := smaCache.sma50[i] / smaCache.sma50[i]
				prevem21 := smaCache.sma21[i] / smaCache.sma21[i]

				em200score := CompsmaRatios(em200, prevem200)
				em100score := CompsmaRatios(em100, prevem100)
				em50score := CompsmaRatios(em50, prevem50)
				em21score := CompsmaRatios(em21, prevem21)
			*/
			// okay so em score is just comparing current to previous sma. What is the goal.
			// want to determine if middle ratio averages are smaller than ennd ratio averages.
			// func to check shit
			sum200 += em200
			sum100 += em100
			sum50 += em50
			sum21 += em21

			// determine the support/resistance
			// how to make persistant? keep an array of all ratios? Or just back a certain extended timeframe
			//
			f, v := FindCross(sma200[i], close[i], sma200[i-1], close[i-1])
			sma200score += f
			crosses200 += v
			//100
			f, v = FindCross(sma100[i], close[i], sma100[i-1], close[i-1])
			sma100score += f
			crosses100 += v
			//50
			f, v = FindCross(sma50[i], close[i], sma50[i-1], close[i-1])
			sma50score += f
			crosses50 += v
			//21
			f, v = FindCross(sma21[i], close[i], sma21[i-1], close[i-1])
			sma21score += f
			crosses21 += v
		}
	}
	// this is assessing the sum of the sma ratios to close price,
	assScore200 := Assess(sum200)
	assScore100 := Assess(sum100)
	assScore50 := Assess(sum50)
	assScore21 := Assess(sum21)
	assScore := assScore100 + assScore200 + assScore50 + assScore21
	// this is calculating line of best fit and estimating next moving averages
	// find up/down trend with each
	em200 := ChecksmaCross(sma200score, crosses200)
	em100 := ChecksmaCross(sma100score, crosses100)
	em50 := ChecksmaCross(sma50score, crosses50)
	em21 := ChecksmaCross(sma21score, crosses21)
	/*
		I think it might be better to take an average of the front end of the curve and the back end, as well as overall then compare
		this way we can get more specific to the timeperiod relevant to us and not lose dat

		do multivariate gaussian distribiution
		1/sqrt(abs(2piK)) * exp (-1/2 N sigma u,v=1  z vu(K ^-1)fvuv Zv)

		lets see if we can't shove all the sma prices with a cumulative distribution function instead of PDF
		Kuv is an N by N symmetric positive definite matrix and inverse is k-1uv
		z is the vector in the data set?
	*/
	//pscore := FindPrediction(arr200, arr100, arr50, arr21, close)

	score = assScore
	pcscore := em200 + em100 + em50 + em21
	//
	return score, pcscore
}
func Assess(score float64) int {
	avg := score / 21
	var rscore int
	if avg > 1 {
		if avg <= 1.005 {
			// small average, meaning that the close price is trending close below the average
			rscore = 3
		} else if avg > 1.005 && avg < 1.05 {
			rscore = 2
		} else {
			rscore = 0
		}
	} else if avg < 1 {
		if avg >= .995 {
			// means average is trending
			rscore = 0
		} else if avg < .995 && avg > .95 {
			rscore = 2
		} else {
			rscore = 0
		}
	}
	return rscore
}

///  CHECK ALL THE MUTIVARIABTE GAUSSIAN FUNCTIONS TO ENSURE THEY ARE FUNCTIONING PROPERLY
func multivariateGaussian(x1, x2, x3, x4, y []float64) float64 {
	// Calculate the mean vector and covariance matrix of the input data
	mean := []float64{mean(x1), mean(x2), mean(x3), mean(x4)}
	covariance := covariance(x1, x2, x3, x4)

	// Calculate the determinant and inverse of the covariance matrix
	det := determinant(covariance)
	inverse := inverse(covariance)

	// Calculate the value of the multivariate Gaussian distribution at the input point (x1, x2, x3, x4, y)
	x := []float64{x1[len(x1)-1], x2[len(x2)-1], x3[len(x3)-1], x4[len(x4)-1], y[len(y)-1]}
	exponent := -0.5 * mahalanobisDistance(x, mean, inverse)
	coefficient := 1 / (math.Pow(2*math.Pi, 2) * math.Sqrt(det))
	return coefficient * math.Exp(exponent)
}

func mean(arr []float64) float64 {
	sum := 0.0
	for _, x := range arr {
		sum += x
	}
	return sum / float64(len(arr))
}

func covariance(x1, x2, x3, x4 []float64) [][]float64 {
	n := len(x1)
	cov := make([][]float64, 4)
	for i := range cov {
		cov[i] = make([]float64, 4)
	}
	for i := 0; i < n; i++ {
		cov[0][0] += (x1[i] - mean(x1)) * (x1[i] - mean(x1))
		cov[0][1] += (x1[i] - mean(x1)) * (x2[i] - mean(x2))
		cov[0][2] += (x1[i] - mean(x1)) * (x3[i] - mean(x3))
		cov[0][3] += (x1[i] - mean(x1)) * (x4[i] - mean(x4))
		cov[1][0] += (x2[i] - mean(x2)) * (x1[i] - mean(x1))
		cov[1][1] += (x2[i] - mean(x2)) * (x2[i] - mean(x2))
		cov[1][2] += (x2[i] - mean(x2)) * (x3[i] - mean(x3))
		cov[1][3] += (x2[i] - mean(x2)) * (x4[i] - mean(x4))
		cov[2][0] += (x3[i] - mean(x3)) * (x1[i] - mean(x1))
		cov[2][1] += (x3[i] - mean(x3)) * (x2[i] - mean(x2))
		cov[2][2] += (x3[i] - mean(x3)) * (x3[i] - mean(x3))
		cov[2][3] += (x3[i] - mean(x3)) * (x4[i] - mean(x4))
		cov[3][0] += (x4[i] - mean(x4)) * (x1[i] - mean(x1))
		cov[3][1] += (x4[i] - mean(x4)) * (x2[i] - mean(x2))
		cov[3][2] += (x4[i] - mean(x4)) * (x3[i] - mean(x3))
		cov[3][3] += (x4[i] - mean(x4)) * (x4[i] - mean(x4))
	}
	for i := range cov {
		for j := range cov[i] {
			cov[i][j] /= float64(n - 1)
		}
	}
	return cov
}

func determinant(matrix [][]float64) float64 {
	if len(matrix) == 1 {
		return matrix[0][0]
	}
	if len(matrix) == 2 {
		return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
	}
	var det float64
	for i := range matrix[0] {
		minor := make([][]float64, len(matrix)-1)
		for j := range minor {
			minor[j] = make([]float64, len(matrix[0])-1)
		}
		for j := 1; j < len(matrix); j++ {
			for k := range matrix[j] {
				if k != i {
					var kay int
					if k > i {
						kay = k
					} else {
						kay = i
					}
					minor[j-1][k-kay] = matrix[j][k]
				}
			}
		}
		if i%2 == 0 {
			det += matrix[0][i] * determinant(minor)
		} else {
			det -= matrix[0][i] * determinant(minor)
		}
	}
	return det
}

func inverse(matrix [][]float64) [][]float64 {
	det := determinant(matrix)
	if det == 0 {
		panic("Matrix is singular")
	}
	if len(matrix) == 1 {
		return [][]float64{{1 / matrix[0][0]}}
	}
	adjugate := make([][]float64, len(matrix))
	for i := range adjugate {
		adjugate[i] = make([]float64, len(matrix[0]))
	}
	for i := range matrix {
		for j := range matrix[0] {
			minor := make([][]float64, len(matrix)-1)
			for k := range minor {
				minor[k] = make([]float64, len(matrix[0])-1)
			}
			for k := range matrix {
				if k != i {
					for l := range matrix[0] {
						if l != j {
							var kay int
							if k > i {
								kay = k
							} else {
								kay = i
							}
							var el int
							if l > j {
								el = l
							} else {
								el = j
							}
							minor[k-kay][l-el] = matrix[k][l]
						}
					}
				}
			}
			if (i+j)%2 == 0 {
				adjugate[j][i] = determinant(minor)
			} else {
				adjugate[j][i] = -determinant(minor)
			}
		}
	}
	for i := range adjugate {
		for j := range adjugate[i] {
			adjugate[i][j] /= det
		}
	}
	return adjugate
}

func mahalanobisDistance(x, mean []float64, inverse [][]float64) float64 {
	var distance float64
	for i := range x {
		distance += (x[i] - mean[i]) * inverse[i][0] * (x[0] - mean[0])
		for j := 1; j < len(x); j++ {
			distance += (x[i] - mean[i]) * inverse[i][j] * (x[j] - mean[j])
		}
	}
	return math.Sqrt(distance)
}
func FindPrediction(em200, em100, em50, em21, close []float64) int {
	// multivariate gaussian distribution
	var prediction float64 = multivariateGaussian(em200, em100, em50, em21, close)

	// misleading, we want to return a value
	if prediction > close[0] {
		return 3
	} else if prediction < close[0] {
		return 0
	} else {
		return 1
	}

}

// need to adjust this
func CompsmaRatios(current float64, previous float64) int {
	var rscore int

	if current > previous {
		rscore = 1
	} else if current < previous {
		rscore = 0
	}

	return rscore
}

/*
type Ct struct {
	sma200 int
	sma100 int
	sma50  int
	sma21  int
}

func popAndPush(arr []float64, val float64) []float64 {
	// Remove last element of array
	for i := len(arr) - 1; i > 0; i-- {
		arr[i] = arr[i-1]
	}

	// Insert new value at the beginning of array
	arr[0] = val

	return arr
}
*/
