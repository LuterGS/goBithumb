package gobithumb

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type publicOrder string
type privateOrder string

type BithumbRequester struct {
	requester *httpRequester

	basicUrl string

	ticker       publicOrder
	orderbook    publicOrder
	trHistory    publicOrder
	assetsStatus publicOrder
	btci         publicOrder
	candlestick  publicOrder

	balance       privateOrder
	account       privateOrder
	walletAddress privateOrder
	tickerUser    privateOrder
	orders        privateOrder
	orderDetail   privateOrder
	transactions  privateOrder

	place          privateOrder
	cancel         privateOrder
	marketBuy      privateOrder
	marketSell     privateOrder
	stopLimit      privateOrder
	withdrawalCoin privateOrder
	withdrawalKRW  privateOrder
}

func NewBithumb(connectKey string, secretKey string) *BithumbRequester {

	bithumbRequester := BithumbRequester{}

	bithumbRequester.requester = newHttpRequester(connectKey, secretKey)

	// init public API address
	bithumbRequester.ticker = "/public/ticker"
	bithumbRequester.orderbook = "/public/orderbook"
	bithumbRequester.trHistory = "/public/transaction_history"
	bithumbRequester.assetsStatus = "/public/assetsstatus"
	bithumbRequester.btci = "/public/btci"
	bithumbRequester.candlestick = "/public/candlestick"

	// init private API address
	bithumbRequester.balance = "/info/balance"
	bithumbRequester.account = "/info/account"
	bithumbRequester.walletAddress = "/info/wallet_address"
	bithumbRequester.tickerUser = "/info/ticker"
	bithumbRequester.orders = "/info/orders"
	bithumbRequester.orderDetail = "/info/order_detail"
	bithumbRequester.transactions = "/info/user_transactions"

	bithumbRequester.place = "/trade/place"
	bithumbRequester.cancel = "/trade/cancel"
	bithumbRequester.marketBuy = "/trade/market_buy"
	bithumbRequester.marketSell = "/trade/market_sell"
	bithumbRequester.stopLimit = "/trade/stop_limit"
	bithumbRequester.withdrawalCoin = "/trade/btc_withdrawal"
	bithumbRequester.withdrawalKRW = "/trade/krw_withdrawal"

	return &bithumbRequester
}

func (b *BithumbRequester) publicRequest(reqUrl publicOrder, reqBody string) map[string]interface{} {
	requestResult := b.requester.requestPublic(reqUrl, reqBody)
	var result map[string]interface{}
	_ = json.Unmarshal(requestResult, &result)
	return result
}

func (b *BithumbRequester) GetTradableCoinList() []Currency {
	requestResult := b.requester.requestPublic(b.ticker, "all_krw")
	var tempResult map[string]interface{}
	var stringResult []string
	var result []Currency
	_ = json.Unmarshal(requestResult, &tempResult)
	datas := tempResult["data"].(map[string]interface{})
	for index := range datas {
		stringResult = append(stringResult, index)
	}
	sort.Strings(stringResult)
	for _, data := range stringResult {
		result = append(result, Currency(strings.ToLower(data)))
	}
	return result[:len(result)-1] // 마지막 하나가 date이므로, date를 제외하고 전달함
}

func (b *BithumbRequester) GetTicker(orderCurrency Currency, paymentCurrency Currency) (map[Currency]Ticker, time.Time, error) {
	reqResult := b.publicRequest(b.ticker, string(orderCurrency)+"_"+string(paymentCurrency))
	result := make(map[Currency]Ticker)

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetTicker failed : ", reqResult["message"].(string))
		return result, time.Now(), errors.New(errNo)
	}

	// Convert data
	datas := reqResult["data"].(map[string]interface{})

	// Time check
	reqTime := milliStringToTime(datas["date"].(string))
	delete(datas, "date")

	if orderCurrency != ALL {
		result[orderCurrency] = newTicker(datas)
	} else {
		for index, data := range datas {
			result[Currency(index)] = newTicker(data.(map[string]interface{}))
		}
	}
	return result, reqTime, nil
}

