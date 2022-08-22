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

//запрос POST с передаваемым эндпоинтом и urlencoded телом
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

//запрос POST с передаваемым эндпоинтом и queryparams
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

//тестовый запрос на создание события
func testCreate() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdf")
	create := "/create_event"
	postReq(create, data)
}

//тестовый запрос на изменение события
func testUpdate() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdffadsfdsaf")
	data.Set("event_num", "1")
	update := "/update_event"
	postReq(update, data)
}

//тестовый запрос на удаление события
func testDelete() {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("event_num", "1")
	rdelete := "/delete_event"
	postReq(rdelete, data)
}

//запрос на получение событий
func testGet(endpoint string) {
	data := url.Values{}
	data.Set("user_id", "John")
	data.Set("date", "2019-09-09")
	data.Set("description", "asdffadsfdsaf")
	data.Set("event_num", "1")
	getReq(endpoint, data)
}

//тестовый запрос на получение событий за день
func testGetDay() {
	endpoint := "/events_for_day?"
	testGet(endpoint)
}

//тестовый запрос на получение событий за неделю
func testGetWeek() {
	endpoint := "/events_for_week?"
	testGet(endpoint)
}

//тестовый запрос на получение событий за месяц
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
