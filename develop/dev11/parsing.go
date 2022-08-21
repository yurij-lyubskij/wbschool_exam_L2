package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

func parseGet(r *http.Request, getNum bool) (Event, error) {
	query := r.URL.Query()
	err := errors.New("parameter not found")
	var event Event
	event.UserID = query.Get("user_id")
	if event.UserID == "" {
		return event, err
	}
	dateStr := query.Get("date")
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
	numStr := query.Get("event_num")
	if numStr == "" {
		return event, err
	}
	event.Num, err = strconv.Atoi(numStr)
	return event, err
}
