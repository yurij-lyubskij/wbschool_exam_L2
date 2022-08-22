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

func postReq(endpoint string, data url.Values) {
	encodedData := data.Encode()
	fmt.Println(encodedData)
	url := "http://127.0.0.1:8080"
	method := "POST"
	req, err := http.NewRequest(method, url+endpoint, strings.NewReader(encodedData))
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

func getReq(endpoint string, data url.Values) {
	encodedData := data.Encode()
	fmt.Println(encodedData)
	url := "http://127.0.0.1:8080"
	method := "GET"
	req, err := http.NewRequest(method, url+endpoint+encodedData, nil)
	req.Header.Set("Content-Type", ctype)
	if err != nil {
		fmt.Println("error happend", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error happend", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("resp %#v\n", string(respBody))
}

func testCreate() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdf")
	create := "/create_event"
	postReq(create, data)
}

func testUpdate() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdffadsfdsaf")
	data.Set("event_num", "1")
	update := "/update_event"
	postReq(update, data)
}

func testDelete() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("event_num", "1")
	rdelete := "/delete_event"
	postReq(rdelete, data)
}

func testGet(endpoint string) {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdffadsfdsaf")
	data.Set("event_num", "1")
	getReq(endpoint, data)
}

func testGetDay() {
	endpoint := "/events_for_day?"
	testGet(endpoint)
}

func testGetWeek() {
	endpoint := "/events_for_week?"
	testGet(endpoint)
}

func testGetMonth() {
	endpoint := "/events_for_month?"
	testGet(endpoint)
}

func main() {
	testGetDay()
	testCreate()
	testGetDay()
	testCreate()
	testUpdate()
	testGetDay()
	testDelete()
	testGetDay()
	testGetWeek()
	testGetMonth()
}