func (b *BithumbRequester) GetOrderbook(orderCurrency Currency, paymentCurrency Currency) (map[Currency]Orderbook, time.Time, error) {
	reqResult := b.publicRequest(b.orderbook, string(orderCurrency)+"_"+string(paymentCurrency))
	result := make(map[Currency]Orderbook)

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetOrderbook failed : ", reqResult["message"].(string))
		return result, time.Now(), errors.New(errNo)
	}

	// Convert data
	datas := reqResult["data"].(map[string]interface{})

	// Time check
	reqTime := milliStringToTime(datas["timestamp"].(string))
	delete(datas, "timestamp")

	if orderCurrency != ALL {
		result[orderCurrency] = newOrderbook(datas)
	} else {
		delete(datas, "payment_currency")
		for index, data := range datas {
			result[Currency(index)] = newOrderbook(data.(map[string]interface{}))
		}
	}
	return result, reqTime, nil
}

func (b *BithumbRequester) GetTransactionHistory(orderCurrency Currency, paymentCurrency Currency) ([]OneTransaction, error) {
	reqResult := b.publicRequest(b.trHistory, string(orderCurrency)+"_"+string(paymentCurrency))

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetTransactionHistory failed : ", reqResult["message"].(string))
		return nil, errors.New(errNo)
	}

	// Convert data and return
	return newTransactionHistory(reqResult["data"].([]interface{})), nil
}

func (b *BithumbRequester) GetAssetsStatus(orderCurrency Currency) (bool, bool, error) {
	reqResult := b.publicRequest(b.assetsStatus, string(orderCurrency))

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetAssetsStatus failed : ", reqResult["message"].(string))
		return false, false, errors.New(errNo)
	}

	// Convert data and return
	datas := reqResult["data"].(map[string]interface{})
	depositStatus, withdrawlStatus := false, false
	if val, _ := fmt.Print(datas["deposit_status"]); val == 1 {
		depositStatus = true
	}
	if val2, _ := fmt.Print(datas["withdrawal_status"]); val2 == 1 {
		withdrawlStatus = true
	}
	return depositStatus, withdrawlStatus, nil
}

func (b *BithumbRequester) GetBTCI() (BTCI, time.Time, error) {
	reqResult := b.publicRequest(b.btci, "")

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetBTCI failed : ", reqResult["message"].(string))
		var noneBTCI BTCI
		return noneBTCI, time.Now(), errors.New(errNo)
	}

	// Convert data
	datas := reqResult["data"].(map[string]interface{})

	// Time check
	reqTime := milliStringToTime(datas["date"].(string))
	delete(datas, "date")

	// return result
	return newBTCI(datas["btai"], datas["btmi"]), reqTime, nil
}

func (b *BithumbRequester) GetCandleStick(orderCurreny Currency, paymentCurrency Currency, chartInterval TimeInterval) ([]OneCandleStick, error) {
	body := string(orderCurreny) + "_" + string(paymentCurrency) + "/" + string(chartInterval)
	requestResult := b.requester.requestPublic(b.candlestick, body)
	var rawResult RawCandleStick
	_ = json.Unmarshal(requestResult, &rawResult)

	timelog(string(requestResult))

	if rawResult.Status != 0 {
		timelog("GetCandleStick failed : ", rawResult.Message)
		var result []OneCandleStick
		return result, errors.New(strconv.Itoa(rawResult.Status))

	}
	return newCandleStick(rawResult), nil
}

func (b *BithumbRequester) privateRequest(reqUrl privateOrder, requestVal map[string]string) map[string]interface{} {
	requestVal["endpoint"] = string(reqUrl)
	reqResult := b.requester.requestPrivate(requestVal)
	var result map[string]interface{}
	_ = json.Unmarshal(reqResult, &result)
	return result
}

func (b *BithumbRequester) GetAccount(orderCurrency Currency, paymentCurrency Currency) (Account, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	reqResult := b.privateRequest(b.account, passVal)
	var result Account

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetAccount failed : ", reqResult["message"].(string))
		return result, errors.New(errNo)
	}

	return newAccount(reqResult["data"].(map[string]interface{})), nil
}

