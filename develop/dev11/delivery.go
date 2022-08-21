package main

import "net/http"

const inputError = 400
const serverError = 500

func createHandler(w http.ResponseWriter, r *http.Request) {
	event, err := parsePost(r, false, true)
	var errResp ErrorResponse
	var resResp ResultResponse
	var bytes []byte
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		errResp.ErrorMsg = err.Error()
		bytes, err = errResp.MarshalJSON()
		if err != nil {
			w.WriteHeader(serverError)
			return
		}
		w.WriteHeader(inputError)
		w.Write(bytes)
		return
	}
	resResp.Result = CRUDResult{Num: event.Num, Success: true}
	bytes, err = resResp.MarshalJSON()
	if err != nil {
		w.WriteHeader(serverError)
		return
	}
	w.Write(bytes)
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
