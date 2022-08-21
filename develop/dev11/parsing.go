package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func parseParams(getNum bool, values url.Values) (Event, error) {
	var event Event
	err := errors.New("parameter not found")
	event.UserID = values.Get("user_id")
	if event.UserID == "" {
		return event, err
	}
	dateStr := values.Get("date")
	if dateStr == "" {
		return event, err
	}
	var errTime error
	event.Date, errTime = time.Parse(dateForm, dateStr)
	if errTime != nil {
		return event, err
	}
	if !getNum {
		return event, nil
	}
	numStr := values.Get("event_num")
	if numStr == "" {
		return event, err
	}
	event.Num, err = strconv.Atoi(numStr)
	return event, err
}

func parseGet(r *http.Request, getNum bool) (Event, error) {
	query := r.URL.Query()
	return parseParams(getNum, query)
}

func parsePost(r *http.Request, getNum bool) (Event, error) {
	var event Event
	err := r.ParseForm()
	if err != nil {
		return event, err
	}
	values := r.PostForm
	return parseParams(getNum, values)
}
