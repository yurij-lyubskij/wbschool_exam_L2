package main

import (
	"net/http"
	"time"
)

const inputError = 400
const businessError = 503
const serverError = 500

type GetEvents func(id string, date time.Time) []Event

func errWrite(w http.ResponseWriter, err error, code int) {
	var bytes []byte
	var errResp ErrorResponse
	errResp.ErrorMsg = err.Error()
	bytes, err = errResp.MarshalJSON()
	if err != nil {
		w.WriteHeader(serverError)
		return
	}
	w.WriteHeader(code)
	w.Write(bytes)
}

func resWrite(w http.ResponseWriter, resResp ResultResponse) {
	var bytes []byte
	bytes, err := resResp.MarshalJSON()
	if err != nil {
		w.WriteHeader(serverError)
		return
	}
	w.Write(bytes)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parsePost(r, false, true)
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	num, err := storage.createEvent(event)
	if err != nil {
		errWrite(w, err, businessError)
		return
	}
	resResp.Result = CRUDResult{Num: num, Success: true}
	resWrite(w, resResp)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	event, err := parsePost(r, true, true)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	err = storage.updateEvent(event)
	if err != nil {
		errWrite(w, err, businessError)
		return
	}
	resResp.Result = CRUDResult{Num: event.Num, Success: true}
	resWrite(w, resResp)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	event, err := parsePost(r, true, false)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	err = storage.deleteEvent(event)
	if err != nil {
		errWrite(w, err, businessError)
		return
	}
	resResp.Result = CRUDResult{Num: event.Num, Success: true}
	resWrite(w, resResp)
}

func eventsHandler(w http.ResponseWriter, r *http.Request, get GetEvents) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	event, err := parseGet(r)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	events := get(event.UserID, event.Date)
	resResp.Result = GetResult{events}
	resWrite(w, resResp)
}

func dayEventsHandler(w http.ResponseWriter, r *http.Request) {
	eventsHandler(w, r, storage.getEventsByDay)
}

func weekEventsHandler(w http.ResponseWriter, r *http.Request) {
	eventsHandler(w, r, storage.getEventsByWeek)
}

func monthEventsHandler(w http.ResponseWriter, r *http.Request) {
	eventsHandler(w, r, storage.getEventsByMonth)
}
