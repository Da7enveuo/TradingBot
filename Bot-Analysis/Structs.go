package main

type Ema_return struct {
	Ema_analysis []Analysis
	Timeperiod   string
}
type Analysis struct {
	Side       string
	Streak     int
	Ema_Change float64
	// may want to add an "emcompassedtimes" and a "timeframe" in here
}

// this is what gets returned from the lambda function
type ReturnList struct {
	Minute15 ReturnJson
	Hour1    ReturnJson
}
type CciData struct {
	Value float64 `json:"value"`
}

type TopLevelJson struct {
	BTC ReturnList `json:"BTC"`
	ETH ReturnList `json:"ETH"`
}

type EmaOver struct {
	Hour1    EmaCache
	Minute15 EmaCache
}
type EmaSym struct {
	Data EmaOver
}
type MacdData struct {
	Val      float64 `json:"valueMACD"`
	Val9DEma float64 `json:"valueMACDSingal"`
	ValHist  float64 `json:"valueMACDHist"`
}

//lets say we do 200 symbols, for each symbol we need 25 * 5 * 64 bits
type EmaCache struct {
	Ema21        []float64
	Ema50        []float64
	Ema100       []float64
	Ema200       []float64
	High         []float64
	Low          []float64
	Close        []float64
	Volume       []int64
	Conversion   []float64
	Base         []float64
	SpanA        []float64
	CurrentSpanA []float64
	SpanB        []float64
	CurrentSpanB []float64
	Width        []float64
}
type ReturnJson struct {
	EMA21      TAapi_Response `json:"Ema20"`
	EMA50      TAapi_Response `json:"Ema50"`
	EMA100     TAapi_Response `json:"Ema100"`
	EMA200     TAapi_Response `json:"Ema200"`
	SMA21      TAapi_Response `json:"Sma20"`
	SMA50      TAapi_Response `json:"Sma50"`
	SMA100     TAapi_Response `json:"Sma100"`
	SMA200     TAapi_Response `json:"Sma200"`
	StochData  Stochrsi       `json:"Stoch"`
	CCIData    TAapi_Response `json:"Cci"`
	MFIData    TAapi_Response `json:"Mfi"`
	CandleData CandleData     `json:"Candle"`
	MacdData   MacdData       `json:"Macd"`
}
type CandleData struct {
	TimestampH string  `json:"timestampHuman"`
	Timestamp  string  `json:"timestamp"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     int64   `json:"volume"`
}
type TAapi_Response struct {
	Value float64
}

type Stochrsi struct {
	Fastk float64 `json:"valueFastK"`
	Fastd float64 `json:"valueFastD"`
}

type Ichimoku struct {
	Conversion float64 `json:"conversion"`
	Base       float64 `json:"base"`
	SpanA      float64 `json:"spanA"`
	SpanB      float64 `json:"spanB"`
	// the span a and b anr the ahead of time/expaned ones, the current is the cloud at the present time with close price.
	CurrentSpanA float64 `json:"CurrentspanA"`
	CurrentSpanB float64 `json:"CurrentspanB"`
}
type TAAPISend struct {
	Symbols []string
}

type Data struct {
	Data []D
}

type D struct {
	Id     string   `json:"id"`
	Result Result   `json:"result"`
	Errors []string `json:"errors"`
}

type Result struct {
	Value      float64 `json:"value,omitempty"`
	Fastk      float64 `json:"valueFastK,omitempty"`
	Fastd      float64 `json:"valueFastD,omitempty"`
	Val        float64 `json:"valueMACD,omitempty"`
	Val9DEma   float64 `json:"valueMACDSingal,omitempty"`
	ValHist    float64 `json:"valueMACDHist,omitempty"`
	TimestampH string  `json:"timestampHuman,omitempty"`
	Timestamp  string  `json:"timestamp,omitempty"`
	Open       float64 `json:"open,omitempty"`
	High       float64 `json:"high,omitempty"`
	Low        float64 `json:"low,omitempty"`
	Close      float64 `json:"close,omitempty"`
	Volume     float64 `json:"volume,omitempty"`
}

type Bnace struct {
	Symbol string
	Test   bool
}

type Ema struct {
	Value float64
}
