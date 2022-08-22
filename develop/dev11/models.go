package main

import (
	"encoding/json"
	"time"
)

const dateForm = "2006-01-02"

//Event - событие в календаре
//Num - не уникальный айдишник, а порядковый номер события
//на определенную дату
type Event struct {
	UserID      string    `json:"user_id"`
	Date        time.Time `json:"date"`
	Num         int       `json:"event_num"`
	Description string    `json:"description"`
}

//ErrorResponse - ошибка
type ErrorResponse struct {
	ErrorMsg string `json:"error"`
}

//ResultResponse - формат ответа на запрос,
//содержащего результат. Общий для
//GET и POST, но у них разный
//Result. Сериализуется в JSON
type ResultResponse struct {
	Result json.Marshaler `json:"result"`
}

//CRUDResult - формат результата ответа на
//POST запрос. Сериализуется в JSON
type CRUDResult struct {
	Num     int  `json:"event_num"`
	Success bool `json:"success"`
}

//GetResult - формат результата ответа на
//GET запрос. Сериализуется в JSON
type GetResult struct {
	Events []Event `json:"events"`
}

//Repository - интерфейс для
//работы с событиями
type Repository interface {
	createEvent(event Event) (num int, err error)
	updateEvent(event Event) error
	deleteEvent(event Event) error
	getEventsByDay(id string, date time.Time) []Event
	getEventsByMonth(id string, date time.Time) []Event
	getEventsByWeek(id string, date time.Time) []Event
}
