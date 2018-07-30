package etherscan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"stone/common"
)

// EtherscanAPIKey etherscan apikey
const EtherscanAPIKey = "AUI76EBSSTAZ3TYU12IRGDX8YY7U9U5RXR"

// EtherScanAPIURL Etherscan api url
const EtherScanAPIURL = "https://api.etherscan.io/api"

// RopstenAPIURL ropsten api url
const RopstenAPIURL = "https://ropsten.etherscan.io/api"

// SendRawTransaction SendRawTransaction
func SendRawTransaction(isTestnet bool, hexstr string) error {
	apiURL := EtherScanAPIURL
	if isTestnet {
		apiURL = RopstenAPIURL
	}
	url := fmt.Sprintf("%s?module=proxy&action=eth_sendRawTransaction&hex=%s&apikey=%s", apiURL, hexstr, EtherscanAPIKey)
	return doBytesPost(url, nil)
}

func doBytesPost(url string, data []byte) error {
	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		common.Logger.Info("http.NewRequest,[err=%s][url=%s]", err, url)
		return err
	}
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		common.Logger.Error(err)
		return err
	}
	defer resp.Body.Close()
	resIO, _ := ioutil.ReadAll(resp.Body)
	res := map[string]interface{}{}
	json.Unmarshal(resIO, &res)
	common.Logger.Info("etherscan send at :", res)
	return err
}
