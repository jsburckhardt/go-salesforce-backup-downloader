package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func sendRequest(verb string, url string, headers map[string]string, data []byte) (int, []byte, error, []string) {

	fmt.Printf("verb: %s\n", verb)
	fmt.Printf("url: %s\n", url)
	fmt.Printf("headers: %v\n", headers)
	fmt.Printf("data: %s\n", string(data))

	req, err := http.NewRequest(verb, url, bytes.NewBuffer(data))
	if err != nil {
		return http.StatusInternalServerError, nil, err, nil
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err, nil
	}

	head := resp.Header["Content-Length"]

	return http.StatusOK, body, nil, head
}
