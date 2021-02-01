goBithumb
=======================

Go용 Bithumb API 클라이언트 프로그램입니다.
----

# Requirements
* Go (version 1.15.6 이상)
  


# How to Use
* 라이브러리 Import
```go
import(
	b "github.com/lutergs/gobithumb"
)
```
* 라이브러리 Init
```go
    BithumbClient := b.NewBithumb("YOUR CONNECT KEY", "YOUR SECRET KEY")
```

* 라이브러리 사용 예제
```go
    // Public API 사용 예시
    ticker, reqTime, err := BithumbClient.GetTicker(b.BTC, b.KRW)
    if err != nil{
        panic(err)	
    }
    
    // Private API 사용 시
    buyId, err := BithumbClient.MarketBuy(b.BTC, b.KRW, 0.0002)예
    if err != nil{
        panic(err)	
    }
    orderStatus, err := BithumbClient.GetOrderDetail(b.BTC, b.KRW, buyId)
    if err != nil{
        panic(err)	
    }
    fmt.Println("buy process : ", orderStatus)
```


# Docs
[여기](https://github.com/LuterGS/goBithumb/wiki) 를 참고
