package gobithumb

import (
	"strconv"
	"strings"
	"time"
)

//==============================TICKER SETTING======================================

type Ticker struct {
	OpeningPrice     float64
	ClosingPrice     float64
	MinPrice         float64
	MaxPrice         float64
	UnitsTraded      float64
	AccTradeValue    float64
	PrevClosingPrice float64
	UnitsTraded24H   float64
	AccTradeValue24H float64
	Fluctate24H      float64
	FluctateRate24H  float64
}

func newTicker(rawTicker map[string]interface{}) Ticker {
	newTicker := Ticker{}
	newTicker.OpeningPrice, _ = strconv.ParseFloat(rawTicker["opening_price"].(string), 64)
	newTicker.ClosingPrice, _ = strconv.ParseFloat(rawTicker["closing_price"].(string), 64)
	newTicker.MinPrice, _ = strconv.ParseFloat(rawTicker["min_price"].(string), 64)
	newTicker.MaxPrice, _ = strconv.ParseFloat(rawTicker["max_price"].(string), 64)
	newTicker.UnitsTraded, _ = strconv.ParseFloat(rawTicker["units_traded"].(string), 64)
	newTicker.AccTradeValue, _ = strconv.ParseFloat(rawTicker["acc_trade_value"].(string), 64)
	newTicker.PrevClosingPrice, _ = strconv.ParseFloat(rawTicker["prev_closing_price"].(string), 64)
	newTicker.UnitsTraded24H, _ = strconv.ParseFloat(rawTicker["units_traded_24H"].(string), 64)
	newTicker.AccTradeValue24H, _ = strconv.ParseFloat(rawTicker["acc_trade_value_24H"].(string), 64)
	newTicker.Fluctate24H, _ = strconv.ParseFloat(rawTicker["fluctate_24H"].(string), 64)
	newTicker.FluctateRate24H, _ = strconv.ParseFloat(rawTicker["fluctate_rate_24H"].(string), 64)

	return newTicker
}

//==============================ORDERBOOK SETTING======================================

type Bidask struct {
	Price    float64
	Quantity float64
}

type Orderbook struct {
	Bids []Bidask
	Asks []Bidask
}

func newOrderbook(rawOrderbook map[string]interface{}) Orderbook {
	newOrderbook := Orderbook{}

	bids := rawOrderbook["bids"].([]interface{})
	asks := rawOrderbook["asks"].([]interface{})

	newOrderbook.Bids = make([]Bidask, len(bids))
	newOrderbook.Asks = make([]Bidask, len(asks))
	for index, data := range bids {
		oneBid := data.(map[string]interface{})
		newOrderbook.Bids[index].Price, _ = strconv.ParseFloat(oneBid["price"].(string), 64)
		newOrderbook.Bids[index].Quantity, _ = strconv.ParseFloat(oneBid["quantity"].(string), 64)
	}
	for index, data := range asks {
		oneAsk := data.(map[string]interface{})
		newOrderbook.Asks[index].Price, _ = strconv.ParseFloat(oneAsk["price"].(string), 64)
		newOrderbook.Asks[index].Quantity, _ = strconv.ParseFloat(oneAsk["quantity"].(string), 64)
	}
	return newOrderbook
}

//==============================TRANSACTIONHISOTRY SETTING======================================

type OneTransaction struct {
	TransactionDate time.Time
	Type            string
	UnitsTraded     float64
	Price           float64
	Total           float64
}

const trTimeForm = "2006-01-02 15:04:05"

func newTransactionHistory(rawTransactionHistory []interface{}) []OneTransaction {
	result := make([]OneTransaction, len(rawTransactionHistory))
	for index, data := range rawTransactionHistory {
		dataMap := data.(map[string]interface{})
		result[index].UnitsTraded, _ = strconv.ParseFloat(dataMap["units_traded"].(string), 64)
		result[index].Price, _ = strconv.ParseFloat(dataMap["price"].(string), 64)
		result[index].Total, _ = strconv.ParseFloat(dataMap["total"].(string), 64)
		result[index].TransactionDate, _ = time.Parse(trTimeForm, dataMap["transaction_date"].(string))
		result[index].Type, _ = dataMap["type"].(string)
	}
	return result
}

