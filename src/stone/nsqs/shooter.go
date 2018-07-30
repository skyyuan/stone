package nsqs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"stone/common"
)

// ShootMessage Shoot message
func ShootMessage(address, topic string, payload interface{}) error {
	url := fmt.Sprintf("http://%s/pub?topic=%s", address, topic)
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return doBytesPost(url, data)
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
		return err
	}
	resp.Body.Close()
	return err
}
