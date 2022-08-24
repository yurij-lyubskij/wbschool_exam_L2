package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func parseParams(getNum bool, getDescr bool, values url.Values) (Event, error) {
	var event Event
	var err error
	errString := " parameter not found"

	param := "user_id"
	if !values.Has(param) {
		return event, errors.New(param + errString)
	}
	event.UserID = values.Get(param)

	param = "date"
	if !values.Has(param) {
		return event, errors.New(param + errString)
	}
	dateStr := values.Get(param)
	event.Date, err = time.Parse(dateForm, dateStr)
	if err != nil {
		return event, err
	}

	if !getNum {
		return event, nil
	}

	param = "event_num"
	if !values.Has(param) {
		return event, errors.New(param + errString)
	}
	numStr := values.Get(param)
	event.Num, err = strconv.Atoi(numStr)

	if !getDescr {
		return event, nil
	}

	param = "description"
	if !values.Has(param) {
		return event, errors.New(param + errString)
	}
	event.Description = values.Get(param)
	return event, nil
}

func parseGet(r *http.Request) (Event, error) {
	query := r.URL.Query()
	return parseParams(false, false, query)
}

func parsePost(r *http.Request, getNum bool, getDescr bool) (Event, error) {
	var event Event
	err := r.ParseForm()
	if err != nil {
		return event, err
	}
	values := r.PostForm
	return parseParams(getNum, getDescr, values)
}