//==============================ASSETSSTATUS SETTING======================================

//==============================BTCI SETTING======================================

type BT_I struct {
	MarketIndex float64
	Rate        float64
	Width       float64
}

type BTCI struct {
	BTAI BT_I
	BTMI BT_I
}

func newBTCI(rawBTAI, rawBTMI interface{}) BTCI {
	btai := rawBTAI.(map[string]interface{})
	btmi := rawBTMI.(map[string]interface{})

	newBTCI := BTCI{}
	newBTCI.BTAI.MarketIndex, _ = strconv.ParseFloat(btai["market_index"].(string), 64)
	newBTCI.BTAI.Rate, _ = strconv.ParseFloat(btai["rate"].(string), 64)
	newBTCI.BTAI.Width, _ = strconv.ParseFloat(btai["width"].(string), 64)
	newBTCI.BTMI.MarketIndex, _ = strconv.ParseFloat(btmi["market_index"].(string), 64)
	newBTCI.BTMI.Rate, _ = strconv.ParseFloat(btmi["rate"].(string), 64)
	newBTCI.BTMI.Width, _ = strconv.ParseFloat(btmi["width"].(string), 64)
	return newBTCI
}

//==============================CANDLESTICK SETTING======================================

type RawCandleStick struct {
	Status  int    `json:"status,string"`
	Message string `json:"message"`
	Data    [][]interface{}
}

type OneCandleStick struct {
	Time         time.Time
	OpeningPrice float64
	ClosingPrice float64
	HighPrice    float64
	LowPrice     float64
	UnitsTraded  float64
}

func newCandleStick(rawCandleStick RawCandleStick) []OneCandleStick {
	candleStick := make([]OneCandleStick, len(rawCandleStick.Data))
	for index, data := range rawCandleStick.Data {
		candleStick[index].Time = milliStringToTime(strconv.FormatInt(int64(int(data[0].(float64))), 10))
		candleStick[index].OpeningPrice, _ = strconv.ParseFloat(data[1].(string), 64)
		candleStick[index].ClosingPrice, _ = strconv.ParseFloat(data[2].(string), 64)
		candleStick[index].HighPrice, _ = strconv.ParseFloat(data[3].(string), 64)
		candleStick[index].LowPrice, _ = strconv.ParseFloat(data[4].(string), 64)
		candleStick[index].UnitsTraded, _ = strconv.ParseFloat(data[5].(string), 64)
	}
	return candleStick
}

/*====================== Private API 관련 ========================
1. Account
2. Balance

*/

//==============================ACCOUNT SETTING======================================

type Account struct {
	ID       string
	Created  time.Time
	Balance  float64
	TradeFee float64
}

func newAccount(rawAccount map[string]interface{}) Account {
	newAccount := Account{}
	newAccount.ID = rawAccount["account_id"].(string)
	newAccount.Created = milliStringToTime(rawAccount["created"].(string))
	newAccount.Balance, _ = strconv.ParseFloat(rawAccount["balance"].(string), 64)
	newAccount.TradeFee, _ = strconv.ParseFloat(rawAccount["trade_fee"].(string), 64)
	return newAccount
}

//==============================BALANCE SETTING======================================

type Balance struct {
	Total     float64
	InUse     float64
	Available float64
	XCoinLast float64
}

