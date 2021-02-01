package gobithumb

import (
	"fmt"
	"strconv"
	"time"
)

func timelog(a ...interface{}) {
	fmt.Print(time.Now().Format(time.StampMilli) + "\t")
	fmt.Println(a...)
}

func milliStringToTime(milliString string) time.Time {

	timeSec, _ := strconv.ParseInt(milliString[:len(milliString)-3], 10, 64)
	timeMilli, _ := strconv.ParseInt(milliString[len(milliString)-3:], 10, 64)
	return time.Unix(timeSec, timeMilli*1000000)
}

func microStringToTime(microString string) time.Time {

	timeSec, _ := strconv.ParseInt(microString[:len(microString)-6], 10, 64)
	timeMilli, _ := strconv.ParseInt(microString[len(microString)-6:], 10, 64)
	return time.Unix(timeSec, timeMilli*1000)
}

func rawBalanceStringToBalance(raw string) (string, int) {

	if raw[0] == 't' {
		return raw[6:], 1 //total_coin
	} else if raw[0] == 'i' {
		return raw[7:], 2 //in_use_coin
	} else if raw[0] == 'a' {
		return raw[10:], 3 //available_coin
	} else {
		return raw[11:], 4 //xcoin
	}
}
