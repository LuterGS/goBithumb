package gobithumb

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func Timelog(a ...interface{}) {
	fmt.Print(time.Now().Format(time.StampMilli) + "\t")
	fmt.Println(a...)
}

func ReadFromFile(fileName string) map[string]string {

	curLoc, _ := os.Getwd()
	Timelog(curLoc + fileName)
	dbData, err := ioutil.ReadFile(curLoc + fileName)
	if err != nil {
		Timelog("파일을 불러오는 데 실패했습니다. 프로그램을 종료합니다.")
		panic("파일 Read 실패")
	}

	var fileMap map[string]string
	fileMap = make(map[string]string)
	lines := strings.Split(string(dbData), "\n")
	for i := 0; i < len(lines); i++ {
		singleLine := strings.Split(lines[i], "=")
		fileMap[singleLine[0]] = singleLine[1]
	}

	return fileMap
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
