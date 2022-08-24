package main

import (
	"errors"
	"net/http"
	"time"
)

const inputError = 400
const businessError = 503
const serverError = 500

//GetEvents - общий тип функций для получения
//массива событий
type GetEvents func(id string, date time.Time) []Event

//выводим ошибку и возвращаем статус http
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

//выводим ответ, как JSON
func resWrite(w http.ResponseWriter, resResp ResultResponse) {
	var bytes []byte
	bytes, err := resResp.MarshalJSON()
	if err != nil {
		w.WriteHeader(serverError)
		return
	}
	w.Write(bytes)
}

//хэндлер, возвращающий 400, если эндпоинта не существует
//(обычно - 404, но по условию задачи - 400)
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	err := errors.New("page not found")
	errWrite(w, err, inputError)
	return
}

//хэндлер создания события
func createHandler(w http.ResponseWriter, r *http.Request) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/json")
	//парсим и валидируем
	event, err := parsePost(r, false, true)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	//создаем событие
	num, err := storage.createEvent(event)
	if err != nil {
		errWrite(w, err, businessError)
		return
	}
	//создаем ответ
	resResp.Result = CRUDResult{Num: num, Success: true}
	//отправляем ответ
	resWrite(w, resResp)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/json")
	//парсим и валидируем
	event, err := parsePost(r, true, true)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	//обновляем событие
	err = storage.updateEvent(event)
	if err != nil {
		errWrite(w, err, businessError)
		return
	}
	//создаем ответ
	resResp.Result = CRUDResult{Num: event.Num, Success: true}
	//отправляем ответ
	resWrite(w, resResp)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/json")
	//парсим и валидируем
	event, err := parsePost(r, true, false)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	//удаляем событие
	err = storage.deleteEvent(event)
	if err != nil {
		errWrite(w, err, businessError)
		return
	}
	//создаем ответ
	resResp.Result = CRUDResult{Num: event.Num, Success: true}
	//отправляем ответ
	resWrite(w, resResp)
}

func eventsHandler(w http.ResponseWriter, r *http.Request, get GetEvents) {
	var resResp ResultResponse
	w.Header().Set("Content-Type", "application/json")
	//парсим и валидируем
	event, err := parseGet(r)
	if err != nil {
		errWrite(w, err, inputError)
		return
	}
	//получаем список событий
	events := get(event.UserID, event.Date)
	//создаем ответ
	resResp.Result = GetResult{events}
	//отправляем ответ
	resWrite(w, resResp)
}

func dayEventsHandler(w http.ResponseWriter, r *http.Request) {
	//передаем метод получения списка событий текущего дня
	eventsHandler(w, r, storage.getEventsByDay)
}

func weekEventsHandler(w http.ResponseWriter, r *http.Request) {
	//передаем метод получения списка событий текущей недели
	eventsHandler(w, r, storage.getEventsByWeek)
}

func monthEventsHandler(w http.ResponseWriter, r *http.Request) {
	//передаем метод получения списка событий текущего месяца
	eventsHandler(w, r, storage.getEventsByMonth)
}