func newBalance(rawBalance map[string]interface{}, coin string) *Balance {
	newBalance := Balance{}
	newBalance.Total, _ = strconv.ParseFloat(rawBalance["total_"+coin].(string), 64)
	newBalance.InUse, _ = strconv.ParseFloat(rawBalance["in_use_"+coin].(string), 64)
	newBalance.Available, _ = strconv.ParseFloat(rawBalance["available_"+coin].(string), 64)
	newBalance.XCoinLast, _ = strconv.ParseFloat(rawBalance["xcoin_last_"+coin].(string), 64)
	return &newBalance
}

//==============================BALANCE SETTING======================================

type UserTicker struct {
	OpeningPrice    float64
	ClosingPrice    float64
	AveragePrice    float64
	MaxPrice        float64
	MinPrice        float64
	UnitsTraded     float64
	Volume1Day      float64
	Volume7Day      float64
	Fluctate24H     float64
	FluctateRate24H float64
}

func newUserTicker(rawUserTicker map[string]interface{}) UserTicker {
	newUserTicker := UserTicker{}
	newUserTicker.OpeningPrice, _ = strconv.ParseFloat(rawUserTicker["opening_price"].(string), 64)
	newUserTicker.ClosingPrice, _ = strconv.ParseFloat(rawUserTicker["closing_price"].(string), 64)
	newUserTicker.AveragePrice, _ = strconv.ParseFloat(rawUserTicker["average_price"].(string), 64)
	newUserTicker.MinPrice, _ = strconv.ParseFloat(rawUserTicker["min_price"].(string), 64)
	newUserTicker.MaxPrice, _ = strconv.ParseFloat(rawUserTicker["max_price"].(string), 64)
	newUserTicker.UnitsTraded, _ = strconv.ParseFloat(rawUserTicker["units_traded"].(string), 64)
	newUserTicker.Volume1Day, _ = strconv.ParseFloat(rawUserTicker["volume_1day"].(string), 64)
	newUserTicker.Volume7Day, _ = strconv.ParseFloat(rawUserTicker["volume_7day"].(string), 64)
	newUserTicker.Fluctate24H, _ = strconv.ParseFloat(rawUserTicker["fluctate_24H"].(string), 64)
	newUserTicker.FluctateRate24H, _ = strconv.ParseFloat(rawUserTicker["fluctate_rate_24H"].(string), 64)
	return newUserTicker
}

//==============================ORDER SETTING======================================

type Order struct {
	OrderDate       time.Time
	OrderCurrency   Currency
	PaymentCurrency Currency
	OrderID         string
	Price           float64
	Type            string
	Units           float64
	UnitsRemaining  float64
	WatchPrice      float64
}

func newOrder(rawOrder map[string]interface{}) Order {
	newOrder := Order{}
	newOrder.OrderDate = microStringToTime(rawOrder["order_date"].(string))
	newOrder.OrderCurrency = Currency(strings.ToLower(rawOrder["order_currency"].(string)))
	newOrder.PaymentCurrency = Currency(strings.ToLower(rawOrder["payment_currency"].(string)))
	newOrder.OrderID = rawOrder["order_id"].(string)
	newOrder.Price, _ = strconv.ParseFloat(rawOrder["price"].(string), 64)
	newOrder.Type = rawOrder["type"].(string)
	newOrder.Units, _ = strconv.ParseFloat(rawOrder["units"].(string), 64)
	newOrder.UnitsRemaining, _ = strconv.ParseFloat(rawOrder["units_remaining"].(string), 64)
	newOrder.WatchPrice, _ = strconv.ParseFloat(rawOrder["watch_price"].(string), 64)
	return newOrder
}

//==============================ORDER SETTING======================================

type SingleOrderDetail struct {
	TransactionDate time.Time
	Price           float64
	Units           float64
	FeeCurrency     Currency
	Fee             float64
	Total           float64
}

type OrderDetail struct {
	OrderDate       time.Time
	Type            string
	OrderStatus     string
	OrderCurrency   Currency
	PaymentCurrency Currency
	OrderPrice      float64
	OrderQty        float64
	CancelDate      time.Time
	CancelType      string
	Contract        []SingleOrderDetail
}

