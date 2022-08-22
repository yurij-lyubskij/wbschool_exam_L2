package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const ctype = `application/x-www-form-urlencoded`

func main() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdf")
	encodedData := data.Encode()
	fmt.Println(encodedData)
	url := "http://127.0.0.1:8080"
	create := "/create_event"
	method := "POST"
	req, err := http.NewRequest(method, url+create, strings.NewReader(encodedData))
	req.Header.Set("Content-Type", ctype)
	if err != nil {
		fmt.Println("error happend", err)
		return
	}
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error happend", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("resp %#v\n", string(respBody))
}