func (b *BithumbRequester) GetBalance(orderCurrency Currency) (map[Currency]*Balance, error) {
	passVal := make(map[string]string)
	passVal["currency"] = string(orderCurrency)
	reqResult := b.privateRequest(b.balance, passVal)

	// Error check
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetBalance failed : ", reqResult["message"].(string))
		return nil, errors.New(errNo)
	}

	// Convert data
	result := make(map[Currency]*Balance)
	datas := reqResult["data"].(map[string]interface{})

	if orderCurrency != ALL {
		result[orderCurrency] = newBalance(datas, string(orderCurrency))
		result[KRW] = newBalance(datas, string(KRW))
	} else {
		for _, data := range COIN_ALL() {
			result[data] = &Balance{}
		}
		for index, data := range datas {
			coin, value := rawBalanceStringToBalance(index)
			if value == 1 {
				result[Currency(coin)].Total, _ = strconv.ParseFloat(data.(string), 64)
			}
			if value == 2 {
				result[Currency(coin)].InUse, _ = strconv.ParseFloat(data.(string), 64)
			}
			if value == 3 {
				result[Currency(coin)].Available, _ = strconv.ParseFloat(data.(string), 64)
			}
			if value == 4 {
				result[Currency(coin)].XCoinLast, _ = strconv.ParseFloat(data.(string), 64)
			}
		}
	}
	return result, nil
}

// TODO : Docs 쓸 때, 만약 주소가 없으면 정상 처리는 되나 아무 값도 리턴하지 않는다고 서술해야함.
func (b *BithumbRequester) GetWalletAddress(orderCurrency Currency) (string, error) {
	passVal := make(map[string]string)
	passVal["currency"] = string(orderCurrency)
	reqResult := b.privateRequest(b.walletAddress, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetAccount failed : ", reqResult["message"].(string))
		return "", errors.New(errNo)
	}

	return reqResult["data"].(map[string]interface{})["wallet_address"].(string), nil
}

func (b *BithumbRequester) GetUserTicker(orderCurrency Currency, paymentCurrency Currency) (UserTicker, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	reqResult := b.privateRequest(b.tickerUser, passVal)
	var result UserTicker

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetAccount failed : ", reqResult["message"].(string))
		return result, errors.New(errNo)
	}

	result = newUserTicker(reqResult["data"].(map[string]interface{}))
	return result, nil
}

// -> date에 값이 들어올 경우, 최측 하나만 사용
func (b *BithumbRequester) GetOrder(orderCurrency Currency, paymentCurrency Currency, count int, date ...time.Time) ([]Order, error) {

	// parameter 정상 체크
	if !(count > 0 && count < 1001) {
		return nil, errors.New("주문의 개수는 1~1000 사이의 정수여야 합니다.")
	}

	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["count"] = strconv.Itoa(count)
	if len(date) > 0 {
		passVal["after"] = strconv.FormatInt(date[0].Unix(), 10)
	}
	reqResult := b.privateRequest(b.orders, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetOrder failed : ", reqResult["message"].(string))
		return nil, errors.New(errNo)
	}

	// Data parse and return val create
	datas := reqResult["data"].([]interface{})
	var result []Order

	for _, data := range datas {
		result = append(result, newOrder(data.(map[string]interface{})))
	}

	return result, nil
}

func (b *BithumbRequester) GetOrderDetail(orderCurrency Currency, paymentCurrency Currency, orderId string) (OrderDetail, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["order_id"] = orderId
	reqResult := b.privateRequest(b.orderDetail, passVal)
	var result OrderDetail

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetOrderDetail failed : ", reqResult["message"].(string))
		return result, errors.New(errNo)
	}

	result = newOrderDetail(result, reqResult["data"].(map[string]interface{}))

	return result, errors.New(errNo)
}

func (b *BithumbRequester) GetTransactions(orderCurrency Currency, paymentCurrency Currency, search SearchType, offset_count ...int) ([]Transactions, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["searchGb"] = string(search)
	// 파라미터 정의 및 값 넣기
	if len(offset_count) == 2 {
		passVal["offset"] = strconv.Itoa(offset_count[0])
		passVal["count"] = strconv.Itoa(offset_count[1])
	}
	reqResult := b.privateRequest(b.transactions, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("GetTransactions failed : ", reqResult["message"].(string))
		return nil, errors.New(errNo)
	}

	// Convert data and input
	var result []Transactions
	datas := reqResult["data"].([]interface{})
	for _, data := range datas {
		result = append(result, newTransaction(data.(map[string]interface{})))
	}
	return result, nil
}

