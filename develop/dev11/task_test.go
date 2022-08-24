package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestCase struct {
	endpoint   string
	method     string
	query      string
	Response   string
	StatusCode int
}

const localhost = "http://127.0.0.1:8080"
const ctype = `application/x-www-form-urlencoded`

func TestCreate(t *testing.T) {
	cases := []TestCase{
		TestCase{
			endpoint:   "/create_event",
			method:     "POST",
			query:      "date=2019-09-09&description=asdffadsfdsaf&user_id=John",
			Response:   `{"result":"{"event_num":1,"success":true}"}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := localhost + item.endpoint
		req := httptest.NewRequest(item.method, url, strings.NewReader(item.query))
		req.Header.Set("Content-Type", ctype)
		w := httptest.NewRecorder()
		createHandler(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []TestCase{
		TestCase{
			endpoint:   "/update_event",
			method:     "POST",
			query:      "date=2019-09-09&description=asdffadsfdsaf&event_num=1&user_id=John",
			Response:   `{"result":"{"event_num":1,"success":true}"}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := localhost + item.endpoint
		req := httptest.NewRequest(item.method, url, strings.NewReader(item.query))
		req.Header.Set("Content-Type", ctype)
		w := httptest.NewRecorder()
		updateHandler(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []TestCase{
		TestCase{
			endpoint:   "/delete_event",
			method:     "POST",
			query:      "date=2019-09-09&description=asdffadsfdsaf&event_num=1&user_id=John",
			Response:   `{"result":"{"event_num":1,"success":true}"}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := localhost + item.endpoint
		req := httptest.NewRequest(item.method, url, strings.NewReader(item.query))
		req.Header.Set("Content-Type", ctype)
		w := httptest.NewRecorder()
		deleteHandler(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}

func TestGetbyDay(t *testing.T) {
	cases := []TestCase{
		TestCase{
			endpoint:   "/events_for_day?",
			method:     "GET",
			query:      "date=2019-09-09&description=asdffadsfdsaf&event_num=1&user_id=John",
			Response:   `{"result":"{"events":[]}"}`,
			StatusCode: http.StatusOK,
		},
	}
	for caseNum, item := range cases {
		url := localhost + item.endpoint + item.query
		req := httptest.NewRequest(item.method, url, nil)
		req.Header.Set("Content-Type", ctype)
		w := httptest.NewRecorder()
		dayEventsHandler(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}
