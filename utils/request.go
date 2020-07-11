package utils

import (
	"io/ioutil"
	"net/http"
)

func GetCurl(hostUrl string, token string) (response *http.Response, body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", token)
	response, err = client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	return
}
