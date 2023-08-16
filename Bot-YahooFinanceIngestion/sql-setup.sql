
Create Database DailyData;

Use DailyData;


-- make this DailyData   
CREATE TABLE DD (
    -- links the data from yahoo to the Technical indicators for that entry
   Dbid int AUTO_INCREMENT Primary Key,
   Ticker VARCHAR(255),
   Time INT ,
   Open DOUBLE,
   High DOUBLE,
   Low DOUBLE,
   Close DOUBLE,
   -- no more adjusted close since we try to account for it
   Volume DOUBLE,
   -- this is to help determine crypto vs stock vs forex vs commodity
   Type VARCHAR(255)
);

-- make this DailyStockTechnicals
CREATE TABLE D_TA (
    Dbid int,
    Ticker VARCHAR(255),
    Time INT,
    Macd FLOAT,
    MacdHist FLOAT,
    MacdSignal FLOAT,
    FastRsi FLOAT,
    SlowRsi FLOAT,
    ExponentialMovingAverage200 FLOAT,
    ExponentialMovingAverage100 FLOAT,
    ExponentialMovingAverage50 FLOAT,
    ExponentialMovingAverage20 FLOAT,
    SmoothedMovingAverage200 FLOAT,
    SmoothedMovingAverage100 FLOAT,
    SmoothedMovingAverage50 FLOAT,
    SmoothedMovingAverage20 FLOAT,
    Conversion FLOAT,
    Base FLOAT,
    SpanA FLOAT,
    SpanB FLOAT,
    CurrentSpanA FLOAT,
    CurrentSpanB FLOAT,
    CommodityChannelIndex FLOAT,
    MoneyFlowIndex FLOAT,
    Type VARCHAR(255),
    FOREIGN KEY (Dbid) REFERENCES DD(Dbid)
);


