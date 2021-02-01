package gobithumb

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type httpRequester struct {
	connectKey string
	secretKey  string

	basicUrl string

	publicClient  *http.Client
	privateClient *http.Client
}

func newHttpRequester(connectKey string, secretKey string) *httpRequester {

	httpRequester := httpRequester{}

	httpRequester.connectKey = connectKey
	httpRequester.secretKey = secretKey
	httpRequester.basicUrl = "https://api.bithumb.com"

	httpRequester.publicClient = &http.Client{}
	httpRequester.privateClient = &http.Client{}

	return &httpRequester
}

func (h *httpRequester) requestPublic(order publicOrder, data string) []byte {

	request, err := http.NewRequest("GET", h.basicUrl+string(order)+"/"+data, nil)
	if err != nil {
		panic("Failed to create Request")
	}

	response, err := h.publicClient.Do(request)
	if err != nil {
		panic("Failed to receive Data, check server status")
	}

	byteResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		timelog("Failed to receive Data")
	}
	return byteResponse
}

func (h *httpRequester) requestPrivate(passVal map[string]string) []byte {

	// request body 설정
	requestBody := url.Values{}
	for index, data := range passVal {
		requestBody.Set(index, data)
	}
	requestBodyString := requestBody.Encode()

	// nonce 및 api-sign 가져오기
	nonce := fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))
	apiSignVal := h.encryptData(passVal["endpoint"], requestBodyString, nonce)

	// request 객체 생성
	request, err := http.NewRequest("POST", h.basicUrl+passVal["endpoint"], bytes.NewBufferString(requestBodyString))
	if err != nil {
		panic("Failed to create Request")
	}

	request.Header.Add("Api-Key", h.connectKey)
	request.Header.Add("Api-Sign", apiSignVal)
	request.Header.Add("Api-Nonce", nonce)
	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(requestBodyString)))

	response, err := h.privateClient.Do(request)
	if err != nil {
		panic("Failed to receive Data, check server stauts")
	}

	byteResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to receive Data")
	}
	return byteResponse
}

func (h *httpRequester) encryptData(endpoint string, body string, nonce string) string {
	reqRawString := endpoint + string(0) + body + string(0) + nonce

	hmacParsed := hmac.New(sha512.New, []byte(h.secretKey))
	hmacParsed.Write([]byte(reqRawString))

	hexData := hex.EncodeToString(hmacParsed.Sum(nil))
	byteHexData := []byte(hexData)
	hmacParsed.Reset()

	result := base64.StdEncoding.EncodeToString(byteHexData)
	return result
}