func (b *BithumbRequester) PlaceOrder(orderCurrency Currency, paymentCurrency Currency, amount float64, price float64, order string) (string, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["units"] = strconv.FormatFloat(amount, 'f', -1, 64)
	passVal["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	passVal["type"] = order
	reqResult := b.privateRequest(b.place, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("PlaceOrder failed : ", reqResult["message"].(string))
		return "", errors.New(errNo)
	}

	return reqResult["order_id"].(string), nil
}

func (b *BithumbRequester) CancelOrder(orderCurrency Currency, paymentCurrency Currency, orderId string, order string) error {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["order_id"] = orderId
	passVal["type"] = order
	reqResult := b.privateRequest(b.cancel, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("CancelOrder failed : ", reqResult["message"].(string))
		return nil
	}
	return errors.New(errNo)
}

func (b *BithumbRequester) MarketBuy(orderCurrency Currency, paymentCurrency Currency, amount float64) (string, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["units"] = strconv.FormatFloat(amount, 'f', -1, 64)
	reqResult := b.privateRequest(b.marketBuy, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("MarketBuy failed : ", reqResult["message"].(string))
		return "", errors.New(errNo)
	}

	return reqResult["order_id"].(string), nil
}

func (b *BithumbRequester) MarketSell(orderCurrency Currency, paymentCurrency Currency, amount float64) (string, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["units"] = strconv.FormatFloat(amount, 'f', -1, 64)
	reqResult := b.privateRequest(b.marketSell, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("MarketSell failed : ", reqResult["message"].(string))
		return "", errors.New(errNo)
	}

	return reqResult["order_id"].(string), nil
}

func (b *BithumbRequester) StopLimit(orderCurrency Currency, paymentCurrency Currency, watchPrice float64, price float64, amount float64, order string) (string, error) {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["payment_currency"] = string(paymentCurrency)
	passVal["watch_price"] = strconv.FormatFloat(watchPrice, 'f', -1, 64)
	passVal["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	passVal["units"] = strconv.FormatFloat(amount, 'f', -1, 64)
	passVal["type"] = order
	reqResult := b.privateRequest(b.stopLimit, passVal)

	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("StopLimit failed : ", reqResult["message"].(string))
		return "", errors.New(errNo)
	}

	return reqResult["order_id"].(string), nil
}

func (b *BithumbRequester) WithDrawCoin(orderCurrency Currency, amount float64, address string, destination ...interface{}) error {
	passVal := make(map[string]string)
	passVal["order_currency"] = string(orderCurrency)
	passVal["units"] = strconv.FormatFloat(amount, 'f', -1, 64)
	passVal["address"] = address

	//destination tag 설정
	if orderCurrency == XRP || orderCurrency == STEEM {
		if len(destination) == 1 && reflect.TypeOf(destination) == reflect.TypeOf(1) {
			passVal["destination"] = strconv.Itoa(destination[0].(int))
		} else {
			return errors.New("XRP 출금 시 destination tag(int) 를 지정해주지 않음, 또는 STEEM 출금 시 입금 메모를 지정해주지 않음")
		}
	}
	if orderCurrency == XMR {
		if len(destination) == 1 && reflect.TypeOf(destination) == reflect.TypeOf("1") {
			passVal["destination"] = destination[0].(string)
		} else {
			return errors.New("XMR 출금 시 Payment ID를 지정해주지 않음")
		}
	}

	reqResult := b.privateRequest(b.withdrawalCoin, passVal)
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("WithdrawCoin failed : ", reqResult["message"].(string))
		return errors.New(errNo)
	}
	return nil
}

func (b *BithumbRequester) WithdrawKRW(account string, price int) error {
	passVal := make(map[string]string)
	passVal["bank"] = "011_농협은행"
	passVal["account"] = account
	passVal["price"] = strconv.Itoa(price)
	reqResult := b.privateRequest(b.withdrawalKRW, passVal)
	errNo := reqResult["status"].(string)
	if errNo != "0000" {
		timelog("WithdrawKRW failed : ", reqResult["message"].(string))
		return errors.New(errNo)
	}
	return nil
}