func newOrderDetail(newOrderDetail OrderDetail, rawOrderDetail map[string]interface{}) OrderDetail {
	newOrderDetail.OrderDate = microStringToTime(rawOrderDetail["order_date"].(string))
	newOrderDetail.Type = rawOrderDetail["type"].(string)
	newOrderDetail.OrderStatus = rawOrderDetail["order_status"].(string)
	newOrderDetail.OrderCurrency = Currency(strings.ToLower(rawOrderDetail["order_currency"].(string)))
	newOrderDetail.PaymentCurrency = Currency(strings.ToLower(rawOrderDetail["payment_currency"].(string)))
	newOrderDetail.OrderPrice, _ = strconv.ParseFloat(rawOrderDetail["order_price"].(string), 64)
	newOrderDetail.OrderQty, _ = strconv.ParseFloat(rawOrderDetail["payment_currency"].(string), 64)
	if rawOrderDetail["cancel_date"].(string) != "" {
		newOrderDetail.CancelDate = microStringToTime(rawOrderDetail["cancel_date"].(string))
	}
	newOrderDetail.CancelType = rawOrderDetail["cancel_type"].(string)

	contracts := rawOrderDetail["contract"].([]interface{})
	newOrderDetail.Contract = make([]SingleOrderDetail, len(contracts))
	for index, data := range contracts {
		singleContract := data.(map[string]interface{})
		newOrderDetail.Contract[index].TransactionDate = microStringToTime(singleContract["transaction_date"].(string))
		newOrderDetail.Contract[index].Price, _ = strconv.ParseFloat(singleContract["price"].(string), 64)
		newOrderDetail.Contract[index].Units, _ = strconv.ParseFloat(singleContract["units"].(string), 64)
		newOrderDetail.Contract[index].FeeCurrency = Currency(strings.ToLower(singleContract["fee_currency"].(string)))
		newOrderDetail.Contract[index].Fee, _ = strconv.ParseFloat(singleContract["fee"].(string), 64)
		newOrderDetail.Contract[index].Total, _ = strconv.ParseFloat(singleContract["total"].(string), 64)
	}

	return newOrderDetail
}

//==============================ORDER SETTING======================================

type Transactions struct {
	Search          SearchType
	TransferDate    time.Time
	OrderCurrency   Currency
	PaymentCurrency Currency
	Units           float64
	Price           float64
	Amount          float64
	FeeCurrency     Currency
	Fee             float64
	OrderBalance    float64
	PaymentBalance  float64
}

func newTransaction(rawTransaction map[string]interface{}) Transactions {
	newTransaction := Transactions{}
	newTransaction.Search = SearchType(rawTransaction["search"].(string))
	newTransaction.TransferDate = microStringToTime(rawTransaction["transfer_date"].(string))
	newTransaction.OrderCurrency = Currency(strings.ToLower(rawTransaction["order_currency"].(string)))
	newTransaction.PaymentCurrency = Currency(strings.ToLower(rawTransaction["payment_currency"].(string)))
	newTransaction.FeeCurrency = Currency(strings.ToLower(rawTransaction["fee_currency"].(string)))
	newTransaction.Units, _ = strconv.ParseFloat(rawTransaction["units"].(string), 64)
	newTransaction.Price, _ = strconv.ParseFloat(rawTransaction["price"].(string), 64)
	newTransaction.Amount, _ = strconv.ParseFloat(rawTransaction["amount"].(string), 64)
	newTransaction.Fee, _ = strconv.ParseFloat(rawTransaction["fee"].(string), 64)
	newTransaction.OrderBalance, _ = strconv.ParseFloat(rawTransaction["order_balance"].(string), 64)
	newTransaction.PaymentBalance, _ = strconv.ParseFloat(rawTransaction["payment_balance"].(string), 64)
	return newTransaction
}

//==============================PLACE SETTING======================================
