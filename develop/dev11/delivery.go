package main

import "net/http"

const inputError = 400
const serverError = 500

func errWrite(w http.ResponseWriter, err error) {
	var bytes []byte
	var errResp ErrorResponse
	errResp.ErrorMsg = err.Error()
	bytes, err = errResp.MarshalJSON()
	if err != nil {
		w.WriteHeader(serverError)
		return
	}
	w.WriteHeader(inputError)
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
		errWrite(w, err)
		return
	}
	num, err := storage.createEvent(event)
	if err != nil {
		errWrite(w, err)
		return
	}
	resResp.Result = CRUDResult{Num: num, Success: true}
	resWrite(w, resResp)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

}

func dayEventsHandler(w http.ResponseWriter, r *http.Request) {

}

func weekEventsHandler(w http.ResponseWriter, r *http.Request) {

}

func monthEventsHandler(w http.ResponseWriter, r *http.Request) {

}
